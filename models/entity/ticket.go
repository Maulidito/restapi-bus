package entity

type Ticket struct {
	TicketId       int
	AgencyId       int
	BusId          int
	DriverId       int
	CustomerId     int
	DeparturePlace string
	ArrivalPlace   string
	Price          int
	Date           string
	Arrived        bool `default:"false"`
}
