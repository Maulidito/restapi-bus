package controller

import (
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

type BusControllerInterface interface {
	GetAllBus(ctx *gin.Context)
	AddBus(ctx *gin.Context)
	GetOneBusOnSpecificAgency(ctx *gin.Context)
	DeleteOneBus(ctx *gin.Context)
	GetAllBusOnSpecificAgency(ctx *gin.Context)
	RouterMount(g gin.IRouter)
}

type BusControllerImplementation struct {
	service entity.BusServiceInterface
	Rdb     *middleware.RedisClientDb
}

func NewBusController(service entity.BusServiceInterface, rdb *middleware.RedisClientDb) BusControllerInterface {
	return &BusControllerImplementation{service: service, Rdb: rdb}
}

func (ctrl *BusControllerImplementation) RouterMount(g gin.IRouter) {
	grouterBus := g.Group("/bus")
	grouterBusAuth := grouterBus.Group("", middleware.MiddlewareAuth)
	grouterBusRdb := grouterBus.Group("", ctrl.Rdb.MiddlewareGetDataRedis)
	grouterBus.GET("", ctrl.GetAllBus)
	grouterBusAuth.POST("", ctrl.AddBus)
	grouterBusRdb.GET("/:busId", ctrl.GetOneBusOnSpecificAgency, ctrl.Rdb.MiddlewareSetDataRedis)
	grouterBus.GET("/agency/:agencyId", ctrl.GetAllBusOnSpecificAgency)
	grouterBusAuth.DELETE("/:busId", ctrl.DeleteOneBus)
}

func (ctrl *BusControllerImplementation) GetAllBus(ctx *gin.Context) {
	requestBusFilter := request.BusFilter{}
	err := ctx.Bind(&requestBusFilter)
	helper.PanicIfError(err)

	busResponse := ctrl.service.GetAllBus(ctx, &requestBusFilter)

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
		panic(exception.NewBadRequestError("ERROR BUS ID NOT FOUND"))
	}

	idIntBus, err := strconv.Atoi(idBus)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR BUS ID NOT INTEGER"))
	}
	busResponse := ctrl.service.GetOneBusSpecificAgency(ctx, idIntBus)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: busResponse}

	ctx.Set("response", finalResponse)

	ctx.JSON(http.StatusOK, finalResponse)

}

func (ctrl *BusControllerImplementation) GetAllBusOnSpecificAgency(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("agencyId")

	if !idBool {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT FOUND"))
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT INTEGER"))
	}

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
		panic(exception.NewBadRequestError("ERROR BUS ID NOT FOUND"))
	}

	idBusIdInt, err := strconv.Atoi(busId)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR BUS ID NOT INTEGER"))
	}

	busResponse := ctrl.service.DeleteOneBus(ctx, idBusIdInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: busResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
