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

func (service *BusServiceImplemtation) GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) (responseBus []response.Bus, responseAgency response.Agency) {
	tx, err := service.Db.Begin()

	helper.PanicIfError(err)
	var busEntity []entity.Bus

	chanErr := make(chan string, 1)

	var agencyEntity entity.Agency = entity.Agency{AgencyId: idAgency}

	go func() {
		defer func() {

			tempRecover := recover()

			if tempRecover != nil {
				chanErr <- tempRecover.(string)
			}

			close(chanErr)

		}()

		service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)
		busEntity = service.RepoBus.GetAllBusSpecificAgency(ctx, tx, idAgency)
	}()
	helper.PanicIfErrorString(<-chanErr)

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

	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)

	if agencyEntity.Name == "" {
		helper.PanicIfError(fmt.Errorf("agency id %d , not found", busEntity.AgencyId))
	}

	service.RepoBus.DeleteOneBus(ctx, tx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)
}
