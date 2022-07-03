package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
	"sync"
)

type BusServiceInterface interface {
	GetAllBus(ctx context.Context) []response.Bus
	AddBus(ctx context.Context, bus *request.Bus) error
	GetOneBusSpecificAgency(ctx context.Context, idAgency int, idBus int) response.Bus
	DeleteOneBus(ctx context.Context, idAgency int, idBus int) response.Bus
	GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) ([]response.Bus, response.Agency)
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
	helper.PanicIfError(err)

	listBus := service.RepoBus.GetAllBus(ctx, tx)
	listBusResponse := []response.Bus{}

	for _, bus := range listBus {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&bus))

	}

	return listBusResponse

}
func (service *BusServiceImplemtation) AddBus(ctx context.Context, bus *request.Bus) error {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	busEntity := helper.BusRequestToEntity(bus)
	err = service.RepoBus.AddBus(ctx, tx, &busEntity)
	return err
}

func (service *BusServiceImplemtation) GetOneBusSpecificAgency(ctx context.Context, idAgency int, idBus int) response.Bus {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	agencyEntity := service.RepoAgency.GetOneAgency(ctx, tx, idAgency)

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

func (service *BusServiceImplemtation) GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) (responseBus []response.Bus, responseAgency response.Agency) {
	tx, err := service.Db.Begin()
	waitGroup := sync.WaitGroup{}
	helper.PanicIfError(err)
	var busEntity []entity.Bus
	var agencyEntity entity.Agency

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		agencyEntity = service.RepoAgency.GetOneAgency(ctx, tx, idAgency)
		busEntity = service.RepoBus.GetAllBusSpecificAgency(ctx, tx, idAgency)
	}()
	waitGroup.Wait()

	for _, val := range busEntity {
		responseBus = append(responseBus, helper.BusEntityToResponse(&val))
	}

	responseAgency = helper.AgencyEntityToResponse(&agencyEntity)

	return responseBus, responseAgency

}
func (service *BusServiceImplemtation) DeleteOneBus(ctx context.Context, idAgency int, idBus int) response.Bus {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	busEntity := entity.Bus{
		BusId:    idBus,
		AgencyId: idAgency,
	}

	agencyEntity := service.RepoAgency.GetOneAgency(ctx, tx, busEntity.AgencyId)

	if agencyEntity.Name == "" {
		helper.PanicIfError(errors.New(fmt.Sprintf("Agency id %d , Not Found", busEntity.AgencyId)))
	}

	service.RepoBus.DeleteOneBus(ctx, tx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)
}
