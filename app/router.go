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

func Router(customer controller.CustomerControllerInterface, agency controller.AgencyControllerInterface, bus controller.BusControllerInterface, driver controller.ControllerDriverInterface) *gin.Engine {

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

	grouterBusOnSpecificAgency := grouterAgency.Group("/:agencyId/bus")

	grouter.GET("/bus", bus.GetAllBus)
	grouterBusOnSpecificAgency.GET("/", bus.GetAllBusOnSpecificAgency)
	grouterBusOnSpecificAgency.POST("/", bus.AddBus)
	grouterBusOnSpecificAgency.GET("/:busId", bus.GetOneBusOnSpecificAgency)
	grouterBusOnSpecificAgency.DELETE("/busId", bus.DeleteOneBus)

	grouterDriverOnSpecificAgency := grouterAgency.Group("/:agencyId/driver")

	grouter.GET("/driver", driver.GetAllDriver)
	grouterDriverOnSpecificAgency.GET("/", driver.GetAllDriverOnSpecificAgency)
	grouterDriverOnSpecificAgency.POST("/", driver.AddDriver)
	grouterDriverOnSpecificAgency.GET("/:driverId", driver.GetOneDriverOnSpecificAgency)
	grouterDriverOnSpecificAgency.DELETE("/:driverId", driver.DeleteDriver)

	return g
}
