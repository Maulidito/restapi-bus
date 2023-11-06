package cronmodel

import "time"

type CronJobSchedule struct {
	Spec      string
	Job       any
	Desc      string
	CreatedAt time.Time
}

type ResponseCronJob struct {
	Id   string
	Spec string
	Desc string
}

type CronJobConfig map[string]CronJobSchedule
