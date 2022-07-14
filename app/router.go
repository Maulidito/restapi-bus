package app

import (
	"os"
	"restapi-bus/controller"
	"restapi-bus/middleware"

	"github.com/gin-gonic/gin"
)

func configurationRouter() *gin.Engine {
	fileLog, _ := os.OpenFile("./log/logging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileRecovery, _ := os.OpenFile("./log/recoveryLog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	g := gin.Default()

	g.Use(gin.LoggerWithWriter(fileLog))
	g.Use(gin.CustomRecoveryWithWriter(fileRecovery, middleware.MiddlewarePanic))

	return g
}

func Router(customer controller.CustomerControllerInterface, agency controller.AgencyControllerInterface, bus controller.BusControllerInterface, driver controller.ControllerDriverInterface, ticket controller.ControllerTicketInterface) *gin.Engine {

	g := configurationRouter()

	grouter := g.Group("/v1")

	grouterCustomer := grouter.Group("/customer")

	grouterCustomer.GET("/", customer.GetAllCustomer)
	grouterCustomer.POST("/", customer.AddCustomer)
	grouterCustomer.GET("/:customerId", customer.GetOneCustomer)
	grouterCustomer.DELETE("/:customerId", customer.DeleteOneCustomer)

	grouterAgency := grouter.Group("/agency")

	grouterAgency.GET("/", agency.GetAllAgency)
	grouterAgency.POST("/", agency.AddAgency)
	grouterAgency.GET("/:agencyId", agency.GetOneAgency)
	grouterAgency.DELETE("/:agencyId", agency.DeleteOneAgency)

	grouterBus := grouter.Group("/bus")

	grouterBus.GET("/", bus.GetAllBus)
	grouterBus.POST("/", bus.AddBus)
	grouterBus.GET("/:busId", bus.GetOneBusOnSpecificAgency)
	grouterBus.GET("/agency/:agencyId", bus.GetAllBusOnSpecificAgency)
	grouterBus.DELETE("/:busId", bus.DeleteOneBus)

	grouterDriver := grouter.Group("/driver")
	grouterDriver.GET("/", driver.GetAllDriver)
	grouterDriver.GET("/filter", driver.GetAllDriver)
	grouterDriver.GET("/:driverId", driver.GetOneDriverOnSpecificAgency)
	grouterDriver.POST("/", driver.AddDriver)
	grouterDriver.GET("/agency/:agencyId", driver.GetAllDriverOnSpecificAgency)
	grouterDriver.DELETE("/:driverId", driver.DeleteDriver)

	grouterTicket := grouter.Group("/ticket")

	grouterTicket.GET("/", ticket.GetAllTicket)
	grouterTicket.GET("/:ticketId", ticket.GetOneTicket)
	grouterTicket.GET("/driver/:driverId", ticket.GetAllTicketOnSpecificDriver)
	grouterTicket.GET("/customer/:customerId", ticket.GetAllTicketOnSpecificCustomer)
	grouterTicket.GET("/bus/:busId", ticket.GetAllTicketOnSpecificBus)
	grouterTicket.GET("/agency/:agencyId", ticket.GetAllTicketOnSpecificAgency)

	grouterTicket.POST("/", ticket.AddTicket)
	grouterTicket.DELETE("/:ticketId", ticket.DeleteTicket)

	grouterTicket.PATCH("/:ticketId/arrived", ticket.UpdateArrivedTicket)

	return g
}
