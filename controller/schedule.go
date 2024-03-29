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
	"sync"

	"github.com/gin-gonic/gin"
)

type ControllerScheduleInterface interface {
	GetAllSchedule(ctx *gin.Context)
	GetOneSchedule(ctx *gin.Context)
	AddSchedule(ctx *gin.Context)
	DeleteSchedule(ctx *gin.Context)
	SetAutoSchedule(ctx *gin.Context)
	GetAutoSchedule(ctx *gin.Context)
	DeleteAutoSchedule(ctx *gin.Context)
	UpdateArrivedSchedule(ctx *gin.Context)
	RouterMount(g gin.IRouter)
}

type ControllerScheduleImplementation struct {
	Service entity.ScheduleServiceInterface
	Rdb     *middleware.RedisClientDb
}

func NewScheduleController(serv entity.ScheduleServiceInterface, rdb *middleware.RedisClientDb) ControllerScheduleInterface {
	return &ControllerScheduleImplementation{Service: serv, Rdb: rdb}
}

func (controller *ControllerScheduleImplementation) RouterMount(g gin.IRouter) {
	grouterSchedule := g.Group("/schedule")
	grouterScheduleAuth := grouterSchedule.Group("", middleware.MiddlewareAuth)
	grouterScheduleRdb := grouterSchedule.Group("", controller.Rdb.MiddlewareGetDataRedis)
	grouterSchedule.GET("", controller.GetAllSchedule)
	grouterScheduleAuth.POST("", controller.AddSchedule)
	grouterAutoSchedule := grouterScheduleAuth.Group("/autoschedule")
	grouterAutoSchedule.POST("", controller.SetAutoSchedule)
	grouterAutoSchedule.GET("", controller.GetAutoSchedule)
	grouterAutoSchedule.DELETE("", controller.DeleteAutoSchedule)
	grouterScheduleRdb.GET("/:scheduleId", controller.GetOneSchedule, controller.Rdb.MiddlewareSetDataRedis)
	grouterScheduleAuth.DELETE("/:scheduleId", controller.DeleteSchedule)
	grouterScheduleAuth.PATCH("/:scheduleId/arrived", controller.UpdateArrivedSchedule)
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
	ctx.Set("response", finalResponse)
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
	var err error
	scheduleArrived := request.ScheduleArrived{}
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

		ctx.Bind(&scheduleArrived)
	}()

	wg.Wait()

	helper.PanicIfError(err)

	responseSchedule := controller.Service.UpdateArrivedSchedule(ctx, scheduleIdInt, *scheduleArrived.IsArrived)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: responseSchedule}

	ctx.JSON(finalResponse.Code, finalResponse)
}

func (controller *ControllerScheduleImplementation) SetAutoSchedule(ctx *gin.Context) {

	automateSchedule := request.AutoSchedule{}

	err := ctx.ShouldBind(&automateSchedule)
	helper.PanicIfError(err)

	controller.Service.SetAutoSchedule(ctx, &automateSchedule)

	finalResponse := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}

	ctx.JSON(http.StatusOK, finalResponse)

}

func (controller *ControllerScheduleImplementation) GetAutoSchedule(ctx *gin.Context) {

	listCronJob := controller.Service.GetAutoSchedule(ctx)

	if len(listCronJob) == 0 {
		ctx.JSON(
			http.StatusNotFound,
			web.WebResponseNoData{
				Code: http.StatusOK, Status: http.StatusText(http.StatusOK),
			})
		return
	}
	ctx.JSON(
		http.StatusOK,
		web.WebResponse{
			Code: http.StatusOK, Status: http.StatusText(http.StatusOK),
			Data: listCronJob,
		})

}

func (controller *ControllerScheduleImplementation) DeleteAutoSchedule(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		panic(exception.NewBadRequestError("ERROR ID NOT FOUND"))

	}
	controller.Service.DeleteAutoSchedule(ctx, id)
}
