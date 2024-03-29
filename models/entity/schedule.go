package entity

import (
	"context"
	cronmodel "restapi-bus/models/cron_model"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type Schedule struct {
	ScheduleId   int
	FromAgencyId int
	ToAgencyId   int
	BusId        int
	DriverId     int
	Price        int
	Date         string
	Arrived      bool `default:"false"`
}

type ScheduleServiceInterface interface {
	GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []response.Schedule
	GetOneSchedule(ctx context.Context, scheduleId int) response.Schedule
	AddSchedule(ctx context.Context, schedule *request.Schedule)
	DeleteSchedule(ctx context.Context, scheduleId int) response.Schedule
	UpdateArrivedSchedule(ctx context.Context, scheduleId int, isArrived bool) response.Schedule
	SetAutoSchedule(ctx context.Context, autoSchedule *request.AutoSchedule)
	GetAutoSchedule(ctx context.Context) []cronmodel.ResponseCronJob
	DeleteAutoSchedule(ctx context.Context, id string)
}

type ScheduleRepositoryInterface interface {
	GetAllSchedule(ctx context.Context, filter *request.ScheduleFilter) []Schedule
	GetOneSchedule(ctx context.Context, schedule *Schedule)
	GetOneDetailSchedule(ctx context.Context, schedule *response.DetailSchedule)
	AddSchedule(ctx context.Context, schedule *Schedule)
	DeleteSchedule(ctx context.Context, schedule *Schedule)
	UpdateArrivedSchedule(ctx context.Context, schedule *Schedule)
}
