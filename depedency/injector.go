//go:build wireinject
// +build wireinject

package depedency

import (
	"database/sql"
	"restapi-bus/app"
	"restapi-bus/controller"
	croncustom "restapi-bus/cron_custom"
	"restapi-bus/external"
	"restapi-bus/middleware"
	"restapi-bus/repository"
	"restapi-bus/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
)

func InitializedControllerCustomer(db *sql.DB, rdb *middleware.RedisClientDb) controller.CustomerControllerInterface {
	wire.Build(
		repository.NewCustomerRepository,
		service.NewCustomerService,
		controller.NewCustomerController)
	return nil
}

func InitializedControllerAgency(db *sql.DB, rdb *middleware.RedisClientDb) controller.AgencyControllerInterface {
	wire.Build(
		repository.NewAgencyRepository,
		service.NewAgencyService,
		controller.NewAgencyController)
	return nil
}

func InitializedControllerBus(db *sql.DB, rdb *middleware.RedisClientDb) controller.BusControllerInterface {
	wire.Build(repository.NewBusRepository,
		repository.NewAgencyRepository,
		service.NewBusService,
		controller.NewBusController)
	return nil
}

func InitializedControllerDriver(db *sql.DB, rdb *middleware.RedisClientDb) controller.ControllerDriverInterface {
	wire.Build(repository.NewDiverRepository,
		repository.NewAgencyRepository,
		service.NewServiceDriver,
		controller.NewDriverController)
	return nil
}

func InitializedControllerSchedule(db *sql.DB, rdb *middleware.RedisClientDb, cronJob croncustom.InterfaceCronJob) controller.ControllerScheduleInterface {
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

func InitializedControllerTicket(
	db *sql.DB,
	rdb *middleware.RedisClientDb,
	rmq *amqp091.Channel,
	paymid external.InterfacePayment,
	cronJob croncustom.InterfaceCronJob,
) controller.ControllerTicketInterface {
	wire.Build(
		repository.NewTicketRepository,
		repository.NewCustomerRepository,
		repository.NewDiverRepository,
		repository.NewBusRepository,
		repository.NewScheduleRepository,
		repository.NewAgencyRepository,
		repository.BindMqChannel,
		service.NewTicketService,
		controller.NewTicketController)
	return nil
}

func InitializedServer(*sql.DB, *middleware.RedisClientDb, *amqp091.Channel, external.InterfacePayment, croncustom.InterfaceCronJob) *gin.Engine {
	wire.Build(
		InitializedControllerCustomer,
		InitializedControllerAgency,
		InitializedControllerBus,
		InitializedControllerDriver,
		InitializedControllerTicket,
		InitializedControllerSchedule,
		app.Router)
	return nil
}
