package testing

import (
	"fmt"
	CronCustom "restapi-bus/cron_custom"
	"testing"
	"time"
)

func TestCronEverySecond(t *testing.T) {

	customCron := CronCustom.NewCronJob()

	dataString := ""
	go func() {

		err := customCron.SetCronJobOnce("1", func() {
			dataString += "hale"
		}, "* * * * * *")

		if err != nil {
			t.Errorf("GOT ERROR SETCRON JOB %v", err)
		}
	}()

	time.Sleep(5 * time.Second)
	customCron.StopCronJob("1")
	fmt.Println("CLOSE CRON JOB", dataString)
	time.Sleep(3 * time.Second)

	t.Log("DONE")

}
