package service

import (
	"context"
	"fmt"
	croncustom "restapi-bus/cron_custom"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"sync"
	"time"
)

type ScheduleServiceImplementation struct {
	RepoSchedule entity.ScheduleRepositoryInterface
	RepoAgency   entity.AgencyRepositoryInterface
	RepoDriver   entity.DriverRepositoryInterface
	RepoBus      entity.BusRepositoryInterface
	CronCustom   croncustom.InterfaceCronJob
}

func NewScheduleService(
	repoSchedule entity.ScheduleRepositoryInterface,
	repoAgency entity.AgencyRepositoryInterface,
	repoDriver entity.DriverRepositoryInterface,
	cronCustom croncustom.InterfaceCronJob,
	repoBus entity.BusRepositoryInterface) entity.ScheduleServiceInterface {
	return &ScheduleServiceImplementation{
		RepoSchedule: repoSchedule,
		RepoAgency:   repoAgency,
		RepoDriver:   repoDriver,
		RepoBus:      repoBus,
		CronCustom:   cronCustom,
	}
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

func (service *ScheduleServiceImplementation) AutoSchedule(ctx context.Context, autoSchedule *request.AutoSchedule) {
	countBusFirstAgency := 0
	countBusSecondAgency := 0
	listDriverFirstAgency := service.RepoDriver.GetAllDriverOnSpecificAgency(ctx, autoSchedule.FirstAgencyId)
	listBusFirstAgency := service.RepoBus.GetAllBusSpecificAgency(ctx, autoSchedule.FirstAgencyId)
	listBusSecondAgency := service.RepoBus.GetAllBusSpecificAgency(ctx, autoSchedule.SecondAgencyId)
	listDriverSecondAgency := service.RepoDriver.GetAllDriverOnSpecificAgency(ctx, autoSchedule.SecondAgencyId)
	entityAgency1 := entity.Agency{AgencyId: autoSchedule.FirstAgencyId}
	entityAgency2 := entity.Agency{AgencyId: autoSchedule.SecondAgencyId}
	service.RepoAgency.GetOneAgency(ctx, &entityAgency1)
	service.RepoAgency.GetOneAgency(ctx, &entityAgency2)
	idCronJob := fmt.Sprintf("%s-%s", entityAgency1.Place, entityAgency2.Place)
	if service.CronCustom.IsCronJobRunning(idCronJob) {
		panic(fmt.Errorf("this route have a cron job running with id %s", idCronJob))
	}

	countBusFirstAgency = len(listBusFirstAgency)

	countBusSecondAgency = len(listBusSecondAgency)

	if countBusFirstAgency == 0 {
		helper.PanicIfError(exception.NewBadRequestError(fmt.Sprintf("got bus in first agency %d , the total bus cannot 0", countBusSecondAgency)))
	}
	if autoSchedule.BothAgency {
		if countBusSecondAgency == 0 {
			helper.PanicIfError(exception.NewBadRequestError(fmt.Sprintf("got bus in  second agency %d, the total bus cannot 0", countBusSecondAgency)))
		}
	}

	startHour, err := time.Parse(time.TimeOnly, autoSchedule.StartHour)
	helper.PanicIfError(err)
	EndHour, err := time.Parse(time.TimeOnly, autoSchedule.EndHour)
	helper.PanicIfError(err)
	EstimateTime, err := time.Parse(time.TimeOnly, autoSchedule.EstimateTime)
	helper.PanicIfError(err)

	timeCurrent := time.Date(
		time.Now().Local().Year(),
		time.Now().Local().Month(),
		time.Now().Local().Day(),
		0, 0, 0, 0,
		time.UTC)
	timeEnd := timeCurrent.AddDate(0, autoSchedule.RangeSchedule, 0).
		Add(
			(time.Duration(EndHour.Hour()) * time.Hour) +
				(time.Duration(EndHour.Minute()) * time.Minute) +
				(time.Duration(EndHour.Second()) * time.Second),
		)
	timeCurrent = timeCurrent.
		Add(
			(time.Duration(startHour.Hour()) * time.Hour) +
				(time.Duration(startHour.Minute()) * time.Minute) +
				(time.Duration(startHour.Second()) * time.Second),
		)
	startHour = startHour.AddDate(timeCurrent.Year(), int(timeCurrent.Month())-int(startHour.Month()), timeCurrent.Day()-startHour.Day())
	EndHour = EndHour.AddDate(timeCurrent.Year(), int(timeCurrent.Month())-int(EndHour.Month()), timeCurrent.Day()-EndHour.Day())

	listSchedule, timeCurrent, timeEnd, startHour, EndHour := generateAutoSchedule(
		timeCurrent,
		timeEnd, EstimateTime,
		startHour, EndHour,
		autoSchedule,
		listBusFirstAgency,
		listBusSecondAgency,
		listDriverFirstAgency,
		listDriverSecondAgency,
	)
	for i := 0; i < len(listSchedule); i++ {

		service.RepoSchedule.AddSchedule(ctx, &listSchedule[i])

	}

	service.CronCustom.SetCronJob(
		fmt.Sprintf("%s-%s", entityAgency1.Place, entityAgency2.Place),
		func() {

			timeEnd = timeEnd.AddDate(0, 0, 1)
			startHour = startHour.AddDate(0, 0, 1)
			EndHour = EndHour.AddDate(0, 0, 1)
			timeCurrent = startHour

			listSchedule = []entity.Schedule{}
			listSchedule, timeCurrent, timeEnd, startHour, EndHour = generateAutoSchedule(
				timeCurrent,
				timeEnd, EstimateTime,
				startHour, EndHour,
				autoSchedule,
				listBusFirstAgency,
				listBusSecondAgency,
				listDriverFirstAgency,
				listDriverSecondAgency,
			)
			for i := 0; i < len(listSchedule); i++ {

				service.RepoSchedule.AddSchedule(ctx, &listSchedule[i])

			}

		},
		"0 * * * *", //every day
	)

}

func generateAutoSchedule(
	timeCurrent time.Time,
	timeEnd time.Time,
	EstimateTime time.Time,
	startHour time.Time,
	EndHour time.Time,
	autoSchedule *request.AutoSchedule,
	listBusFirstAgency []entity.Bus,
	listBusSecondAgency []entity.Bus,
	listDriverFirstAgency []entity.Driver,
	listDriverSecondAgency []entity.Driver,
) ([]entity.Schedule, time.Time, time.Time, time.Time, time.Time) {
	counter := 0
	listSchedule := []entity.Schedule{}
	timeNowWithoutTZ := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.UTC)
	for timeCurrent.Before(timeEnd) {

		if timeCurrent.Equal(EndHour) || timeCurrent.After(EndHour) || timeCurrent.Before(startHour) {
			startHour = startHour.AddDate(0, 0, 1)
			EndHour = EndHour.AddDate(0, 0, 1)
			timeCurrent = startHour
		}

		if timeCurrent.After(timeNowWithoutTZ) {
			listSchedule = append(listSchedule, entity.Schedule{
				FromAgencyId: autoSchedule.FirstAgencyId,
				ToAgencyId:   autoSchedule.SecondAgencyId,
				Date:         fmt.Sprintf("%s %s", timeCurrent.Format(time.DateOnly), timeCurrent.Format(time.TimeOnly)),
				Price:        autoSchedule.Price,
				Arrived:      false,
				BusId:        listBusFirstAgency[counter%len(listBusFirstAgency)].BusId,
				DriverId:     listDriverFirstAgency[counter%len(listDriverFirstAgency)].DriverId,
			})
			if autoSchedule.BothAgency {
				listSchedule = append(listSchedule, entity.Schedule{
					FromAgencyId: autoSchedule.SecondAgencyId,
					ToAgencyId:   autoSchedule.FirstAgencyId,
					Date:         fmt.Sprintf("%s %s", timeCurrent.Format(time.DateOnly), timeCurrent.Format(time.TimeOnly)),
					Price:        autoSchedule.Price,
					Arrived:      false,
					BusId:        listBusSecondAgency[counter%len(listBusSecondAgency)].BusId,
					DriverId:     listDriverSecondAgency[counter%len(listDriverSecondAgency)].DriverId,
				})
			}
		}
		timeCurrent = timeCurrent.
			Add((time.Duration(EstimateTime.Hour())) * time.Hour).
			Add((time.Duration(EstimateTime.Minute())) * time.Minute)

		counter++
	}

	return listSchedule, timeCurrent, timeEnd, startHour, EndHour

}
