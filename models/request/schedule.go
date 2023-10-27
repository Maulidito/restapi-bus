package request

type Schedule struct {
	FromAgencyId int    `form:"fromAgencyId" binding:"required,number,nefield=ToAgencyId"`
	ToAgencyId   int    `form:"toAgencyId" binding:"required,number"`
	BusId        int    `form:"busId" binding:"required,number"`
	DriverId     int    `form:"driverId" binding:"required,number"`
	Price        int    `form:"price" binding:"required,number"`
	Date         string `form:"date" binding:"required,datetime=2006-01-02,validatedateafternow"`
	Arrived      bool   `default:"false" form:"arrived" binding:"omitempty" `
}

type ScheduleFilter struct {
	Limit      int    `form:"limit" binding:"omitempty,number"`
	Arrived    *bool  `form:"arrived" binding:"omitempty"`
	FromAgency int    `form:"fromAgency" binding:"omitempty,number"`
	ToAgency   int    `form:"toAgency" binding:"omitempty,number"`
	PriceAbove int    `form:"priceAbove" binding:"omitempty,number,gtecsfield=PriceBelow"`
	PriceBelow int    `form:"priceBelow" binding:"omitempty,number,ltecsfield=PriceAbove"`
	OnDate     string `form:"onDate" binding:"omitempty,datetime=2006-01-02"`
}

type ScheduleArrived struct {
	IsArrived *bool `form:"arrived" binding:"required"`
}

type AutoSchedule struct {
	FirstAgencyId  int    `form:"firstAgencyId" binding:"required,number,nefield=SecondAgencyId"`
	SecondAgencyId int    `form:"secondAgencyId" binding:"required,number"`
	Price          int    `form:"price" binding:"required,number"`
	StartHour      string `form:"startHour" binding:"required,datetime=15:04:05"`
	EndHour        string `form:"startEnd" binding:"required,datetime=15:04:05"`
	RangeSchedule  string `form:"rangeSchedule,default=3" binding:"number,min=1,max=12"`
}
