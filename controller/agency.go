package controller

import (
	"net/http"
	"restapi-bus/constant"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/middleware"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/web"

	"strconv"

	"github.com/gin-gonic/gin"
)

type AgencyControllerInterface interface {
	GetAllAgency(ctx *gin.Context)
	RegisterAgency(ctx *gin.Context)
	GetOneAgency(ctx *gin.Context)
	DeleteOneAgency(ctx *gin.Context)
	LoginAgency(ctx *gin.Context)
	RouterMount(g gin.IRouter)
}

type AgencyControllerImplementation struct {
	service entity.AgencyServiceInterface
	Rdb     *middleware.RedisClientDb
}

func NewAgencyController(service entity.AgencyServiceInterface, rdb *middleware.RedisClientDb) AgencyControllerInterface {
	return &AgencyControllerImplementation{service: service, Rdb: rdb}
}

func (ctrl *AgencyControllerImplementation) RouterMount(g gin.IRouter) {

	grouterAgency := g.Group("/agency")
	grouterAgencyAuth := grouterAgency.Group("", middleware.MiddlewareAuth)
	grouterAgencyRdb := grouterAgency.Group("", ctrl.Rdb.MiddlewareGetDataRedis)
	grouterAgency.POST("/login", ctrl.LoginAgency)
	grouterAgency.GET("", ctrl.GetAllAgency)
	grouterAgency.POST("", ctrl.RegisterAgency)
	grouterAgencyRdb.GET("/:agencyId", ctrl.GetOneAgency, ctrl.Rdb.MiddlewareSetDataRedis)
	grouterAgencyAuth.DELETE("/:agencyId", ctrl.DeleteOneAgency)
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
func (ctrl *AgencyControllerImplementation) RegisterAgency(ctx *gin.Context) {

	agencyRequest := request.Agency{}
	agencyAuth := request.AgencyAuth{}
	err := ctx.ShouldBind(&agencyRequest)
	helper.PanicIfError(err)
	err = ctx.ShouldBind(&agencyAuth)
	helper.PanicIfError(err)
	agencyRequest.Auth = &agencyAuth
	ctrl.service.RegisterAgency(ctx, &agencyRequest)
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

	ctx.Set("response", finalResponse)

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

func (ctrl *AgencyControllerImplementation) LoginAgency(ctx *gin.Context) {

	agencyAuth := request.AgencyAuth{}
	err := ctx.Bind(&agencyAuth)

	helper.PanicIfError(err)
	token, maxAge, _ := ctrl.service.LoginAgency(ctx, &agencyAuth)
	ctx.SetCookie(constant.X_API_KEY, token, maxAge, "/", "localhost", true, true)
	webToken := web.Token{Token: token}

	ctx.JSON(http.StatusOK, web.WebResponseToken{Code: http.StatusOK, Status: "OK", Data: &webToken})
}
