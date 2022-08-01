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
	PriceAbove int    `form:"priceAbove" binding:"omitempty,number,gte=PriceBelow"`
	PriceBelow int    `form:"priceBelow" binding:"omitempty,number,lte=PriceAbove"`
	OnDate     string `form:"onDate" binding:"omitempty,datetime=2006-01-02"`
}

type ScheduleArrived struct {
	IsArrived *bool `form:"arrived" binding:"required"`
}
