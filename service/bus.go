package service

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type BusServiceInterface interface {
	GetAllBus(ctx context.Context) []response.Bus
	AddBus(ctx context.Context, bus *request.Bus) error
	GetOneBus(ctx context.Context, id int) response.Bus
	DeleteOneBus(ctx context.Context, id int) response.Bus
}

type BusServiceImplemtation struct {
	Db   *sql.DB
	Repo repository.BusRepositoryInterface
}

func NewBusService(db *sql.DB, repo repository.BusRepositoryInterface) BusServiceInterface {
	return &BusServiceImplemtation{Db: db, Repo: repo}
}

func (service *BusServiceImplemtation) GetAllBus(ctx context.Context) []response.Bus {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	listBus := service.Repo.GetAllBus(ctx, tx)
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
	err = service.Repo.AddBus(ctx, tx, &busEntity)
	return err
}
func (service *BusServiceImplemtation) GetOneBus(ctx context.Context, id int) response.Bus {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	busEntity := service.Repo.GetOneBus(ctx, tx, id)

	return helper.BusEntityToResponse(&busEntity)

}
func (service *BusServiceImplemtation) DeleteOneBus(ctx context.Context, id int) response.Bus {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	busEntity := service.Repo.DeleteOneBus(ctx, tx, id)

	return helper.BusEntityToResponse(&busEntity)
}
