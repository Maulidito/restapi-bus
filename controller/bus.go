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

type BusControllerInterface interface {
	GetAllBus(ctx *gin.Context)
	AddBus(ctx *gin.Context)
	GetOneBusOnSpecificAgency(ctx *gin.Context)
	DeleteOneBus(ctx *gin.Context)
	GetAllBusOnSpecificAgency(ctx *gin.Context)
}

type BusControllerImplementation struct {
	service service.BusServiceInterface
}

func NewBusController(service service.BusServiceInterface) BusControllerInterface {
	return &BusControllerImplementation{service: service}
}

func (ctrl *BusControllerImplementation) GetAllBus(ctx *gin.Context) {
	busResponse := ctrl.service.GetAllBus(ctx)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: busResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
func (ctrl *BusControllerImplementation) AddBus(ctx *gin.Context) {
	busRequest := request.Bus{}
	err := ctx.ShouldBind(&busRequest)
	helper.PanicIfError(err)
	ctrl.service.AddBus(ctx, &busRequest)

	finalResponse := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}
	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *BusControllerImplementation) GetOneBusOnSpecificAgency(ctx *gin.Context) {
	idBus, idBoolBus := ctx.Params.Get("busId")

	if !idBoolBus {
		panic("ERROR ID BUS PARAMAETER NOT FOUND")
	}

	idIntBus, err := strconv.Atoi(idBus)
	helper.PanicIfError(err)
	busResponse := ctrl.service.GetOneBusSpecificAgency(ctx, idIntBus)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: busResponse}

	ctx.JSON(http.StatusOK, finalResponse)

}

func (ctrl *BusControllerImplementation) GetAllBusOnSpecificAgency(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("agencyId")

	if !idBool {
		panic("ERROR ID agencyId PARAMETER NOT FOUND")
	}

	idInt, err := strconv.Atoi(id)
	helper.PanicIfError(err)

	dataResponse := ctrl.service.GetAllBusOnSpecificAgency(ctx, idInt)

	finalResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   dataResponse,
	}

	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *BusControllerImplementation) DeleteOneBus(ctx *gin.Context) {
	busId, idBoolBus := ctx.Params.Get("busId")

	if !idBoolBus {
		panic("ERROR ID busId PARAMAETER NOT FOUND")
	}

	idBusIdInt, err := strconv.Atoi(busId)
	helper.PanicIfError(err)

	busResponse := ctrl.service.DeleteOneBus(ctx, idBusIdInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: busResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
