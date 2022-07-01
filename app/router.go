package app

import (
	"restapi-bus/controller"
	"restapi-bus/middleware"

	"github.com/gin-gonic/gin"
)

func Router(customer controller.CustomerControllerInterface) *gin.Engine {
	g := gin.Default()
	g.Use(middleware.MiddlewarePanicHandler)

	grouter := g.Group("/v1")

	grouter.GET("/customer", customer.GetAllCustomer)
	grouter.POST("/customer/", customer.AddCustomer)
	grouter.GET("/customer/:customerId", customer.GetOneCustomer)
	grouter.DELETE("/customer/:customerId", customer.DeleteOneCustomer)

	return g
}
