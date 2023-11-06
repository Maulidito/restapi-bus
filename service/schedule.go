package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	croncustom "restapi-bus/cron_custom"
	"restapi-bus/exception"
	"restapi-bus/helper"
	cronmodel "restapi-bus/models/cron_model"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"time"
)

type ScheduleServiceImplementation struct {
	RepoSchedule entity.ScheduleRepositoryInterface
	RepoAgency   entity.AgencyRepositoryInterface
	RepoDriver   entity.DriverRepositoryInterface
	RepoBus      entity.BusRepositoryInterface
	CronCustom   croncustom.InterfaceCronJob
	Tx           database.TrInterface
}

func NewScheduleService(
	repoSchedule entity.ScheduleRepositoryInterface,
	repoAgency entity.AgencyRepositoryInterface,
	repoDriver entity.DriverRepositoryInterface,
	cronCustom croncustom.InterfaceCronJob,
	tx database.TrInterface,
	repoBus entity.BusRepositoryInterface) entity.ScheduleServiceInterface {

	scheduleService := &ScheduleServiceImplementation{
		RepoSchedule: repoSchedule,
		RepoAgency:   repoAgency,
		RepoDriver:   repoDriver,
		RepoBus:      repoBus,
		CronCustom:   cronCustom,
		Tx:           tx,
	}
	scheduleService.InitAutoSchedule()
	return scheduleService
}

func (service *ScheduleServiceImplementation) GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []response.Schedule {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	listSchedule := service.RepoSchedule.GetAllSchedule(ctx, filter)

	listResponseSchedule := []response.Schedule{}

	for _, val := range listSchedule {

		listResponseSchedule = append(listResponseSchedule, helper.ScheduleEntityToResponse(&val))

	}

	return listResponseSchedule

}

func (service *ScheduleServiceImplementation) GetOneSchedule(ctx context.Context, scheduleId int) response.Schedule {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) AddSchedule(ctx context.Context, schedule *request.Schedule) {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
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
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, &entitySchedule)
	service.RepoSchedule.DeleteSchedule(ctx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) UpdateArrivedSchedule(ctx context.Context, scheduleId int, isArrived bool) response.Schedule {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	entitySchedule := entity.Schedule{ScheduleId: scheduleId}

	service.RepoSchedule.GetOneSchedule(ctx, &entitySchedule)
	entitySchedule.Arrived = isArrived
	service.RepoSchedule.UpdateArrivedSchedule(ctx, &entitySchedule)

	return helper.ScheduleEntityToResponse(&entitySchedule)

}

