package request

type Ticket struct {
	AgecnyId       int    `json:"agencyId" binding:"required" form:"agencyId"`
	BusId          int    `json:"busId" binding:"required" form:"busId"`
	DriverId       int    `json:"driverId" binding:"required" form:"driverId"`
	CustomerId     int    `json:"customerId" binding:"required" form:"customerId"`
	DeparturePlace string `json:"depaturePlace" binding:"required" form:"departurePlace"`
	ArrivalPlace   string `json:"arrivalPlace" binding:"required" form:"arrivalPlace"`
	Price          int    `json:"price" binding:"required" form:"price"`
	Arrived        bool   `json:"arrived" binding:"required" form:"arrived"`
}
