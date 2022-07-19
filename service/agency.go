package service

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type AgencyServiceInterface interface {
	GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []response.Agency
	AddAgency(ctx context.Context, agency *request.Agency)
	GetOneAgency(ctx context.Context, id int) response.Agency
	DeleteOneAgency(ctx context.Context, id int) response.Agency
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
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	listAgency := service.Repo.GetAllAgency(ctx, tx, helper.RequestFilterAgencyToString(filter))
	listAgencyResponse := []response.Agency{}

	for _, agency := range listAgency {
		listAgencyResponse = append(listAgencyResponse, helper.AgencyEntityToResponse(&agency))

	}

	return listAgencyResponse

}
func (service *AgencyServiceImplemtation) AddAgency(ctx context.Context, agency *request.Agency) {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)
	agencyEntity := helper.AgencyRequestToEntity(agency)
	service.Repo.AddAgency(ctx, tx, &agencyEntity)

}
func (service *AgencyServiceImplemtation) GetOneAgency(ctx context.Context, id int) response.Agency {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, tx, &agencyEntity)
	return helper.AgencyEntityToResponse(&agencyEntity)

}
func (service *AgencyServiceImplemtation) DeleteOneAgency(ctx context.Context, id int) response.Agency {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, tx, &agencyEntity)
	service.Repo.DeleteOneAgency(ctx, tx, &agencyEntity)

	return helper.AgencyEntityToResponse(&agencyEntity)
}
