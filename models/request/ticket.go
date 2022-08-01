package request

type Ticket struct {
	ScheduleId int `form:"scheduleId" binding:"required, numeric"`
	CustomerId int `form:"customerId" binding:"required, numeric" `
}

type TicketFilter struct {
	FromDate   string `form:"fromDate" binding:"required_with=ToDate,omitempty,datetime=2006-01-02,validatefromTodate=ToDate"`
	ToDate     string `form:"toDate" binding:"omitempty,datetime=2006-01-02"`
	OnDate     string `form:"onDate" binding:"excluded_with=ToDate FromDate,omitempty,datetime=2006-01-02"`
	FromAgency int    `form:"fromAgency" binding:"omitempty,number"`
	ToAgency   int    `form:"toAgency" binding:"omitempty,number"`
	Arrived    *bool  `form:"arrived" binding:"omitempty"`
	Limit      int    `form:"limit" binding:"omitempty,number"`
	PriceAbove int    `form:"priceAbove" binding:"omitempty,number"`
	PriceBelow int    `form:"priceBelow" binding:"omitempty,number"`
}
