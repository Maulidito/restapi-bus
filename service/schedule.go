package service

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"sync"
)

type ScheduleServiceImplementation struct {
	RepoSchedule entity.ScheduleRepositoryInterface
	RepoAgency   entity.AgencyRepositoryInterface
	RepoDriver   entity.DriverRepositoryInterface
	RepoBus      entity.BusRepositoryInterface
}

func NewScheduleService(repoSchedule entity.ScheduleRepositoryInterface, repoAgency entity.AgencyRepositoryInterface,
	repoDriver entity.DriverRepositoryInterface,
	repoBus entity.BusRepositoryInterface) entity.ScheduleServiceInterface {
	return &ScheduleServiceImplementation{RepoSchedule: repoSchedule, RepoAgency: repoAgency, RepoDriver: repoDriver, RepoBus: repoBus}
}

func (service *ScheduleServiceImplementation) GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []response.Schedule {

	listSchedule := service.RepoSchedule.GetAllSchedule(ctx, filter)

	listResponseSchedule := []response.Schedule{}

	for _, val := range listSchedule {

		listResponseSchedule = append(listResponseSchedule, helper.ScheduleEntityToResponse(&val))

	}

	return listResponseSchedule

}

func (service *ScheduleServiceImplementation) GetOneSchedule(ctx context.Context, scheduleId int) response.Schedule {

	entitySchedule := entity.Schedule{ScheduleId: scheduleId}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.RepoSchedule.GetOneSchedule(ctx, &entitySchedule)
	}()
	wg.Wait()

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) AddSchedule(ctx context.Context, schedule *request.Schedule) {

	chanErr := make(chan error, 1)

	go func() {
		defer func() {
			defer close(chanErr)
			errRecover := recover()
			if errRecover != nil {
				chanErr <- errRecover.(error)
			}

		}()

		service.RepoAgency.GetOneAgency(ctx, &entity.Agency{AgencyId: schedule.FromAgencyId})

		service.RepoAgency.GetOneAgency(ctx, &entity.Agency{AgencyId: schedule.ToAgencyId})

		service.RepoBus.GetOneBus(ctx, &entity.Bus{BusId: schedule.BusId})

		service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &entity.Driver{DriverId: schedule.DriverId})

	}()

	if err := <-chanErr; err != nil {

		helper.PanicIfError(err)
	}

	entitySchedule := helper.ScheduleRequestToEntity(schedule)
	service.RepoSchedule.AddSchedule(ctx, &entitySchedule)

}

func (service *ScheduleServiceImplementation) DeleteSchedule(ctx context.Context, scheduleId int) response.Schedule {

	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, &entitySchedule)
	service.RepoSchedule.DeleteSchedule(ctx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) UpdateArrivedSchedule(ctx context.Context, scheduleId int, isArrived bool) response.Schedule {

	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, &entitySchedule)
	entitySchedule.Arrived = isArrived
	service.RepoSchedule.UpdateArrivedSchedule(ctx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}
