package request

type Ticket struct {
	AgecnyId       int    `json:"agencyId" binding:"required" form:"agencyId"`
	BusId          int    `json:"busId" binding:"required" form:"busId"`
	DriverId       int    `json:"driverId" binding:"required" form:"driverId"`
	CustomerId     int    `json:"customerId" binding:"required" form:"customerId"`
	DeparturePlace string `json:"depaturePlace" binding:"required" form:"departurePlace"`
	ArrivalPlace   string `json:"arrivalPlace" binding:"required" form:"arrivalPlace"`
	Price          int    `json:"price" binding:"required" form:"price"`
	Arrived        bool   `json:"arrived"  form:"arrived"`
}

type TicketFilter struct {
	FromDate       string `form:"fromDate" binding:"required_with=ToDate,omitempty,datetime=2006-01-02,ltefield=ToDate,validatefromTodate=ToDate"`
	ToDate         string `form:"toDate" binding:"omitempty,datetime=2006-01-02"`
	OnDate         string `form:"onDate" binding:"excluded_with=ToDate FromDate,omitempty,datetime=2006-01-02"`
	DeparturePlace string `form:"departurePlace" binding:"omitempty,alpha"`
	ArrivalPlace   string `form:"arrivalPlace" binding:"omitempty,alpha"`
	Arrived        *bool  `form:"arrived" binding:"omitempty"`
	Limit          int    `form:"limit" binding:"omitempty,number"`
	PriceAbove     int    `form:"priceAbove" binding:"omitempty,number"`
	PriceBelow     int    `form:"priceBelow" binding:"omitempty,number"`
}
