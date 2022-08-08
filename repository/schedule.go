package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type ScheduleRepositoryInterface interface {
	GetAllSchedule(ctx context.Context, tx *sql.Tx, filter string) []entity.Schedule
	GetOneSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule)
	AddSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule)
	DeleteSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule)
	UpdateArrivedSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule)
}

type ScheduleRepositoryImplementation struct {
}

func NewScheduleRepository() ScheduleRepositoryInterface {
	return &ScheduleRepositoryImplementation{}
}

func (repo *ScheduleRepositoryImplementation) GetAllSchedule(ctx context.Context, tx *sql.Tx, filter string) []entity.Schedule {

	defer func() {
		err := recover()
		if err != nil {
			panic(err)
		}
	}()
	row, err := tx.QueryContext(ctx, "SELECT * FROM schedule "+filter)

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

func (repo *ScheduleRepositoryImplementation) GetOneSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule) {

	err := tx.QueryRowContext(ctx, "SELECT * FROM schedule WHERE schedule_id = ?", schedule.ScheduleId).
		Scan(&schedule.ScheduleId, &schedule.FromAgencyId, &schedule.ToAgencyId, &schedule.BusId,
			&schedule.DriverId, &schedule.Price, &schedule.Date, &schedule.Arrived)

	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("SCHEDULE ID %d NOT FOUND", schedule.ScheduleId)))
	}
}

func (repo *ScheduleRepositoryImplementation) AddSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule) {

	res, err := tx.ExecContext(ctx, "INSERT INTO schedule(from_agency_id,to_agency_id,bus_id,driver_id,price,date,arrived) VALUES (?,?,?,?,?,?,?)", schedule.FromAgencyId, schedule.ToAgencyId, schedule.BusId, schedule.DriverId, schedule.Price, schedule.Date, schedule.Arrived)
	helper.PanicIfError(err)
	scheduleId, err := res.LastInsertId()
	helper.PanicIfError(err)
	schedule.ScheduleId = int(scheduleId)
}

func (repo *ScheduleRepositoryImplementation) DeleteSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule) {

	_, err := tx.ExecContext(ctx, "DELETE FROM schedule WHERE schedule_id = ?", schedule.ScheduleId)
	helper.PanicIfError(err)

}

func (repo *ScheduleRepositoryImplementation) UpdateArrivedSchedule(ctx context.Context, tx *sql.Tx, schedule *entity.Schedule) {

	_, err := tx.ExecContext(ctx, "UPDATE schedule SET arrived = ? WHERE schedule_id = ?", schedule.Arrived, schedule.ScheduleId)
	helper.PanicIfError(err)
}
