package entity

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
