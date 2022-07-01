package main

import (
	"restapi-bus/app"
	"restapi-bus/controller"
	"restapi-bus/repository"
	"restapi-bus/service"
)

func main() {

	db := app.NewDatabase()
	repo := repository.NewCustomerRepository()
	service := service.NewCustomerService(db, &repo)

	controller := controller.NewCustomerController(service)

	server := app.Router(controller)

	server.Run(":8080")
}
