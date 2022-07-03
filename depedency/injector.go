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

func InitializedServer() *gin.Engine {
	wire.Build(
		app.NewDatabase,
		InitializedControllerCustomer,
		InitializedControllerAgency,
		InitializedControllerBus,
		app.Router)
	return nil
}
