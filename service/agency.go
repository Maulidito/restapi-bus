package service

import (
	"context"
	"fmt"
	"os"

	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/models/web"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AgencyServiceImplemtation struct {
	Repo entity.AgencyRepositoryInterface
}

func NewAgencyService(repo entity.AgencyRepositoryInterface) entity.AgencyServiceInterface {
	return &AgencyServiceImplemtation{Repo: repo}
}

func (service *AgencyServiceImplemtation) GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []response.Agency {

	listAgency := service.Repo.GetAllAgency(ctx, filter)
	listAgencyResponse := []response.Agency{}

	for _, agency := range listAgency {
		listAgencyResponse = append(listAgencyResponse, helper.AgencyEntityToResponse(&agency))

	}

	return listAgencyResponse

}
func (service *AgencyServiceImplemtation) RegisterAgency(ctx context.Context, agency *request.Agency) {
	var (
		err error
	)
	if service.Repo.IsUsenameAgencyExist(ctx, agency.Auth.Username) {
		panic(exception.NewBadRequestError("email already registered"))
	}
	salt := fmt.Sprint(time.Now().UnixNano())
	agency.Auth.Password = helper.HashPassword(agency.Auth.Password, salt)
	helper.PanicIfError(err)
	agencyEntity := helper.AgencyRequestToEntity(agency)
	agencyEntity.Salt = salt
	service.Repo.RegisterAgency(ctx, &agencyEntity)

}
func (service *AgencyServiceImplemtation) GetOneAgency(ctx context.Context, id int) response.Agency {

	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, &agencyEntity)
	return helper.AgencyEntityToResponse(&agencyEntity)

}
func (service *AgencyServiceImplemtation) DeleteOneAgency(ctx context.Context, id int) response.Agency {
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, &agencyEntity)
	service.Repo.DeleteOneAgency(ctx, &agencyEntity)

	return helper.AgencyEntityToResponse(&agencyEntity)
}

func (service *AgencyServiceImplemtation) LoginAgency(ctx context.Context, agencyAuth *request.AgencyAuth) (string, int, response.Agency) {

	salt := service.Repo.GetSaltAgencyWithUsername(ctx, agencyAuth.Username)

	passEncrypted := helper.HashPassword(agencyAuth.Password, salt)

	agency := entity.Agency{Username: agencyAuth.Username, Password: passEncrypted}

	service.Repo.GetOneAgencyAuth(ctx, &agency)

	claim := web.Claim{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
		Username: agency.Username,
	}

	jwToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := os.Getenv("SECRET_KEY_AUTH")

	token, err := jwToken.SignedString([]byte(secret))

	helper.PanicIfError(err)

	return token, claim.ExpiresAt.Second(), helper.AgencyEntityToResponse(&agency)

}
