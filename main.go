package main

import (
	"restapi-bus/app"
	"restapi-bus/depedency"
)

func main() {

	db := app.NewDatabase()

	server := depedency.InitializedServer(db)

	server.Run(":8080")
}
