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
	GetAllAgency(ctx context.Context) []response.Agency
	AddAgency(ctx context.Context, agency *request.Agency) error
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

func (service *AgencyServiceImplemtation) GetAllAgency(ctx context.Context) []response.Agency {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	listAgency := service.Repo.GetAllAgency(ctx, tx)
	listAgencyResponse := []response.Agency{}

	for _, agency := range listAgency {
		listAgencyResponse = append(listAgencyResponse, helper.AgencyEntityToResponse(&agency))

	}

	return listAgencyResponse

}
func (service *AgencyServiceImplemtation) AddAgency(ctx context.Context, agency *request.Agency) error {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	agencyEntity := helper.AgencyRequestToEntity(agency)
	err = service.Repo.AddAgency(ctx, tx, &agencyEntity)
	return err
}
func (service *AgencyServiceImplemtation) GetOneAgency(ctx context.Context, id int) response.Agency {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.GetOneAgency(ctx, tx, &agencyEntity)
	return helper.AgencyEntityToResponse(&agencyEntity)

}
func (service *AgencyServiceImplemtation) DeleteOneAgency(ctx context.Context, id int) response.Agency {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	agencyEntity := entity.Agency{AgencyId: id}
	service.Repo.DeleteOneAgency(ctx, tx, &agencyEntity)

	return helper.AgencyEntityToResponse(&agencyEntity)
}
