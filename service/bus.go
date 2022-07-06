package service

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type BusServiceInterface interface {
	GetAllBus(ctx context.Context) []response.Bus
	AddBus(ctx context.Context, bus *request.Bus)
	GetOneBusSpecificAgency(ctx context.Context, idAgency int, idBus int) response.Bus
	DeleteOneBus(ctx context.Context, idAgency int, idBus int) response.Bus
	GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) response.AllBusOnAgency
}

type BusServiceImplemtation struct {
	Db         *sql.DB
	RepoBus    repository.BusRepositoryInterface
	RepoAgency repository.AgencyRepositoryInterface
}

func NewBusService(db *sql.DB, repoBus repository.BusRepositoryInterface, repoAgency repository.AgencyRepositoryInterface) BusServiceInterface {
	return &BusServiceImplemtation{Db: db, RepoBus: repoBus, RepoAgency: repoAgency}
}

func (service *BusServiceImplemtation) GetAllBus(ctx context.Context) []response.Bus {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	listBus := service.RepoBus.GetAllBus(ctx, tx)
	listBusResponse := []response.Bus{}

	for _, bus := range listBus {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&bus))

	}

	return listBusResponse

}
func (service *BusServiceImplemtation) AddBus(ctx context.Context, bus *request.Bus) {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)
	busEntity := helper.BusRequestToEntity(bus)
	service.RepoBus.AddBus(ctx, tx, &busEntity)

}

func (service *BusServiceImplemtation) GetOneBusSpecificAgency(ctx context.Context, idAgency int, idBus int) response.Bus {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)

	if agencyEntity.Name == "" {
		helper.PanicIfError(fmt.Errorf("agency id %d , not found", agencyEntity.AgencyId))
	}
	busEntity := entity.Bus{
		BusId:    idBus,
		AgencyId: idAgency,
	}
	service.RepoBus.GetOneBus(ctx, tx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)

}

func (service *BusServiceImplemtation) GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) response.AllBusOnAgency {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)
	var busEntity []entity.Bus

	chanErr := make(chan error, 1)
	listBusResponse := []response.Bus{}
	var agencyEntity entity.Agency = entity.Agency{AgencyId: idAgency}

	go func() {
		defer func() {

			tempRecover := recover()

			if tempRecover != nil {
				chanErr <- tempRecover.(error)
			}

			close(chanErr)

		}()

		service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)
		busEntity = service.RepoBus.GetAllBusSpecificAgency(ctx, tx, idAgency)
	}()

	helper.PanicIfError(<-chanErr)

	for _, val := range busEntity {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&val))
	}

	agencyResponse := helper.AgencyEntityToResponse(&agencyEntity)

	finalResponse := response.AllBusOnAgency{Agency: &agencyResponse, Bus: &listBusResponse}

	return finalResponse

}
func (service *BusServiceImplemtation) DeleteOneBus(ctx context.Context, idAgency int, idBus int) response.Bus {
	tx, err := service.Db.Begin()
	defer helper.DoCommit(tx)
	helper.PanicIfError(err)

	busEntity := entity.Bus{
		BusId:    idBus,
		AgencyId: idAgency,
	}

	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)

	if agencyEntity.Name == "" {
		helper.PanicIfError(fmt.Errorf("agency id %d , not found", busEntity.AgencyId))
	}

	service.RepoBus.DeleteOneBus(ctx, tx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)
}
