package controller

import (
	"fmt"
	"net/http"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/middleware"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/web"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerControllerInterface interface {
	GetAllCustomer(ctx *gin.Context)
	AddCustomer(ctx *gin.Context)
	GetOneCustomer(ctx *gin.Context)
	DeleteOneCustomer(ctx *gin.Context)
	RouterMount(g gin.IRouter)
}

type CustomerControllerImplementation struct {
	service entity.CustomerServiceInterface
	Rdb     *middleware.RedisClientDb
}

func NewCustomerController(service entity.CustomerServiceInterface, rdb *middleware.RedisClientDb) CustomerControllerInterface {
	return &CustomerControllerImplementation{service: service, Rdb: rdb}
}

func (ctrl *CustomerControllerImplementation) RouterMount(g gin.IRouter) {
	grouterCustomer := g.Group("/customer")
	grouterCustomerAuth := grouterCustomer.Group("", middleware.MiddlewareAuth)
	grouterCustomerRdb := g.Group("", ctrl.Rdb.MiddlewareGetDataRedis)
	grouterCustomer.GET("", ctrl.GetAllCustomer)
	grouterCustomerAuth.POST("", ctrl.AddCustomer)
	grouterCustomerRdb.GET("/:customerId", ctrl.GetOneCustomer, ctrl.Rdb.MiddlewareSetDataRedis)
	grouterCustomerAuth.DELETE("/:customerId", ctrl.DeleteOneCustomer)
}

func (ctrl *CustomerControllerImplementation) GetAllCustomer(ctx *gin.Context) {
	filter := request.CustomerFilter{}
	fmt.Println("CHECK ERR", filter)
	err := ctx.ShouldBindQuery(&filter)
	helper.PanicIfError(err)

	customerResponse := ctrl.service.GetAllCustomer(ctx, &filter)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: customerResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
func (ctrl *CustomerControllerImplementation) AddCustomer(ctx *gin.Context) {
	customerRequest := request.Customer{}
	err := ctx.ShouldBind(&customerRequest)
	helper.PanicIfError(err)
	ctrl.service.AddCustomer(ctx, &customerRequest)

	finalResponse := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}
	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *CustomerControllerImplementation) GetOneCustomer(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("customerId")

	if !idBool {
		panic(exception.NewBadRequestError("ERROR CUSTOMER ID NOT FOUND"))
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR CUSTOMER ID NOT INTEGER"))
	}
	customerResponse := ctrl.service.GetOneCustomer(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: customerResponse}

	ctx.Set("response", finalResponse)

	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *CustomerControllerImplementation) DeleteOneCustomer(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("customerId")

	if !idBool {
		panic(exception.NewBadRequestError("ERROR CUSTOMER ID NOT FOUND"))
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR CUSTOMER ID NOT INTEGER"))
	}
	customerResponse := ctrl.service.DeleteOneCustomer(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: customerResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
