// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package depedency

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"restapi-bus/app"
	"restapi-bus/controller"
	"restapi-bus/cron_custom"
	"restapi-bus/external"
	"restapi-bus/middleware"
	"restapi-bus/repository"
	"restapi-bus/service"
)

// Injectors from injector.go:

func InitializedControllerCustomer(db *sql.DB, rdb *middleware.RedisClientDb) controller.CustomerControllerInterface {
	customerRepositoryInterface := repository.NewCustomerRepository(db)
	customerServiceInterface := service.NewCustomerService(customerRepositoryInterface)
	customerControllerInterface := controller.NewCustomerController(customerServiceInterface, rdb)
	return customerControllerInterface
}

func InitializedControllerAgency(db *sql.DB, rdb *middleware.RedisClientDb) controller.AgencyControllerInterface {
	agencyRepositoryInterface := repository.NewAgencyRepository(db)
	agencyServiceInterface := service.NewAgencyService(agencyRepositoryInterface)
	agencyControllerInterface := controller.NewAgencyController(agencyServiceInterface, rdb)
	return agencyControllerInterface
}

func InitializedControllerBus(db *sql.DB, rdb *middleware.RedisClientDb) controller.BusControllerInterface {
	busRepositoryInterface := repository.NewBusRepository(db)
	agencyRepositoryInterface := repository.NewAgencyRepository(db)
	busServiceInterface := service.NewBusService(busRepositoryInterface, agencyRepositoryInterface)
	busControllerInterface := controller.NewBusController(busServiceInterface, rdb)
	return busControllerInterface
}

func InitializedControllerDriver(db *sql.DB, rdb *middleware.RedisClientDb) controller.ControllerDriverInterface {
	driverRepositoryInterface := repository.NewDiverRepository(db)
	agencyRepositoryInterface := repository.NewAgencyRepository(db)
	serviceDriverInterface := service.NewServiceDriver(driverRepositoryInterface, agencyRepositoryInterface)
	controllerDriverInterface := controller.NewDriverController(serviceDriverInterface, rdb)
	return controllerDriverInterface
}

func InitializedControllerSchedule(db *sql.DB, rdb *middleware.RedisClientDb, cronJob croncustom.InterfaceCronJob) controller.ControllerScheduleInterface {
	scheduleRepositoryInterface := repository.NewScheduleRepository(db)
	agencyRepositoryInterface := repository.NewAgencyRepository(db)
	driverRepositoryInterface := repository.NewDiverRepository(db)
	busRepositoryInterface := repository.NewBusRepository(db)
	scheduleServiceInterface := service.NewScheduleService(scheduleRepositoryInterface, agencyRepositoryInterface, driverRepositoryInterface, cronJob, busRepositoryInterface)
	controllerScheduleInterface := controller.NewScheduleController(scheduleServiceInterface, rdb)
	return controllerScheduleInterface
}

func InitializedControllerTicket(db *sql.DB, rdb *middleware.RedisClientDb, rmq *amqp091.Channel, paymid external.InterfacePayment, cronJob croncustom.InterfaceCronJob) controller.ControllerTicketInterface {
	busRepositoryInterface := repository.NewBusRepository(db)
	customerRepositoryInterface := repository.NewCustomerRepository(db)
	driverRepositoryInterface := repository.NewDiverRepository(db)
	ticketRepositoryInterface := repository.NewTicketRepository(db)
	agencyRepositoryInterface := repository.NewAgencyRepository(db)
	scheduleRepositoryInterface := repository.NewScheduleRepository(db)
	iMessageChannel := repository.BindMqChannel(rmq)
	ticketServiceInterface := service.NewTicketService(db, busRepositoryInterface, customerRepositoryInterface, driverRepositoryInterface, ticketRepositoryInterface, agencyRepositoryInterface, scheduleRepositoryInterface, iMessageChannel, paymid, cronJob)
	controllerTicketInterface := controller.NewTicketController(ticketServiceInterface, rdb, rmq)
	return controllerTicketInterface
}

func InitializedServer(db *sql.DB, redisClientDb *middleware.RedisClientDb, channel *amqp091.Channel, interfacePayment external.InterfacePayment, interfaceCronJob croncustom.InterfaceCronJob) *gin.Engine {
	customerControllerInterface := InitializedControllerCustomer(db, redisClientDb)
	agencyControllerInterface := InitializedControllerAgency(db, redisClientDb)
	busControllerInterface := InitializedControllerBus(db, redisClientDb)
	controllerDriverInterface := InitializedControllerDriver(db, redisClientDb)
	controllerTicketInterface := InitializedControllerTicket(db, redisClientDb, channel, interfacePayment, interfaceCronJob)
	controllerScheduleInterface := InitializedControllerSchedule(db, redisClientDb, interfaceCronJob)
	engine := app.Router(customerControllerInterface, agencyControllerInterface, busControllerInterface, controllerDriverInterface, controllerTicketInterface, controllerScheduleInterface, interfaceCronJob)
	return engine
}
