package response

type Schedule struct {
	ScheduleId   int    `json:"scheduleId" `
	FromAgencyId int    `json:"fromAgencyId" `
	ToAgencyId   int    `json:"toAgencyId" `
	BusId        int    `json:"busId" `
	DriverId     int    `json:"driverId" `
	Price        int    `json:"price" `
	Date         string `json:"date" `
	Arrived      bool   `default:"false" json:"arrived"  `
}
