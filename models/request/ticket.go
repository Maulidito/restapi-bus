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
	FromDate       string `form:"formDate" binding:"dateTime"`
	ToDate         string `form:"toDate" binding:"dateTime ,required_with=formDate"`
	OnDate         string `form:"onDate" binding:"dateTime"`
	DeparturePlace string `form:"departurePlace" binding:"alpha"`
	ArrivalPlace   string `form:"arrivalPlace" binding:"alpha"`
	Arrived        bool   `form:"arrived" binding:"boolean"`
	Limit          int    `form:"limit" binding:"number"`
	PriceAbove     int    `form:"priceAbove" binding:"number"`
	PriceBelow     int    `form:"priceBelow" binding:"number"`
}
