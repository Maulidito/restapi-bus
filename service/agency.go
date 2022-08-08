package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/models/web"
	"restapi-bus/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AgencyServiceInterface interface {
	GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []response.Agency
	RegisterAgency(ctx context.Context, agency *request.Agency)
	GetOneAgency(ctx context.Context, id int) response.Agency
	DeleteOneAgency(ctx context.Context, id int) response.Agency
	LoginAgency(ctx context.Context, agencyAuth *request.AgencyAuth) (string, int, response.Agency)
}

type AgencyServiceImplemtation struct {
	Db   *sql.DB
	Repo repository.AgencyRepositoryInterface
}

func NewAgencyService(db *sql.DB, repo repository.AgencyRepositoryInterface) AgencyServiceInterface {
	return &AgencyServiceImplemtation{Db: db, Repo: repo}
}

func (service *AgencyServiceImplemtation) GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []response.Agency {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	listAgency := service.Repo.GetAllAgency(ctx, tx, helper.RequestFilterAgencyToString(filter))
	listAgencyResponse := []response.Agency{}

	for _, agency := range listAgency {
		listAgencyResponse = append(listAgencyResponse, helper.AgencyEntityToResponse(&agency))

	}

	return listAgencyResponse

}
func (service *AgencyServiceImplemtation) RegisterAgency(ctx context.Context, agency *request.Agency) {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	if service.Repo.IsUsenameAgencyExist(ctx, tx, agency.Auth.Username) {
		panic(exception.NewBadRequestError("email already registered"))
	}
	salt := fmt.Sprint(time.Now().UnixNano())
	agency.Auth.Password = helper.HashPassword(agency.Auth.Password, salt)
	helper.PanicIfError(err)
	agencyEntity := helper.AgencyRequestToEntity(agency)
	agencyEntity.Salt = salt
	service.Repo.RegisterAgency(ctx, tx, &agencyEntity)

}
func (service *AgencyServiceImplemtation) GetOneAgency(ctx context.Context, id int) response.Agency {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, tx, &agencyEntity)
	return helper.AgencyEntityToResponse(&agencyEntity)

}
func (service *AgencyServiceImplemtation) DeleteOneAgency(ctx context.Context, id int) response.Agency {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, tx, &agencyEntity)
	service.Repo.DeleteOneAgency(ctx, tx, &agencyEntity)

	return helper.AgencyEntityToResponse(&agencyEntity)
}

func (service *AgencyServiceImplemtation) LoginAgency(ctx context.Context, agencyAuth *request.AgencyAuth) (string, int, response.Agency) {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommitOrRollback(tx)

	salt := service.Repo.GetSaltAgencyWithUsername(ctx, tx, agencyAuth.Username)

	passEncrypted := helper.HashPassword(agencyAuth.Password, salt)

	agency := entity.Agency{Username: agencyAuth.Username, Password: passEncrypted}

	service.Repo.GetOneAgencyAuth(ctx, tx, &agency)

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
