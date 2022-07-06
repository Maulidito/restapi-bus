package controller

import (
	"fmt"
	"net/http"
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
	agencyResponse := ctrl.service.GetAllAgency(ctx)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: agencyResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
func (ctrl *AgencyControllerImplementation) AddAgency(ctx *gin.Context) {
	name, _ := ctx.Params.Get("name")
	place, _ := ctx.Params.Get("place")

	agencyRequest := request.Agency{Name: name, Place: place}
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
		helper.PanicIfError(fmt.Errorf("ERROR ID PARAMAETER NOT FOUND"))
	}
	idInt, err := strconv.Atoi(id)
	helper.PanicIfError(err)
	agencyResponse := ctrl.service.GetOneAgency(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: agencyResponse}

	ctx.JSON(http.StatusOK, finalResponse)

}
func (ctrl *AgencyControllerImplementation) DeleteOneAgency(ctx *gin.Context) {
	id, idBool := ctx.Params.Get("agencyId")

	if !idBool {
		helper.PanicIfError(fmt.Errorf("ERROR ID PARAMAETER NOT FOUND"))

	}
	idInt, err := strconv.Atoi(id)
	helper.PanicIfError(err)
	agencyResponse := ctrl.service.DeleteOneAgency(ctx, idInt)

	finalResponse := web.WebResponse{Code: http.StatusOK, Status: "OK", Data: agencyResponse}

	ctx.JSON(http.StatusOK, finalResponse)
}
