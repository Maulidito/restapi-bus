package main

import (
	"fmt"
	"os"
	"restapi-bus/app"
	"restapi-bus/depedency"
)

func main() {
	//err := godotenv.Load()
	//helper.PanicIfError(err)

	port := os.Getenv("PORT")
	usernameDb := os.Getenv("USERNAME_DB")
	hostDb := os.Getenv("HOST_DB")
	passDb := os.Getenv("PASSWORD_DB")
	nameDb := os.Getenv("NAME_DB")
	db := app.NewDatabase(usernameDb, passDb, nameDb, hostDb)

	server := depedency.InitializedServer(db)
	fmt.Println("SERVER RUNNING ON PORT ", port)
	server.Run(":" + port)

}
