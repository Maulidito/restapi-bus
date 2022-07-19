package controller

import (
	"net/http"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/request"
	"restapi-bus/models/web"
	"restapi-bus/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ControllerDriverInterface interface {
	GetAllDriver(ctx *gin.Context)
	GetAllDriverOnSpecificAgency(ctx *gin.Context)
	GetOneDriverOnSpecificAgency(ctx *gin.Context)
	AddDriver(ctx *gin.Context)
	DeleteDriver(ctx *gin.Context)
}

type ControllerDriverImplementation struct {
	service service.ServiceDriverInterface
}

func NewDriverController(serv service.ServiceDriverInterface) ControllerDriverInterface {
	return &ControllerDriverImplementation{service: serv}
}

func (controller *ControllerDriverImplementation) GetAllDriver(ctx *gin.Context) {
	request := request.DriverFilter{}
	err := ctx.ShouldBindQuery(&request)

	helper.PanicIfError(err)
	listDriver := controller.service.GetAllDriver(ctx, &request)

	ctx.JSON(http.StatusOK, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: &listDriver})
}
func (controller *ControllerDriverImplementation) GetAllDriverOnSpecificAgency(ctx *gin.Context) {
	idAgency, idAgencyBool := ctx.Params.Get("agencyId")

	if !idAgencyBool {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT FOUND"))

	}

	idAgencyInt, err := strconv.Atoi(idAgency)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT INTEGER"))
	}
	finalResponse := controller.service.GetAllDriverOnSpecificAgency(ctx, idAgencyInt)

	ctx.JSON(http.StatusOK, &web.WebResponse{Code: http.StatusOK, Status: "OK", Data: finalResponse})

}
func (controller *ControllerDriverImplementation) GetOneDriverOnSpecificAgency(ctx *gin.Context) {

	idDriver, idDriverBool := ctx.Params.Get("driverId")

	if !idDriverBool {
		panic(exception.NewBadRequestError("DRIVER ID NOT FOUND"))
	}

	idDriverInt, err := strconv.Atoi(idDriver)

	helper.PanicIfError(err)

	finalResponse := controller.service.GetOneDriverOnSpecificAgency(ctx, idDriverInt)

	ctx.JSON(http.StatusOK, &web.WebResponse{Code: http.StatusOK, Status: "OK", Data: finalResponse})
}

func (controller *ControllerDriverImplementation) AddDriver(ctx *gin.Context) {

	req := request.Driver{}
	err := ctx.ShouldBindWith(&req, binding.Form)

	helper.PanicIfError(err)

	controller.service.AddDriver(ctx, &req)
	response := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}
	ctx.JSON(http.StatusOK, response)
}
func (controller *ControllerDriverImplementation) DeleteDriver(ctx *gin.Context) {

	idDriver, idDriverBool := ctx.Params.Get("driverId")

	if !idDriverBool {
		panic(exception.NewBadRequestError("ERROR DRIVER ID NOT FOUND"))
	}

	idDriverInt, err := strconv.Atoi(idDriver)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR DRIVER ID NOT INTEGER"))
	}

	responseData := controller.service.DeleteDriver(ctx, idDriverInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: responseData}
	ctx.JSON(http.StatusOK, &finalResponse)

}
