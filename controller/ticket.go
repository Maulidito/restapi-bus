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
	"github.com/rabbitmq/amqp091-go"
)

type ControllerTicketInterface interface {
	GetAllTicket(ctx *gin.Context)
	AddTicket(ctx *gin.Context)
	GetOneTicket(ctx *gin.Context)
	DeleteTicket(ctx *gin.Context)
	GetAllTicketOnSpecificDriver(ctx *gin.Context)
	GetAllTicketOnSpecificCustomer(ctx *gin.Context)
	GetAllTicketOnSpecificAgency(ctx *gin.Context)
	GetAllTicketOnSpecificBus(ctx *gin.Context)
	GetTotalPriceAllTicket(ctx *gin.Context)
	GetTotalPriceTicketFromSpecificAgency(ctx *gin.Context)
	GetTotalPriceTicketFromSpecificDriver(ctx *gin.Context)
	RouterMount(g gin.IRouter)
}

type ControllerTicketImplementation struct {
	service  entity.TicketServiceInterface
	Rdb      *middleware.RedisClientDb
	Rabbitmq *amqp091.Channel
}

func NewTicketController(serv entity.TicketServiceInterface, rdb *middleware.RedisClientDb, rmq *amqp091.Channel) ControllerTicketInterface {
	return &ControllerTicketImplementation{service: serv, Rdb: rdb, Rabbitmq: rmq}
}

func (ctrl *ControllerTicketImplementation) RouterMount(g gin.IRouter) {

	grouterTicket := g.Group("/ticket")
	grouterTicketAuth := grouterTicket.Group("", middleware.MiddlewareAuth)
	grouterTicketRdb := grouterTicket.Group("", ctrl.Rdb.MiddlewareGetDataRedis)
	grouterTicket.GET("", ctrl.GetAllTicket)
	grouterTicketRdb.GET("/:ticketId", ctrl.GetOneTicket, ctrl.Rdb.MiddlewareSetDataRedis)
	grouterTicket.GET("/driver/:driverId", ctrl.GetAllTicketOnSpecificDriver)
	grouterTicket.GET("/customer/:customerId", ctrl.GetAllTicketOnSpecificCustomer)
	grouterTicket.GET("/bus/:busId", ctrl.GetAllTicketOnSpecificBus)
	grouterTicket.GET("/agency/:agencyId", ctrl.GetAllTicketOnSpecificAgency)
	grouterTicket.GET("/price", ctrl.GetTotalPriceAllTicket)
	grouterTicketAuth.GET("/agency/:agencyId/price", ctrl.GetTotalPriceTicketFromSpecificAgency)
	grouterTicketAuth.GET("/driver/:driverId/price", ctrl.GetTotalPriceTicketFromSpecificDriver)
	grouterTicketAuth.POST("", ctrl.AddTicket)
	grouterTicketAuth.DELETE("/:ticketId", ctrl.DeleteTicket)
}

