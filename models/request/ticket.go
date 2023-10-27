package request

type Ticket struct {
	ScheduleId int    `form:"scheduleId" json:"scheduleId" binding:"required,numeric"`
	CustomerId int    `form:"customerId" json:"customerId" binding:"required,numeric"`
	BankCode   string `json:"bank_code" form:"bank_code" binding:"required,oneof=BNI BJB BRI BSI BNC MANDIRI PERMATA BCA CIMB DBS"`
	SeatNumber int    `form:"seatNumber" json:"seatNumber" binding:"required,numeric,min=1"`
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
