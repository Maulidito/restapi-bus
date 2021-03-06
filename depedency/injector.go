//go:build wireinject
// +build wireinject

package depedency

import (
	"database/sql"
	"restapi-bus/app"
	"restapi-bus/controller"
	"restapi-bus/repository"
	"restapi-bus/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializedControllerCustomer(db *sql.DB) controller.CustomerControllerInterface {
	wire.Build(
		repository.NewCustomerRepository,
		service.NewCustomerService,
		controller.NewCustomerController)
	return nil
}

func InitializedControllerAgency(db *sql.DB) controller.AgencyControllerInterface {
	wire.Build(
		repository.NewAgencyRepository,
		service.NewAgencyService,
		controller.NewAgencyController)
	return nil
}

func InitializedControllerBus(db *sql.DB) controller.BusControllerInterface {
	wire.Build(repository.NewBusRepository,
		repository.NewAgencyRepository,
		service.NewBusService,
		controller.NewBusController)
	return nil
}

func InitializedControllerDriver(db *sql.DB) controller.ControllerDriverInterface {
	wire.Build(repository.NewDiverRepository,
		repository.NewAgencyRepository,
		service.NewServiceDriver,
		controller.NewDriverController)
	return nil
}

func InitializedControllerSchedule(db *sql.DB) controller.ControllerScheduleInterface {
	wire.Build(
		repository.NewScheduleRepository,
		repository.NewBusRepository,
		repository.NewAgencyRepository,
		repository.NewDiverRepository,
		service.NewScheduleService,
		controller.NewScheduleController,
	)
	return nil
}

func InitializedControllerTicket(db *sql.DB) controller.ControllerTicketInterface {
	wire.Build(
		repository.NewTicketRepository,
		repository.NewCustomerRepository,
		repository.NewDiverRepository,
		repository.NewBusRepository,
		repository.NewAgencyRepository,
		service.NewTicketService,
		controller.NewTicketController)
	return nil
}

func InitializedServer() *gin.Engine {
	wire.Build(
		app.NewDatabase,
		InitializedControllerCustomer,
		InitializedControllerAgency,
		InitializedControllerBus,
		InitializedControllerDriver,
		InitializedControllerTicket,
		InitializedControllerSchedule,
		app.Router)
	return nil
}
