package main

import (
	"os"
	"restapi-bus/app"
	"restapi-bus/depedency"
	"restapi-bus/helper"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	helper.PanicIfError(err)

	port := os.Getenv("PORT")
	usernameDb := os.Getenv("USERNAME_DB")
	passDb := os.Getenv("PASSWORD_DB")
	nameDb := os.Getenv("NAME_DB")
	db := app.NewDatabase(usernameDb, passDb, nameDb)

	server := depedency.InitializedServer(db)

	server.Run(port)

}
