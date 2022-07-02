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

var controllerCustomerSet = wire.NewSet(
	repository.NewCustomerRepository,
	service.NewCustomerService,
	controller.NewCustomerController)

var controllerAgencySet = wire.NewSet(
	repository.NewAgencyRepository,
	service.NewAgencyService,
	controller.NewAgencyController)

func InitializedServer(db *sql.DB) *gin.Engine {
	wire.Build(
		controllerCustomerSet,
		controllerAgencySet,
		app.Router)
	return nil
}
