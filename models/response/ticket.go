package response

type Ticket struct {
	TicketId       int    `json:"ticketId"`
	BusId          int    `json:"busId"`
	DriverId       int    `json:"driverId"`
	CustomerId     int    `json:"customerId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
}
