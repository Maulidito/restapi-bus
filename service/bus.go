package service

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type BusServiceImplemtation struct {
	RepoBus    entity.BusRepositoryInterface
	RepoAgency entity.AgencyRepositoryInterface
	Tx         database.TrInterface
}

func NewBusService(repoBus entity.BusRepositoryInterface, repoAgency entity.AgencyRepositoryInterface, tx database.TrInterface) entity.BusServiceInterface {
	return &BusServiceImplemtation{RepoBus: repoBus, RepoAgency: repoAgency, Tx: tx}
}

func (service *BusServiceImplemtation) GetAllBus(ctx context.Context, filter *request.BusFilter) []response.Bus {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	listBus := service.RepoBus.GetAllBus(ctx, filter)
	listBusResponse := []response.Bus{}

	for _, bus := range listBus {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&bus))

	}

	return listBusResponse

}
func (service *BusServiceImplemtation) AddBus(ctx context.Context, bus *request.Bus) {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	service.RepoAgency.GetOneAgency(ctx, &entity.Agency{AgencyId: bus.AgencyId})
	busEntity := helper.BusRequestToEntity(bus)
	service.RepoBus.AddBus(ctx, &busEntity)

}

func (service *BusServiceImplemtation) GetOneBusSpecificAgency(ctx context.Context, idBus int) response.Bus {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	busEntity := entity.Bus{
		BusId: idBus,
	}
	service.RepoBus.GetOneBus(ctx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)

}

func (service *BusServiceImplemtation) GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) []response.Bus {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	var busEntity []entity.Bus

	listBusResponse := []response.Bus{}
	var agencyEntity entity.Agency = entity.Agency{AgencyId: idAgency}

	service.RepoAgency.GetOneAgency(ctx, &agencyEntity)
	busEntity = service.RepoBus.GetAllBusSpecificAgency(ctx, idAgency)

	for _, val := range busEntity {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&val))
	}

	return listBusResponse

}
func (service *BusServiceImplemtation) DeleteOneBus(ctx context.Context, idBus int) response.Bus {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	busEntity := entity.Bus{
		BusId: idBus,
	}
	service.RepoBus.GetOneBus(ctx, &busEntity)
	service.RepoBus.DeleteOneBus(ctx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)
}
