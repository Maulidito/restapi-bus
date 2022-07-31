package controller

import (
	"net/http"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/request"
	"restapi-bus/models/web"
	"restapi-bus/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type ControllerScheduleInterface interface {
	GetAllSchedule(ctx *gin.Context)
	GetOneSchedule(ctx *gin.Context)
	AddSchedule(ctx *gin.Context)
	DeleteSchedule(ctx *gin.Context)
	UpdateArrivedSchedule(ctx *gin.Context)
	RouterMount(g *gin.RouterGroup)
}

type ControllerScheduleImplementation struct {
	Service service.ScheduleServiceInterface
}

func NewScheduleController(serv service.ScheduleServiceInterface) ControllerScheduleInterface {
	return &ControllerScheduleImplementation{Service: serv}
}

func (controller *ControllerScheduleImplementation) RouterMount(g *gin.RouterGroup) {
	grouterSchedule := g.Group("/schedule")
	grouterSchedule.GET("", controller.GetAllSchedule)
	grouterSchedule.POST("", controller.AddSchedule)
	grouterSchedule.GET("/:scheduleId", controller.GetOneSchedule)
	grouterSchedule.DELETE("/:scheduleId", controller.DeleteSchedule)
	grouterSchedule.PATCH("/:scheduleId/arrived", controller.UpdateArrivedSchedule)
}

func (controller *ControllerScheduleImplementation) GetAllSchedule(ctx *gin.Context) {

	requestFilter := request.ScheduleFilter{}
	err := ctx.Bind(&requestFilter)

	helper.PanicIfError(err)

	responseSchedule := controller.Service.GetAllSchedule(ctx, &requestFilter)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: responseSchedule}

	ctx.JSON(http.StatusOK, finalResponse)

}

func (controller *ControllerScheduleImplementation) GetOneSchedule(ctx *gin.Context) {
	scheduleIdString, ok := ctx.Params.Get("scheduleId")
	if !ok {
		panic(exception.NewBadRequestError("ERROR SCHEDULE ID NOT FOUND"))
	}

	scheduleIdInt, err := strconv.Atoi(scheduleIdString)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR SCHEDULE ID NOT NUMBER"))
	}

	responseSchedule := controller.Service.GetOneSchedule(ctx, scheduleIdInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: responseSchedule}
	ctx.JSON(http.StatusOK, finalResponse)
}

func (controller *ControllerScheduleImplementation) AddSchedule(ctx *gin.Context) {
	requestSchedule := request.Schedule{}
	err := ctx.Bind(&requestSchedule)
	helper.PanicIfError(err)

	controller.Service.AddSchedule(ctx, &requestSchedule)

	finalResponse := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}

	ctx.JSON(finalResponse.Code, finalResponse)
}

func (controller *ControllerScheduleImplementation) DeleteSchedule(ctx *gin.Context) {
	scheduleIdString, ok := ctx.Params.Get("scheduleId")
	if !ok {
		panic(exception.NewBadRequestError("ERROR SCHEDULE ID NOT FOUND"))
	}

	scheduleIdInt, err := strconv.Atoi(scheduleIdString)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR SCHEDULE ID NOT NUMBER"))
	}

	responseSchedule := controller.Service.DeleteSchedule(ctx, scheduleIdInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: responseSchedule}
	ctx.JSON(finalResponse.Code, finalResponse)

}

func (controller *ControllerScheduleImplementation) UpdateArrivedSchedule(ctx *gin.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var scheduleIdInt int
	var isArrivedBool bool
	var err error
	go func() {
		defer wg.Done()
		scheduleIdString, ok := ctx.Params.Get("scheduleId")
		if !ok {
			panic(exception.NewBadRequestError("ERROR SCHEDULE ID NOT FOUND"))
		}

		scheduleIdInt, err = strconv.Atoi(scheduleIdString)

		if err != nil {
			panic(exception.NewBadRequestError("ERROR SCHEDULE ID NOT NUMBER"))
		}
	}()

	go func() {
		defer wg.Done()
		isArrived, ok := ctx.GetPostForm("arrived")
		if !ok {
			panic(exception.NewBadRequestError("ERROR ARRIVED NOT FOUND"))
		}

		isArrivedBool, err = strconv.ParseBool(isArrived)

		if err != nil {
			panic(exception.NewBadRequestError("ERROR ARRIVED ID NOT BOOL"))
		}
	}()

	wg.Wait()

	helper.PanicIfError(err)

	responseSchedule := controller.Service.UpdateArrivedSchedule(ctx, scheduleIdInt, isArrivedBool)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: responseSchedule}

	ctx.JSON(finalResponse.Code, finalResponse)
}
