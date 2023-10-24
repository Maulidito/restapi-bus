package response

type Ticket struct {
	TicketId   int    `json:"ticketId"`
	ScheduleId int    `json:"scheduleId"`
	CustomerId int    `json:"customerId"`
	Date       string `json:"date"`
	ExternalId string `json:"external_id"`
	PaymentId  string `json:"payment_id"`
}

type DetailTicket struct {
	TicketId   int `json:"ticketId"`
	Schedule   DetailSchedule
	Customer   Customer
	Date       string `json:"date"`
	PaymentId  string `json:"paymentId"`
	IsPaid     bool   `json:"isPaid"`
	ExternalId string `json:"external_id"`
}

type TicketOrder struct {
	TicketId            int `json:"ticketId"`
	Schedule            DetailSchedule
	Customer            Customer
	Date                string `json:"date"`
	ExpiryDate          string `json:"expirate_date"`
	ExpiryMinute        int    `json:"expirate_minute"`
	ExpiryHour          int    `json:"expirate_hour"`
	ExpiryDay           int    `json:"expirate_day"`
	VirtualAccontNumber string `json:"virtual_account_number"`
	IsPaid              bool   `json:"isPaid"`
	BankCode            string `json:"bank_code"`
	MerchantCode        string `json:"merchant_code"`
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
