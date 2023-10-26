package croncustom

import (
	"fmt"

	cron "github.com/robfig/cron/v3"
)

type InterfaceCronJob interface {
	SetCronJobOnce(idCron string, callback func(), timeFormat string) error
	StopCronJob(idCron string)
}

type CronJob struct {
	ListCron map[string]*cron.Cron
}

var CronJobSingleton *CronJob

func NewCronJob() InterfaceCronJob {

	if CronJobSingleton == nil {
		CronJobSingleton = &CronJob{}
		CronJobSingleton.ListCron = make(map[string]*cron.Cron)
	}

	return CronJobSingleton
}

func (c *CronJob) SetCronJobOnce(idCron string, callback func(), timeFormat string) error {

	cronTicket := cron.New()

	c.ListCron[idCron] = cronTicket

	_, err := cronTicket.AddFunc(fmt.Sprintf("CRON_TZ=Asia/Jakarta %s", timeFormat), func() {
		defer c.closeAndDeleteCron(cronTicket, idCron)
		callback()

	})
	cronTicket.Start()
	return err

}

func (c *CronJob) StopCronJob(idCron string) {
	cronJob := c.ListCron[idCron]
	defer c.closeAndDeleteCron(cronJob, idCron)
}

func (c *CronJob) closeAndDeleteCron(cronJob *cron.Cron, idCron string) {
	cronJob.Stop()
	delete(c.ListCron, idCron)
}
