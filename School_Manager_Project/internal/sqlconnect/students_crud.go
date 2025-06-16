package sqlconnect

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"reflect"
	"restapi/internal/models"
	"restapi/utils"
	"strconv"
	"strings"
)

func GetStudentByID(realID int) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, class FROM students WHERE id = ?",
		realID,
	).Scan(&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.UnitNotFoundError
	} else if err != nil {
		return models.Student{}, utils.DatabaseQueryError
	}
	return student, nil
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

func GetStudentsDBHandler(studentList []models.Student, r *http.Request) ([]models.Student, error) {
	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
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

	var students []models.Student
	err = db.Select(&students, query, args...)
	if err != nil {
		log.Println(err)
		return nil, utils.DatabaseQueryError
	}

	return students, nil
}

func AddStudentsDBHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class) VALUES (?, ?, ?, ?)") // will prepare SQL for execution
	if err != nil {
		log.Println("ERR 1:", err)
		return nil, utils.DatabaseQueryError
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))
	for i, student := range newStudents {
		res, err := stmt.Exec(student.FirstName, student.LastName, student.Email, student.Class)
		if err != nil {
			log.Println("ERR 2:", err)
			if strings.Contains(err.Error(), "Duplicate entry") {
				return nil, utils.DuplicateEmailError
			} else if strings.Contains(err.Error(), "a foreign key constraint fails") {
				return nil, utils.ClassTeacherNotFound
			}
			return nil, utils.DatabaseQueryError
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			log.Println("ERR 3:", err)
			return nil, utils.DatabaseQueryError
		}
		student.ID = int(lastID)
		addedStudents[i] = student
	}
	return addedStudents, nil
}

func UpdateStudentDBHandler(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class)
	if err != nil {
		return models.Student{}, utils.DatabaseQueryError
	}
	updatedStudent.ID = existingStudent.ID
	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		updatedStudent.FirstName,
		updatedStudent.LastName,
		updatedStudent.Email,
		updatedStudent.Class,
		updatedStudent.ID)
	if err != nil {
		return models.Student{}, utils.DatabaseQueryError
	}
	return updatedStudent, nil
}

func PatchStudentsDBHandler(updates []map[string]interface{}) error {
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
		var studentFromDb models.Student
		err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
			&studentFromDb.ID,
			&studentFromDb.FirstName,
			&studentFromDb.LastName,
			&studentFromDb.Email,
			&studentFromDb.Class)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return utils.UnitNotFoundError
		} else if err != nil {
			tx.Rollback()
			return utils.DatabaseQueryError
		}
		// apply updates using reflection
		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		studentType := studentVal.Type()
		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				if field.Tag.Get("json") == k {
					fieldVal := studentVal.Field(i)
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
		_, err = tx.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
			studentFromDb.FirstName,
			studentFromDb.LastName,
			studentFromDb.Email,
			studentFromDb.Class,
			studentFromDb.ID)
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

func PatchOneStudentDBHandler(id int, updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Student{}, utils.DatabaseQueryError
	}
	defer db.Close()

	var existingStudent models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&existingStudent.ID,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.UnitNotFoundError
	} else if err != nil {
		return models.Student{}, utils.DatabaseQueryError
	}

	// apply updates using reflect
	studentVal := reflect.ValueOf(&existingStudent).Elem()
	studentType := studentVal.Type()

	for key, value := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == key {
				if studentVal.Field(i).CanSet() {
					fieldVal := studentVal.Field(i)
					fieldVal.Set(reflect.ValueOf(value).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		existingStudent.FirstName,
		existingStudent.LastName,
		existingStudent.Email,
		existingStudent.Class,
		existingStudent.ID)
	if err != nil {
		return models.Student{}, utils.DatabaseQueryError
	}
	return existingStudent, nil
}

func DeleteOneStudentDBHandler(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ConnectingToDatabaseError
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
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

func DeleteStudentsDBHandler(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return nil, utils.UnableToStartTransactionError
	}

	stm, err := tx.Prepare("DELETE FROM students WHERE id = ?")
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
