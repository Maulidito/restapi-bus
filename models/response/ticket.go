package response

type Ticket struct {
	TicketId   int    `json:"ticketId"`
	ScheduleId int    `json:"scheduleId"`
	CustomerId int    `json:"customerId"`
	Date       string `json:"date"`
}

type DetailTicket struct {
	TicketId int `json:"ticketId"`
	Schedule DetailSchedule
	Customer Customer
	Date     string `json:"date"`
}

type AllTicketOnAgency struct {
	Agency *Agency   `json:"agency"`
	Ticket *[]Ticket `json:"ticket"`
}
type AllTicketOnBus struct {
	Bus    *Bus      `json:"bus"`
	Ticket *[]Ticket `json:"ticket"`
}
type AllTicketOnDriver struct {
	Driver *Driver   `json:"driver"`
	Ticket *[]Ticket `json:"ticket"`
}
type AllTicketOnCustomer struct {
	Customer *Customer `json:"customer"`
	Ticket   *[]Ticket `json:"ticket"`
}

type AllTicketPrice struct {
	TotalPrice  *int64 `json:"totalTicket" default:"0"`
	TicketCount int    `json:"ticketCount" default:"0"`
}

type AllTicketPriceSpecificAgency struct {
	Agency      Agency `json:"agency"`
	TotalPrice  *int64 `json:"totalTicket" default:"0"`
	TicketCount int    `json:"ticketCount" default:"0"`
}

type AllTicketPriceSpecificDriver struct {
	Driver      Driver `json:"driver"`
	TotalPrice  *int64 `json:"totalTicket" default:"0"`
	TicketCount int    `json:"ticketCount" default:"0"`
}
