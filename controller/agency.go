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

type AgencyControllerInterface interface {
	GetAllAgency(ctx *gin.Context)
	AddAgency(ctx *gin.Context)
	GetOneAgency(ctx *gin.Context)
	DeleteOneAgency(ctx *gin.Context)
}

type AgencyControllerImplementation struct {
	service service.AgencyServiceInterface
}

func NewAgencyController(service service.AgencyServiceInterface) AgencyControllerInterface {
	return &AgencyControllerImplementation{service: service}
}

func (ctrl *AgencyControllerImplementation) GetAllAgency(ctx *gin.Context) {

	filter := request.AgencyFilter{}

	err := ctx.Bind(&filter)

	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	agencyResponse := ctrl.service.GetAllAgency(ctx, &filter)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: agencyResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
func (ctrl *AgencyControllerImplementation) AddAgency(ctx *gin.Context) {

	agencyRequest := request.Agency{}
	err := ctx.ShouldBind(&agencyRequest)
	helper.PanicIfError(err)
	ctrl.service.AddAgency(ctx, &agencyRequest)
	helper.PanicIfError(err)

	finalResponse := web.WebResponseNoData{Code: http.StatusOK, Status: "OK"}
	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *AgencyControllerImplementation) GetOneAgency(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("agencyId")

	if !idBool {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT FOUND"))
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT INTEGER"))
	}
	agencyResponse := ctrl.service.GetOneAgency(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: agencyResponse}

	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *AgencyControllerImplementation) DeleteOneAgency(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("agencyId")

	if !idBool {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT FOUND"))

	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		panic(exception.NewBadRequestError("ERROR AGENCY ID NOT INTEGER"))
	}

	agencyResponse := ctrl.service.DeleteOneAgency(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: agencyResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
