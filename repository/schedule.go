package repository

import (
	"context"
	"fmt"
	"log"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

var scheduleRepositorySingleton *ScheduleRepositoryImplementation

type ScheduleRepositoryImplementation struct {
}

func NewScheduleRepository() entity.ScheduleRepositoryInterface {
	if scheduleRepositorySingleton == nil {
		scheduleRepositorySingleton = &ScheduleRepositoryImplementation{}
	}
	return scheduleRepositorySingleton
}

func (repo *ScheduleRepositoryImplementation) GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []entity.Schedule {
	tx := database.GetTransactionContext(ctx)
	defer func() {
		err := recover()
		if err != nil {
			panic(err)
		}
	}()
	filterString := helper.RequestFilterScheduleToString(filter)
	row, err := tx.QueryContext(ctx, "SELECT * FROM schedule "+filterString)

	helper.PanicIfError(err)
	defer row.Close()

	entitySchedule := entity.Schedule{}
	listEntitySchedule := []entity.Schedule{}

	for row.Next() {
		err := row.Scan(&entitySchedule.ScheduleId, &entitySchedule.FromAgencyId, &entitySchedule.ToAgencyId, &entitySchedule.BusId,
			&entitySchedule.DriverId, &entitySchedule.Price, &entitySchedule.Date, &entitySchedule.Arrived)
		helper.PanicIfError(err)
		listEntitySchedule = append(listEntitySchedule, entitySchedule)

	}

	return listEntitySchedule
}

func (repo *ScheduleRepositoryImplementation) GetOneDetailSchedule(ctx context.Context, schedule *response.DetailSchedule) {
	tx := database.GetTransactionContext(ctx)
	row, err := tx.QueryContext(ctx, "SELECT s.schedule_id,a1.name,a1.place,a2.name,a2.place,b.agency_id,b.number_plate,b.total_seat,d.agency_id,d.name,s.price,s.date,s.arrived FROM schedule as s LEFT JOIN agency as a1 ON s.from_agency_id = a1.agency_id LEFT JOIN agency as a2 ON s.to_agency_id = a2.agency_id LEFT JOIN bus as b ON b.bus_id = s.bus_id LEFT JOIN driver as d ON d.driver_id = s.driver_id WHERE s.schedule_id = ?", schedule.ScheduleId)

	helper.PanicIfError(err)
	defer row.Close()

	if row.Next() {

		err = row.Scan(
			&schedule.ScheduleId,
			&schedule.FromAgency.Name, &schedule.FromAgency.Place, &schedule.ToAgency.Name, &schedule.ToAgency.Place,
			&schedule.Bus.AgencyId, &schedule.Bus.NumberPlate, &schedule.Bus.TotalSeat,
			&schedule.Driver.AgencyId, &schedule.Driver.Name, &schedule.Price, &schedule.Date, &schedule.Arrived)
		log.Print(err)

	}
	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("SCHEDULE ID %d NOT FOUND", schedule.ScheduleId)))
	}
}

func (repo *ScheduleRepositoryImplementation) GetOneSchedule(ctx context.Context, schedule *entity.Schedule) {
	tx := database.GetTransactionContext(ctx)

	err := tx.QueryRowContext(ctx, "SELECT * FROM schedule WHERE schedule_id = ?", schedule.ScheduleId).
		Scan(&schedule.ScheduleId, &schedule.FromAgencyId, &schedule.ToAgencyId, &schedule.BusId,
			&schedule.DriverId, &schedule.Price, &schedule.Date, &schedule.Arrived)

	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("SCHEDULE ID %d NOT FOUND", schedule.ScheduleId)))
	}
}

func (repo *ScheduleRepositoryImplementation) AddSchedule(ctx context.Context, schedule *entity.Schedule) {
	tx := database.GetTransactionContext(ctx)

	res, err := tx.ExecContext(ctx, "INSERT INTO schedule(from_agency_id,to_agency_id,bus_id,driver_id,price,date,arrived) VALUES (?,?,?,?,?,?,?)", schedule.FromAgencyId, schedule.ToAgencyId, schedule.BusId, schedule.DriverId, schedule.Price, schedule.Date, schedule.Arrived)
	helper.PanicIfError(err)
	scheduleId, err := res.LastInsertId()
	helper.PanicIfError(err)
	schedule.ScheduleId = int(scheduleId)

}

func (repo *ScheduleRepositoryImplementation) DeleteSchedule(ctx context.Context, schedule *entity.Schedule) {
	tx := database.GetTransactionContext(ctx)
	_, err := tx.ExecContext(ctx, "DELETE FROM schedule WHERE schedule_id = ?", schedule.ScheduleId)
	helper.PanicIfError(err)

}

func (repo *ScheduleRepositoryImplementation) UpdateArrivedSchedule(ctx context.Context, schedule *entity.Schedule) {
	tx := database.GetTransactionContext(ctx)
	_, err := tx.ExecContext(ctx, "UPDATE schedule SET arrived = ? WHERE schedule_id = ?", schedule.Arrived, schedule.ScheduleId)
	helper.PanicIfError(err)
}
