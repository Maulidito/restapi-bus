package service

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
	"sync"
)

type ScheduleServiceInterface interface {
	GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []response.Schedule
	GetOneSchedule(ctx context.Context, scheduleId int) response.Schedule
	AddSchedule(ctx context.Context, schedule *request.Schedule)
	DeleteSchedule(ctx context.Context, scheduleId int) response.Schedule
	UpdateArrivedSchedule(ctx context.Context, scheduleId int, isArrived bool) response.Schedule
}

type ScheduleServiceImplementation struct {
	RepoSchedule repository.ScheduleRepositoryInterface
	RepoAgency   repository.AgencyRepositoryInterface
	RepoDriver   repository.DriverRepositoryInterface
	RepoBus      repository.BusRepositoryInterface
	Db           *sql.DB
}

func NewScheduleService(repoSchedule repository.ScheduleRepositoryInterface, repoAgency repository.AgencyRepositoryInterface,
	repoDriver repository.DriverRepositoryInterface,
	repoBus repository.BusRepositoryInterface, db *sql.DB) ScheduleServiceInterface {
	return &ScheduleServiceImplementation{RepoSchedule: repoSchedule, RepoAgency: repoAgency, RepoDriver: repoDriver, RepoBus: repoBus, Db: db}
}

func (service *ScheduleServiceImplementation) GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []response.Schedule {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommitOrRollback(tx)

	listSchedule := service.RepoSchedule.GetAllSchedule(ctx, tx, helper.RequestFilterScheduleToString(filter))

	listResponseSchedule := []response.Schedule{}

	for _, val := range listSchedule {

		listResponseSchedule = append(listResponseSchedule, helper.ScheduleEntityToResponse(&val))

	}

	return listResponseSchedule

}

func (service *ScheduleServiceImplementation) GetOneSchedule(ctx context.Context, scheduleId int) response.Schedule {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommitOrRollback(tx)
	entitySchedule := entity.Schedule{ScheduleId: scheduleId}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.RepoSchedule.GetOneSchedule(ctx, tx, &entitySchedule)
	}()
	wg.Wait()

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) AddSchedule(ctx context.Context, schedule *request.Schedule) {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommitOrRollback(tx)
	chanErr := make(chan error, 1)

	go func() {
		defer func() {
			defer close(chanErr)
			errRecover := recover()
			if errRecover != nil {
				chanErr <- errRecover.(error)
			}

		}()

		service.RepoAgency.GetOneAgency(ctx, tx, &entity.Agency{AgencyId: schedule.FromAgencyId})

		service.RepoAgency.GetOneAgency(ctx, tx, &entity.Agency{AgencyId: schedule.ToAgencyId})

		service.RepoBus.GetOneBus(ctx, tx, &entity.Bus{BusId: schedule.BusId})

		service.RepoDriver.GetOneDriverOnSpecificAgency(tx, ctx, &entity.Driver{DriverId: schedule.DriverId})

	}()

	if err = <-chanErr; err != nil {

		helper.PanicIfError(err)
	}

	entitySchedule := helper.ScheduleRequestToEntity(schedule)
	service.RepoSchedule.AddSchedule(ctx, tx, &entitySchedule)

}

func (service *ScheduleServiceImplementation) DeleteSchedule(ctx context.Context, scheduleId int) response.Schedule {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommitOrRollback(tx)

	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, tx, &entitySchedule)
	service.RepoSchedule.DeleteSchedule(ctx, tx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) UpdateArrivedSchedule(ctx context.Context, scheduleId int, isArrived bool) response.Schedule {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommitOrRollback(tx)

	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, tx, &entitySchedule)
	entitySchedule.Arrived = isArrived
	service.RepoSchedule.UpdateArrivedSchedule(ctx, tx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}
