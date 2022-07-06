package controller

import (
	"net/http"
	"restapi-bus/helper"
	"restapi-bus/models/request"
	"restapi-bus/models/web"
	"restapi-bus/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerControllerInterface interface {
	GetAllCustomer(ctx *gin.Context)
	AddCustomer(ctx *gin.Context)
	GetOneCustomer(ctx *gin.Context)
	DeleteOneCustomer(ctx *gin.Context)
}

type CustomerControllerImplementation struct {
	service service.CustomerServiceInterface
}

func NewCustomerController(service service.CustomerServiceInterface) CustomerControllerInterface {
	return &CustomerControllerImplementation{service: service}
}

func (ctrl *CustomerControllerImplementation) GetAllCustomer(ctx *gin.Context) {
	customerResponse := ctrl.service.GetAllCustomer(ctx)

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
		panic("ERROR ID NOT FOUND")
	}
	idInt, err := strconv.Atoi(id)
	helper.PanicIfError(err)
	customerResponse := ctrl.service.GetOneCustomer(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: customerResponse}

	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *CustomerControllerImplementation) DeleteOneCustomer(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("customerId")

	if !idBool {
		panic("ERROR ID NOT FOUND")
	}
	idInt, err := strconv.Atoi(id)
	helper.PanicIfError(err)
	customerResponse := ctrl.service.DeleteOneCustomer(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: customerResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
