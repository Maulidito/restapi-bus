package request

type Ticket struct {
	BusId          int    `json:"busId"`
	DriverId       int    `json:"driverId"`
	CustomerId     int    `json:"customerId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
}
