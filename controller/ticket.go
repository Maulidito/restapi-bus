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
}

type ControllerTicketImplementation struct {
	service service.TicketServiceInterface
}

func NewTicketController(serv service.TicketServiceInterface) ControllerTicketInterface {
	return &ControllerTicketImplementation{service: serv}
}

func (ctrl *ControllerTicketImplementation) GetAllTicket(ctx *gin.Context) {
	listTicket := ctrl.service.GetAllTicket(ctx)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: listTicket}
	ctx.JSON(http.StatusOK, &finalResponse)
}
func (ctrl *ControllerTicketImplementation) AddTicket(ctx *gin.Context) {
	ticketRequest := request.Ticket{}
	err := ctx.ShouldBind(&ticketRequest)
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
	driverId, isDriverId := ctx.Params.Get("busId")

	if !isDriverId {
		panic(exception.NewBadRequestError("ERROR BUS ID NOT FOUND"))
	}

	driverIdInt, err := strconv.Atoi(driverId)

	if err != nil {
		panic(exception.NewBadRequestError("ERROR BUS ID NOT INTEGER"))
	}

	Ticket := ctrl.service.GetAllTicketOnDriver(ctx, driverIdInt)
	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: Ticket}
	ctx.JSON(http.StatusOK, &finalResponse)

}
