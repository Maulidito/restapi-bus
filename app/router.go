package app

import (
	"io/fs"
	"os"
	"reflect"
	"restapi-bus/controller"
	"restapi-bus/helper"
	"restapi-bus/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func DefaultConfigurationRouter() *gin.Engine {
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
		v.RegisterValidation("validatedateafternow", helper.ValidateDateAfterNow)
		v.RegisterValidation("validateoneuppercase", helper.AtLeastOneUppercase)
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

}

func Router(
	customer controller.CustomerControllerInterface,
	agency controller.AgencyControllerInterface,
	bus controller.BusControllerInterface,
	driver controller.ControllerDriverInterface,
	ticket controller.ControllerTicketInterface,
	schedule controller.ControllerScheduleInterface,

) *gin.Engine {

	g := DefaultConfigurationRouter()

	IntializedCustomValidation()

	grouter := g.Group("/v1")

	customer.RouterMount(grouter)

	driver.RouterMount(grouter)

	bus.RouterMount(grouter)

	agency.RouterMount(grouter)

	ticket.RouterMount(grouter)

	schedule.RouterMount(grouter)

	return g
}