func (service *ScheduleServiceImplementation) SetAutoSchedule(ctx context.Context, autoSchedule *request.AutoSchedule) {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
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
	startDateTime := time.Now().Local()
	if autoSchedule.StartFrom != "" {
		startDateTime, err = time.ParseInLocation(time.DateOnly, autoSchedule.StartFrom, time.Local)
		helper.PanicIfError(err)

	}
	timeCurrent := time.Date(
		startDateTime.Year(),
		startDateTime.Month(),
		startDateTime.Day(),
		0, 0, 0, 0,
		time.Local)
	var timeEnd time.Time

	timeEnd = timeCurrent.AddDate(0, autoSchedule.AddRangeMonth, autoSchedule.AddRangeDay)

	timeEnd = timeEnd.
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
	startHour = time.Date(timeCurrent.Year(), timeCurrent.Month(), timeCurrent.Day(), startHour.Hour(), startHour.Minute(), startHour.Second(), 0, time.Local)
	EndHour = time.Date(timeCurrent.Year(), timeCurrent.Month(), timeCurrent.Day(), EndHour.Hour(), EndHour.Minute(), EndHour.Second(), 0, time.Local)

	listSchedule, timeCurrent, timeEnd, startHour, EndHour := helper.GenerateAutoSchedule(
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
			ctx = context.Background()
			ctx = service.Tx.BeginTransactionWithContext(ctx)
			defer service.Tx.DoCommitOrRollbackWithContext(ctx)
			timeEnd = timeEnd.AddDate(0, 0, 1)
			startHour = startHour.AddDate(0, 0, 1)
			EndHour = EndHour.AddDate(0, 0, 1)
			timeCurrent = startHour

			listSchedule = []entity.Schedule{}
			listSchedule, timeCurrent, timeEnd, startHour, EndHour = helper.GenerateAutoSchedule(
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
		"0 0 * * *", //every day
		autoSchedule,
		fmt.Sprintf("Everyday - Schedule %s - Time %s", idCronJob, "00:00"),
		true,
	)

}

func (service *ScheduleServiceImplementation) InitAutoSchedule() {
	ctx := context.Background()
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	listCronJob := service.CronCustom.LoadConfigCronJobSchedule()

	for id, val := range listCronJob {
		job := val.Job.(map[string]interface{})
		autoSchedule := &request.AutoSchedule{}
		byteJob, err := json.Marshal(job)

		if err != nil {
			log.Print("something went wrong when decode the json config")
		}
		err = json.Unmarshal(byteJob, autoSchedule)
		if err != nil {
			log.Print("something went wrong when decode the json config")
		}

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
		countBusFirstAgency = len(listBusFirstAgency)

		countBusSecondAgency = len(listBusSecondAgency)

		if countBusFirstAgency == 0 {
			log.Printf("got bus in first agency %d , the total bus cannot 0", countBusSecondAgency)
		}
		if autoSchedule.BothAgency {
			if countBusSecondAgency == 0 {
				log.Printf("got bus in  second agency %d, the total bus cannot 0", countBusSecondAgency)
			}
		}

		startHour, err := time.Parse(time.TimeOnly, autoSchedule.StartHour)
		helper.PanicIfError(err)
		EndHour, err := time.Parse(time.TimeOnly, autoSchedule.EndHour)
		helper.PanicIfError(err)
		EstimateTime, err := time.Parse(time.TimeOnly, autoSchedule.EstimateTime)
		helper.PanicIfError(err)
		startDateTime := time.Now().Local()
		timeCurrent := time.Date(
			startDateTime.Year(),
			startDateTime.Month(),
			startDateTime.Day(),
			0, 0, 0, 0,
			time.Local)
		var timeEnd time.Time

		timeEnd = timeCurrent.AddDate(0, autoSchedule.AddRangeMonth, autoSchedule.AddRangeDay)

		timeEnd = timeEnd.
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
		startHour = time.Date(timeCurrent.Year(), timeCurrent.Month(), timeCurrent.Day(), startHour.Hour(), startHour.Minute(), startHour.Second(), 0, time.Local)
		EndHour = time.Date(timeCurrent.Year(), timeCurrent.Month(), timeCurrent.Day(), EndHour.Hour(), EndHour.Minute(), EndHour.Second(), 0, time.Local)
		listSchedule := []entity.Schedule{}

		service.CronCustom.SetCronJob(
			id,
			func() {
				ctx = context.Background()
				ctx = service.Tx.BeginTransactionWithContext(ctx)
				defer service.Tx.DoCommitOrRollbackWithContext(ctx)
				timeEnd = timeEnd.AddDate(0, 0, 1)
				startHour = startHour.AddDate(0, 0, 1)
				EndHour = EndHour.AddDate(0, 0, 1)
				timeCurrent = startHour

				listSchedule, timeCurrent, timeEnd, startHour, EndHour = helper.GenerateAutoSchedule(
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
			val.Spec, //every day
			autoSchedule,
			val.Desc,
			false,
		)
	}

}

func (service *ScheduleServiceImplementation) GetAutoSchedule(ctx context.Context) []cronmodel.ResponseCronJob {
	listAllCronJob, err := service.CronCustom.GetAllCronJob()

	helper.PanicIfError(err)
	listResponseCronJob := []cronmodel.ResponseCronJob{}
	for id, val := range listAllCronJob {
		listResponseCronJob = append(listResponseCronJob, cronmodel.ResponseCronJob{Id: id, Spec: val.Spec, Desc: val.Description})
	}

	return listResponseCronJob
}

func (service *ScheduleServiceImplementation) DeleteAutoSchedule(ctx context.Context, id string) {
	err := service.CronCustom.DeleteOneCronJob(id)
	helper.PanicIfError(err)
}
