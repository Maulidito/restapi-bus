package service

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type BusServiceImplemtation struct {
	RepoBus    entity.BusRepositoryInterface
	RepoAgency entity.AgencyRepositoryInterface
}

func NewBusService(repoBus entity.BusRepositoryInterface, repoAgency entity.AgencyRepositoryInterface) entity.BusServiceInterface {
	return &BusServiceImplemtation{RepoBus: repoBus, RepoAgency: repoAgency}
}

func (service *BusServiceImplemtation) GetAllBus(ctx context.Context, filter *request.BusFilter) []response.Bus {

	listBus := service.RepoBus.GetAllBus(ctx, filter)
	listBusResponse := []response.Bus{}

	for _, bus := range listBus {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&bus))

	}

	return listBusResponse

}
func (service *BusServiceImplemtation) AddBus(ctx context.Context, bus *request.Bus) {

	service.RepoAgency.GetOneAgency(ctx, &entity.Agency{AgencyId: bus.AgencyId})
	busEntity := helper.BusRequestToEntity(bus)
	service.RepoBus.AddBus(ctx, &busEntity)

}

func (service *BusServiceImplemtation) GetOneBusSpecificAgency(ctx context.Context, idBus int) response.Bus {

	busEntity := entity.Bus{
		BusId: idBus,
	}
	service.RepoBus.GetOneBus(ctx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)

}

func (service *BusServiceImplemtation) GetAllBusOnSpecificAgency(ctx context.Context, idAgency int) []response.Bus {

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

		service.RepoAgency.GetOneAgency(ctx, &agencyEntity)
		busEntity = service.RepoBus.GetAllBusSpecificAgency(ctx, idAgency)
	}()

	helper.PanicIfError(<-chanErr)

	for _, val := range busEntity {
		listBusResponse = append(listBusResponse, helper.BusEntityToResponse(&val))
	}

	return listBusResponse

}
func (service *BusServiceImplemtation) DeleteOneBus(ctx context.Context, idBus int) response.Bus {

	busEntity := entity.Bus{
		BusId: idBus,
	}
	service.RepoBus.GetOneBus(ctx, &busEntity)
	service.RepoBus.DeleteOneBus(ctx, &busEntity)

	return helper.BusEntityToResponse(&busEntity)
}
