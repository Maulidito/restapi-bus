package helper

import (
	"fmt"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"time"
)

func GenerateAutoSchedule(
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
	timeNowWithTZ := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	for timeCurrent.Before(timeEnd) {

		if timeCurrent.Equal(EndHour) || timeCurrent.After(EndHour) || timeCurrent.Before(startHour) {
			startHour = startHour.AddDate(0, 0, 1)
			EndHour = EndHour.AddDate(0, 0, 1)
			timeCurrent = startHour
		}

		if timeCurrent.After(timeNowWithTZ) {
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
