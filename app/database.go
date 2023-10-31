package app

import (
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(username string, password string, dbName string, hostDb string, portDB string) *sql.DB {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostDb, portDB, dbName)
	db, err := sql.Open("mysql", dataSourceName)

	helper.PanicIfError(err)

	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	helper.PanicIfError(err)

	return db

}

func NewDbManager(db *sql.DB) database.TrInterface {
	return &database.Tr{Db: db}
}
