package app

import (
	"io/fs"
	"os"
	"restapi-bus/controller"
	"restapi-bus/helper"
	"restapi-bus/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func configurationRouter() *gin.Engine {
	_, err := os.ReadDir("./log")
	if err != nil {
		err = os.Mkdir("./log", fs.FileMode(int(0766)))
	}
	helper.PanicIfError(err)
	fileLog, _ := os.OpenFile("./log/logging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileRecovery, _ := os.OpenFile("./log/recoveryLog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	g := gin.Default()

	g.Use(gin.LoggerWithWriter(fileLog))
	g.Use(gin.CustomRecoveryWithWriter(fileRecovery, middleware.MiddlewarePanic))

	return g
}

func IntializedCustomValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatefromTodate", helper.ValidateFromToDate)
		v.RegisterValidation("isbool", helper.IsBool)
	}
}

func Router(customer controller.CustomerControllerInterface, agency controller.AgencyControllerInterface, bus controller.BusControllerInterface, driver controller.ControllerDriverInterface, ticket controller.ControllerTicketInterface, schedule controller.ControllerScheduleInterface) *gin.Engine {

	g := configurationRouter()

	IntializedCustomValidation()

	grouter := g.Group("/v1")

	customer.RouterMount(grouter)

	driver.RouterMount(grouter)

	bus.RouterMount(grouter)

	agency.RouterMount(grouter)

	grouterTicket := grouter.Group("/ticket")

	grouterTicket.GET("/", ticket.GetAllTicket)
	grouterTicket.GET("/:ticketId", ticket.GetOneTicket)
	grouterTicket.GET("/driver/:driverId", ticket.GetAllTicketOnSpecificDriver)
	grouterTicket.GET("/customer/:customerId", ticket.GetAllTicketOnSpecificCustomer)
	grouterTicket.GET("/bus/:busId", ticket.GetAllTicketOnSpecificBus)
	grouterTicket.GET("/agency/:agencyId", ticket.GetAllTicketOnSpecificAgency)
	grouterTicket.GET("/price", ticket.GetTotalPriceAllTicket)
	grouterTicket.GET("/agency/:agencyId/price", ticket.GetTotalPriceTicketFromSpecificAgency)
	grouterTicket.GET("/driver/:driverId/price", ticket.GetTotalPriceTicketFromSpecificDriver)

	grouterTicket.POST("/", ticket.AddTicket)
	grouterTicket.DELETE("/:ticketId", ticket.DeleteTicket)

	grouterTicket.PATCH("/:ticketId/arrived", ticket.UpdateArrivedTicket)

	schedule.RouterMount(grouter)

	return g
}
