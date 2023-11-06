package croncustom

import (
	"encoding/json"
	"fmt"
	"os"
	"restapi-bus/helper"
	cronmodel "restapi-bus/models/cron_model"
	"time"

	cron "github.com/robfig/cron/v3"
)

type InterfaceCronJob interface {
	SetCronJobOnce(idCron string, callback func(), timeFormat string, desc string) error
	SetCronJob(idCron string, callback func(), timeFormat string, structData any, desc string, isWriteConfig bool) error
	GetAllCronJob() (map[string]DetailCronJob, error)
	DeleteOneCronJob(id string) error
	StopCronJob(idCron string)
	IsCronJobRunning(idCron string) bool
	MakeIdCronSchedule(place1 string, place2 string) string
	LoadConfigCronJobSchedule() cronmodel.CronJobConfig
}

type CronJob struct {
	ListCron map[string]DetailCronJob
}
type DetailCronJob struct {
	Spec        string
	Description string
	Cron        *cron.Cron
}

var CronJobSingleton *CronJob

func NewCronJob() InterfaceCronJob {

	if CronJobSingleton == nil {
		CronJobSingleton = &CronJob{}
		CronJobSingleton.ListCron = make(map[string]DetailCronJob)
	}

	return CronJobSingleton
}

func (c *CronJob) GetAllCronJob() (map[string]DetailCronJob, error) {
	listCronJob := c.ListCron

	//c.ListCron
	return listCronJob, nil
}

func (c *CronJob) MakeIdCronSchedule(place1 string, place2 string) string {
	return fmt.Sprintf("%s-%s", place1, place2)
}

func (c *CronJob) SetCronJobOnce(idCron string, callback func(), timeFormat string, desc string) error {

	cronAgent := cron.New()
	detailCronJob := DetailCronJob{Spec: timeFormat, Description: desc, Cron: cronAgent}

	if _, ok := c.ListCron[idCron]; ok {
		return fmt.Errorf("you still have cron job with id %s running", idCron)
	}

	c.ListCron[idCron] = detailCronJob

	_, err := cronAgent.AddFunc(fmt.Sprintf("CRON_TZ=Asia/Jakarta %s", timeFormat), func() {
		defer c.closeAndDeleteCron(&detailCronJob, idCron)
		callback()

	})
	cronAgent.Start()
	return err

}

func (c *CronJob) StopCronJob(idCron string) {
	cronJob := c.ListCron[idCron]
	listConfigSchedule := c.LoadConfigCronJobSchedule()
	delete(listConfigSchedule, idCron)
	dataByte, err := json.Marshal(listConfigSchedule)
	helper.PanicIfError(err)
	err = os.WriteFile("../../config/cron_jobs_schedule.json", dataByte, os.ModeAppend)
	helper.PanicIfError(err)

	if cronJob.Cron == nil {
		return
	}
	c.closeAndDeleteCron(&cronJob, idCron)
}

func (c *CronJob) closeAndDeleteCron(cronJob *DetailCronJob, idCron string) {
	cronJob.Cron.Stop()
	delete(c.ListCron, idCron)
}

func (c *CronJob) SetCronJob(idCron string, callback func(), timeFormat string, structData any, desc string, isWriteConfig bool) error {

	cronAgenct := cron.New(cron.WithLocation(time.Local))

	c.ListCron[idCron] = DetailCronJob{Spec: timeFormat, Description: desc, Cron: cronAgenct}

	_, err := cronAgenct.AddFunc(fmt.Sprintf("CRON_TZ=Asia/Jakarta %s", timeFormat), func() {
		callback()

	})
	helper.PanicIfError(err)
	listCronJobSchedule := c.LoadConfigCronJobSchedule()
	listCronJobSchedule[idCron] = cronmodel.CronJobSchedule{Spec: timeFormat, Job: structData, CreatedAt: time.Now().Local(), Desc: desc}

	if isWriteConfig {
		if err = writeCronJobConfig(listCronJobSchedule); err != nil {
			return err
		}
	}
	cronAgenct.Start()
	return nil

}

func (c *CronJob) IsCronJobRunning(idCron string) bool {
	if _, ok := c.ListCron[idCron]; ok {
		return ok
	}
	return false
}

func (c *CronJob) LoadConfigCronJobSchedule() cronmodel.CronJobConfig {

	dataFile, err := os.ReadFile("../../config/cron_jobs_schedule.json")
	if err != nil {
		os.Mkdir("../../config", 0766)
		if _, err = os.Create("../../config/cron_jobs_schedule.json"); err != nil {
			fmt.Println("Something went wrong making config cron job", err)
			return nil
		}
		dataFile, _ = os.ReadFile("../../config/cron_jobs_schedule.json")
	}
	listCronJobSchedule := cronmodel.CronJobConfig{}
	json.Unmarshal(dataFile, &listCronJobSchedule)

	return listCronJobSchedule
}

func (c *CronJob) DeleteOneCronJob(id string) error {
	listCronJob := c.LoadConfigCronJobSchedule()
	if _, ok := listCronJob[id]; !ok {
		return fmt.Errorf("id cron job not found")
	}
	if _, ok := c.ListCron[id]; !ok {
		return fmt.Errorf("id cron job not found, the config cron job not syncronize with started cron job")
	}

	delete(listCronJob, id)
	if err := writeCronJobConfig(listCronJob); err != nil {
		return err
	}
	return nil
}

func writeCronJobConfig(listCronJob cronmodel.CronJobConfig) error {
	jsonStructByte, err := json.Marshal(listCronJob)
	if err != nil {
		return err
	}
	err = os.WriteFile("../../config/cron_jobs_schedule.json", jsonStructByte, os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}
