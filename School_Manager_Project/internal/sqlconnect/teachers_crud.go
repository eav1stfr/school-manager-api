package sqlconnect

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"reflect"
	"restapi/internal/models"
	"restapi/utils"
	"strconv"
)

func GetTeacherByID(realID int) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Teacher{}, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?",
		realID,
	).Scan(&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Class,
		&teacher.Subject)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Teacher{}, utils.UnitNotFoundError
	} else if err != nil {
		return models.Teacher{}, utils.DatabaseQueryError
	}
	return teacher, nil
}

//func GetTeachersDBHandler(teacherList []models.Teacher, r *http.Request) ([]models.Teacher, error) {
//	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
//	var args []interface{}
//
//	query, args = addSearchFilters(r, query, args)
//	query, err := addSortFilters(r, query)
//
//	if err != nil {
//		return nil, err
//	}
//	db, err := ConnectDb()
//	if err != nil {
//		return nil, utils.ConnectingToDatabaseError
//	}
//	defer db.Close()
//
//	rows, err := db.Query(query, args...)
//	if err != nil {
//		log.Println(err)
//		return nil, utils.InvalidSearchOParametersError
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var teacher models.Teacher
//		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
//		if err != nil {
//			return nil, utils.DatabaseQueryError
//		}
//		teacherList = append(teacherList, teacher)
//	}
//	return teacherList, nil
//}

func GetTeachersDBHandler(r *http.Request) ([]models.Teacher, error) {
	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = utils.AddSearchFilters(r, query, args)
	query, err := utils.AddSortFilters(r, query)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("mysql", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var teachers []models.Teacher
	err = db.Select(&teachers, query, args...)
	if err != nil {
		log.Println(err)
		return nil, utils.DatabaseQueryError
	}

	return teachers, nil
}

func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)") // will prepare SQL for execution
	if err != nil {
		return nil, utils.DatabaseQueryError
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			return nil, utils.DatabaseQueryError
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.DatabaseQueryError
		}
		teacher.ID = int(lastID)
		addedTeachers[i] = teacher
	}
	return addedTeachers, nil
}

func UpdateTeacherDBHandler(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Teacher{}, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject)
	if err != nil {
		return models.Teacher{}, utils.DatabaseQueryError
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName,
		updatedTeacher.LastName,
		updatedTeacher.Email,
		updatedTeacher.Class,
		updatedTeacher.Subject,
		updatedTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.DatabaseQueryError
	}
	return updatedTeacher, nil
}

func PatchTeachersDBHandler(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ConnectingToDatabaseError
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.UnableToStartTransactionError
	}
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.InvalidIdError
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.InvalidIdError
		}
		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
			&teacherFromDb.ID,
			&teacherFromDb.FirstName,
			&teacherFromDb.LastName,
			&teacherFromDb.Email,
			&teacherFromDb.Class,
			&teacherFromDb.Subject)
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return utils.UnitNotFoundError
		} else if err != nil {
			tx.Rollback()
			return utils.DatabaseQueryError
		}
		// apply updates using reflection
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()
		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(field.Type) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("cannot convert %v to %v", val.Type(), fieldVal.Type())
							return utils.InvalidUpdateParametersError
						}
					}
					break
				}
			}
		}
		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
			teacherFromDb.FirstName,
			teacherFromDb.LastName,
			teacherFromDb.Email,
			teacherFromDb.Class,
			teacherFromDb.Subject,
			teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.DatabaseQueryError
		}
	}
	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return utils.ErrorCommitingTransaction
	}
	return nil
}

func PatchOneTeacherDBHandler(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Teacher{}, utils.DatabaseQueryError
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Teacher{}, utils.UnitNotFoundError
	} else if err != nil {
		return models.Teacher{}, utils.DatabaseQueryError
	}

	// apply updates using reflect
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for key, value := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == key {
				if teacherVal.Field(i).CanSet() {
					fieldVal := teacherVal.Field(i)
					fieldVal.Set(reflect.ValueOf(value).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName,
		existingTeacher.LastName,
		existingTeacher.Email,
		existingTeacher.Class,
		existingTeacher.Subject,
		existingTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.DatabaseQueryError
	}
	return existingTeacher, nil
}

func DeleteOneTeacherDBHandler(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ConnectingToDatabaseError
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.DatabaseQueryError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.DatabaseQueryError
	}
	if rowsAffected == 0 {
		return utils.UnitNotFoundError
	}
	return nil
}

func DeleteTeachersDBHandler(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return nil, utils.UnableToStartTransactionError
	}

	stm, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.DatabaseQueryError
	}
	defer db.Close()

	deletedIds := []int{}

	for _, id := range ids {
		res, err := stm.Exec(id)
		if err != nil {
			return nil, utils.DatabaseQueryError
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, utils.DatabaseQueryError
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

		if rowsAffected < 1 {
			tx.Rollback()
			localErr := utils.AppErrors{}
			localErr.SetErrorMessage(fmt.Sprintf("unit not found: %d", id))
			localErr.SetErrStatusCode(http.StatusNotFound)
			return nil, &localErr
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorCommitingTransaction
	}
	if len(deletedIds) < 1 {
		return nil, utils.UnitNotFoundError
	}
	return deletedIds, nil
}

func GetStudentsListForTeacherDBHandler(id int) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var studentsList []models.Student
	var class string
	err = db.QueryRow("SELECT class FROM teachers WHERE id = ?", id).Scan(&class)

	if err == sql.ErrNoRows {
		return nil, utils.UnitNotFoundError
	} else if err != nil {
		return nil, utils.DatabaseQueryError
	}

	rows, err := db.Query("SELECT id, first_name, last_name, email, class FROM students WHERE class = ?", class)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Class)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		studentsList = append(studentsList, student)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return studentsList, nil
}

func GetStudentCountForTeacherDBHandler(id int) (int, error) {
	db, err := ConnectDb()
	if err != nil {
		return 0, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var class string
	err = db.QueryRow("SELECT class FROM teachers WHERE id = ?", id).Scan(&class)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, utils.UnitNotFoundError
	} else if err != nil {
		return 0, utils.DatabaseQueryError
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM students WHERE class = ?", class).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, utils.UnitNotFoundError
	} else if err != nil {
		return 0, utils.DatabaseQueryError
	}
	return count, nil
}