func (ctrl *ControllerTicketImplementation) GetAllTicket(ctx *gin.Context) {
	filter := request.TicketFilter{}

	err := ctx.BindQuery(&filter)
	fmt.Println("CHECK FILTER ", filter.Arrived, "err", err)
	if err != nil {

		panic(exception.NewBadRequestError(err.Error()))
	}

	listTicket := ctrl.service.GetAllTicket(ctx, &filter)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: listTicket}
	ctx.JSON(http.StatusOK, &finalResponse)
}
func (ctrl *ControllerTicketImplementation) AddTicket(ctx *gin.Context) {
	ticketRequest := request.Ticket{}
	err := ctx.Bind(&ticketRequest)
	helper.PanicIfError(err)
	ctrl.service.AddTicket(ctx, &ticketRequest)

	finalResponse := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}

	ctx.JSON(http.StatusOK, &finalResponse)

}
func (ctrl *ControllerTicketImplementation) GetOneTicket(ctx *gin.Context) {
	ticketId, isTicketId := ctx.Params.Get("ticketId")

	if !isTicketId {
		panic(exception.NewBadRequestError("ERROR TICKET ID NOT FOUND"))
	}

	ticketIdInt, err := strconv.Atoi(ticketId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR TICKET ID NOT INTEGER"))
	}

	Ticket := ctrl.service.GetOneTicket(ctx, ticketIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.Set("response", finalResponse)
	ctx.JSON(http.StatusOK, &finalResponse)
}
func (ctrl *ControllerTicketImplementation) DeleteTicket(ctx *gin.Context) {
	ticketId, isTicketId := ctx.Params.Get("ticketId")

	if !isTicketId {
		panic(exception.NewBadRequestError("ERROR TICKET ID NOT FOUND"))
	}

	ticketIdInt, err := strconv.Atoi(ticketId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR TICKET ID NOT INTEGER"))
	}

	Ticket := ctrl.service.DeleteTicket(ctx, ticketIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.JSON(http.StatusOK, &finalResponse)
}
func (ctrl *ControllerTicketImplementation) GetAllTicketOnSpecificDriver(ctx *gin.Context) {
	driverId, isDriverId := ctx.Params.Get("driverId")

	if !isDriverId {
		panic(exception.NewBadRequestError("ERROR DRIVER ID NOT FOUND"))
	}

	driverIdInt, err := strconv.Atoi(driverId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR DRIVER ID NOT INTEGER"))
	}

	Ticket := ctrl.service.GetAllTicketOnDriver(ctx, driverIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.JSON(http.StatusOK, &finalResponse)

}
func (ctrl *ControllerTicketImplementation) GetAllTicketOnSpecificCustomer(ctx *gin.Context) {
	customerId, isCustomerId := ctx.Params.Get("customerId")

	if !isCustomerId {
		panic(exception.NewBadRequestError("ERROR CUSTOMER ID NOT FOUND"))
	}

	customerIdInt, err := strconv.Atoi(customerId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR CUSTOMER ID NOT INTEGER"))
	}

	Ticket := ctrl.service.GetAllTicketOnCustomer(ctx, customerIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.JSON(http.StatusOK, &finalResponse)

}
func (ctrl *ControllerTicketImplementation) GetAllTicketOnSpecificAgency(ctx *gin.Context) {
	agencyId, isAgencyId := ctx.Params.Get("agencyId")

	if !isAgencyId {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT FOUND"))
	}

	agencyIdInt, err := strconv.Atoi(agencyId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT INTEGER"))
	}
	Ticket := ctrl.service.GetAllTicketOnAgency(ctx, agencyIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.JSON(http.StatusOK, &finalResponse)

}
func (ctrl *ControllerTicketImplementation) GetAllTicketOnSpecificBus(ctx *gin.Context) {
	busId, isBusId := ctx.Params.Get("busId")

	if !isBusId {
		panic(exception.NewBadRequestError("ERROR BUS ID NOT FOUND"))
	}

	busIdInt, err := strconv.Atoi(busId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR BUS ID NOT INTEGER"))
	}

	Ticket := ctrl.service.GetAllTicketOnBus(ctx, busIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.JSON(http.StatusOK, &finalResponse)

}

func (ctrl *ControllerTicketImplementation) GetTotalPriceAllTicket(ctx *gin.Context) {

	response := ctrl.service.GetTotalPriceAllTicket(ctx)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: response}
	ctx.JSON(http.StatusOK, &finalResponse)
}

func (ctrl *ControllerTicketImplementation) GetTotalPriceTicketFromSpecificAgency(ctx *gin.Context) {
	agencyId, isAgencyId := ctx.Params.Get("agencyId")

	fmt.Println("CHECK agency", agencyId)

	if !isAgencyId {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT FOUND"))
	}

	agencyIdInt, err := strconv.Atoi(agencyId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT INTEGER"))
	}

	response := ctrl.service.GetTotalPriceTicketFromSpecificAgency(ctx, agencyIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: response}
	ctx.JSON(http.StatusOK, &finalResponse)
}
func (ctrl *ControllerTicketImplementation) GetTotalPriceTicketFromSpecificDriver(ctx *gin.Context) {
	driverId, isDriverId := ctx.Params.Get("driverId")

	if !isDriverId {
		panic(exception.NewBadRequestError("ERROR DRIVER ID NOT FOUND"))
	}

	driverIdInt, err := strconv.Atoi(driverId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR DRIVER ID NOT INTEGER"))
	}

	response := ctrl.service.GetTotalPriceTicketFromSpecificDriver(ctx, driverIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: response}
	ctx.JSON(http.StatusOK, &finalResponse)
}
