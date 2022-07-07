package response

type Ticket struct {
	TicketId       int    `json:"ticketId"`
	AgencyId       int    `json:"agencyId"`
	BusId          int    `json:"busId"`
	DriverId       int    `json:"driverId"`
	CustomerId     int    `json:"customerId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
	Arrived        bool   `json:"arrived"`
}

type TicketNoAgency struct {
	TicketId       int    `json:"ticketId"`
	BusId          int    `json:"busId"`
	DriverId       int    `json:"driverId"`
	CustomerId     int    `json:"customerId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
	Arrived        bool   `json:"arrived"`
}

type TicketNoBus struct {
	TicketId       int    `json:"ticketId"`
	AgencyId       int    `json:"agencyId"`
	DriverId       int    `json:"driverId"`
	CustomerId     int    `json:"customerId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
	Arrived        bool   `json:"arrived"`
}

type TicketNoCustomer struct {
	TicketId       int    `json:"ticketId"`
	AgencyId       int    `json:"agencyId"`
	BusId          int    `json:"busId"`
	DriverId       int    `json:"driverId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
	Arrived        bool   `json:"arrived"`
}

type TicketNoDriver struct {
	TicketId       int    `json:"ticketId"`
	AgencyId       int    `json:"agencyId"`
	BusId          int    `json:"busId"`
	CustomerId     int    `json:"customerId"`
	DeparturePlace string `json:"depaturePlace"`
	ArrivalPlace   string `json:"arrivalPlace"`
	Price          int    `json:"price"`
	Date           string `json:"date"`
	Arrived        bool   `json:"arrived"`
}

type AllTicketOnAgency struct {
	Agency *Agency           `json:"agency"`
	Ticket *[]TicketNoAgency `json:"ticket"`
}
type AllTicketOnBus struct {
	Bus    *Bus           `json:"bus"`
	Ticket *[]TicketNoBus `json:"ticket"`
}
type AllTicketOnDriver struct {
	Driver *Driver           `json:"driver"`
	Ticket *[]TicketNoDriver `json:"ticket"`
}
type AllTicketOnCustomer struct {
	Customer *Customer           `json:"customer"`
	Ticket   *[]TicketNoCustomer `json:"ticket"`
}
