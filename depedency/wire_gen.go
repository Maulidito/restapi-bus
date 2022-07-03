// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package depedency

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"restapi-bus/app"
	"restapi-bus/controller"
	"restapi-bus/repository"
	"restapi-bus/service"
)

// Injectors from injector.go:

func InitializedControllerCustomer(db *sql.DB) controller.CustomerControllerInterface {
	customerRepositoryInterface := repository.NewCustomerRepository()
	customerServiceInterface := service.NewCustomerService(db, customerRepositoryInterface)
	customerControllerInterface := controller.NewCustomerController(customerServiceInterface)
	return customerControllerInterface
}

func InitializedControllerAgency(db *sql.DB) controller.AgencyControllerInterface {
	agencyRepositoryInterface := repository.NewAgencyRepository()
	agencyServiceInterface := service.NewAgencyService(db, agencyRepositoryInterface)
	agencyControllerInterface := controller.NewAgencyController(agencyServiceInterface)
	return agencyControllerInterface
}

func InitializedControllerBus(db *sql.DB) controller.BusControllerInterface {
	busRepositoryInterface := repository.NewBusRepository()
	agencyRepositoryInterface := repository.NewAgencyRepository()
	busServiceInterface := service.NewBusService(db, busRepositoryInterface, agencyRepositoryInterface)
	busControllerInterface := controller.NewBusController(busServiceInterface)
	return busControllerInterface
}

func InitializedServer() *gin.Engine {
	db := app.NewDatabase()
	customerControllerInterface := InitializedControllerCustomer(db)
	agencyControllerInterface := InitializedControllerAgency(db)
	busControllerInterface := InitializedControllerBus(db)
	engine := app.Router(customerControllerInterface, agencyControllerInterface, busControllerInterface)
	return engine
}
