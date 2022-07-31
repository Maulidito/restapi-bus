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

func InitializedControllerDriver(db *sql.DB) controller.ControllerDriverInterface {
	driverRepositoryInterface := repository.NewDiverRepository()
	agencyRepositoryInterface := repository.NewAgencyRepository()
	serviceDriverInterface := service.NewServiceDriver(db, driverRepositoryInterface, agencyRepositoryInterface)
	controllerDriverInterface := controller.NewDriverController(serviceDriverInterface)
	return controllerDriverInterface
}

func InitializedControllerSchedule(db *sql.DB) controller.ControllerScheduleInterface {
	scheduleRepositoryInterface := repository.NewScheduleRepository()
	agencyRepositoryInterface := repository.NewAgencyRepository()
	driverRepositoryInterface := repository.NewDiverRepository()
	busRepositoryInterface := repository.NewBusRepository()
	scheduleServiceInterface := service.NewScheduleService(scheduleRepositoryInterface, agencyRepositoryInterface, driverRepositoryInterface, busRepositoryInterface, db)
	controllerScheduleInterface := controller.NewScheduleController(scheduleServiceInterface)
	return controllerScheduleInterface
}

func InitializedControllerTicket(db *sql.DB) controller.ControllerTicketInterface {
	busRepositoryInterface := repository.NewBusRepository()
	customerRepositoryInterface := repository.NewCustomerRepository()
	driverRepositoryInterface := repository.NewDiverRepository()
	ticketRepositoryInterface := repository.NewTicketRepository()
	agencyRepositoryInterface := repository.NewAgencyRepository()
	ticketServiceInterface := service.NewTicketService(db, busRepositoryInterface, customerRepositoryInterface, driverRepositoryInterface, ticketRepositoryInterface, agencyRepositoryInterface)
	controllerTicketInterface := controller.NewTicketController(ticketServiceInterface)
	return controllerTicketInterface
}

func InitializedServer() *gin.Engine {
	db := app.NewDatabase()
	customerControllerInterface := InitializedControllerCustomer(db)
	agencyControllerInterface := InitializedControllerAgency(db)
	busControllerInterface := InitializedControllerBus(db)
	controllerDriverInterface := InitializedControllerDriver(db)
	controllerTicketInterface := InitializedControllerTicket(db)
	controllerScheduleInterface := InitializedControllerSchedule(db)
	engine := app.Router(customerControllerInterface, agencyControllerInterface, busControllerInterface, controllerDriverInterface, controllerTicketInterface, controllerScheduleInterface)
	return engine
}
