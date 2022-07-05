package app

import (
	"database/sql"
	"restapi-bus/helper"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase() *sql.DB {

	db, err := sql.Open("mysql", "root:@/db_bus")

	helper.PanicIfError(err)

	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	helper.PanicIfError(err)

	return db

}
