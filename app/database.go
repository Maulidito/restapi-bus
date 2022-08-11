package app

import (
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(username string, password string, dbName string, hostDb string) *sql.DB {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostDb, dbName)
	db, err := sql.Open("mysql", dataSourceName)

	helper.PanicIfError(err)

	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	helper.PanicIfError(err)

	return db

}
