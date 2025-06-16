package sqlconnect

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"restapi/utils"
)

func ConnectDb() (*sql.DB, error) {
	connectionString := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	if err = db.Ping(); err != nil {
		return nil, utils.ConnectingToDatabaseError
	}
	return db, nil
}
