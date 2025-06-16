package sqlconnect

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"reflect"
	"restapi/internal/models"
	"restapi/utils"
	"strconv"
	"strings"
	"time"
)

func GetExecByID(realID int) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id = ?",
		realID,
	).Scan(&exec.ID,
		&exec.FirstName,
		&exec.LastName,
		&exec.Email,
		&exec.Username,
		&exec.UserCreatedAt,
		&exec.InactiveStatus,
		&exec.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.UnitNotFoundError
	} else if err != nil {
		return models.Exec{}, utils.DatabaseQueryError
	}
	return exec, nil
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

func GetExecsDBHandler(r *http.Request) ([]models.Exec, error) {
	query := "SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"
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

	var execs []models.Exec
	err = db.Select(&execs, query, args...)
	if err != nil {
		log.Println(err)
		return nil, utils.DatabaseQueryError
	}

	return execs, nil
}

func AddExecsDBHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO execs (first_name, last_name, email, username, password, role) VALUES (?, ?, ?, ?, ?, ?)") // will prepare SQL for execution
	if err != nil {
		log.Println("ERR 1:", err)
		return nil, utils.DatabaseQueryError
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))
	for i, exec := range newExecs {
		exec.Password, err = utils.Hash(exec.Password)
		if err != nil {
			return nil, err
		}
		res, err := stmt.Exec(exec.FirstName, exec.LastName, exec.Email, exec.Username, exec.Password, exec.Role)
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
		exec.ID = int(lastID)
		addedExecs[i] = exec
	}
	return addedExecs, nil
}

func PatchExecsDBHandler(updates []map[string]interface{}) error {
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
		var execFromDb models.Exec
		err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(
			&execFromDb.ID,
			&execFromDb.FirstName,
			&execFromDb.LastName,
			&execFromDb.Email,
			&execFromDb.Username)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return utils.UnitNotFoundError
		} else if err != nil {
			tx.Rollback()
			return utils.DatabaseQueryError
		}
		// apply updates using reflection
		execVal := reflect.ValueOf(&execFromDb).Elem()
		execType := execVal.Type()
		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k {
					fieldVal := execVal.Field(i)
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
		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?",
			execFromDb.FirstName,
			execFromDb.LastName,
			execFromDb.Email,
			execFromDb.Username)
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

func PatchOneExecDBHandler(id int, updates map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.DatabaseQueryError
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(
		&existingExec.ID,
		&existingExec.FirstName,
		&existingExec.LastName,
		&existingExec.Email,
		&existingExec.Username)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.UnitNotFoundError
	} else if err != nil {
		return models.Exec{}, utils.DatabaseQueryError
	}

	// apply updates using reflect
	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for key, value := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			if field.Tag.Get("json") == key {
				if execVal.Field(i).CanSet() {
					fieldVal := execVal.Field(i)
					fieldVal.Set(reflect.ValueOf(value).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?",
		existingExec.FirstName,
		existingExec.LastName,
		existingExec.Email,
		existingExec.Username,
		existingExec.ID)
	if err != nil {
		return models.Exec{}, utils.DatabaseQueryError
	}
	return existingExec, nil
}

func DeleteOneExecDBHandler(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ConnectingToDatabaseError
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
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

func LoginDBHandler(username string) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ConnectingToDatabaseError
	}
	defer db.Close()

	var user models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM execs WHERE username = ?", username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.InactiveStatus,
		&user.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.UnitNotFoundError
	} else if err != nil {
		log.Println(err)
		return models.Exec{}, utils.DatabaseQueryError
	}
	return user, nil
}

func UpdatePasswordInDB(userId int, newPassword, currentPassword string) (string, error) {
	db, err := ConnectDb()
	if err != nil {
		return "", utils.ConnectingToDatabaseError
	}
	var username string
	var userPassword string
	var userRole string
	err = db.QueryRow("SELECT username, password, role FROM execs WHERE id = ?", userId).Scan(&username, &userPassword, &userRole)

	if err != nil {
		return "", utils.UnitNotFoundError
	}

	_, err = utils.VerifyPassword(userPassword, currentPassword)
	if err != nil {
		return "", err
	}

	hashedPassword, err := utils.Hash(newPassword)
	if err != nil {
		return "", err
	}

	currentTime := time.Now().Format(time.RFC3339)
	_, err = db.Exec("UPDATE execs SET password = ?, password_changed_at = ? WHERE id = ?", hashedPassword, currentTime, userId)
	if err != nil {
		return "", utils.DatabaseQueryError
	}
	token, err := utils.SignToken(userId, username, userRole)
	if err != nil {
		return "", err
	}
	return token, nil
}
