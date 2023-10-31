package entity

import (
	"context"
	"database/sql"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type Ticket struct {
	TicketId   int
	ScheduleId int
	CustomerId int
	Date       string
	PaymentId  sql.NullString
	ExternalId string
	IsPaid     bool
	SeatNumber int
}

type TicketServiceInterface interface {
	GetAllTicket(ctx context.Context, filter *request.TicketFilter) []response.Ticket
	AddTicket(ctx context.Context, ticket *request.Ticket) response.Ticket
	GetOneTicket(ctx context.Context, ticketId int) response.Ticket
	DeleteTicket(ctx context.Context, ticketId int) response.Ticket
	GetAllTicketOnDriver(ctx context.Context, idDriver int) response.AllTicketOnDriver
	GetAllTicketOnCustomer(ctx context.Context, idCustomer int) response.AllTicketOnCustomer
	GetAllTicketOnBus(ctx context.Context, idBus int) response.AllTicketOnBus
	GetAllTicketOnAgency(ctx context.Context, idAgency int) response.AllTicketOnAgency
	GetTotalPriceAllTicket(ctx context.Context) response.AllTicketPrice
	GetTotalPriceTicketFromSpecificAgency(ctx context.Context, idAgency int) response.AllTicketPriceSpecificAgency
	GetTotalPriceTicketFromSpecificDriver(ctx context.Context, idDriver int) response.AllTicketPriceSpecificDriver
}

type TicketRepositoryInterface interface {
	GetAllTicket(ctx context.Context, filter *request.TicketFilter) []Ticket
	IsCustomerHaveUnpaidPayment(ctx context.Context, customerId int) bool
	AddTicket(ctx context.Context, ticket *Ticket)
	GetOneTicket(ctx context.Context, ticket *Ticket)
	GetOneTicketbyExternalId(ctx context.Context, externalId string) Ticket
	DeleteTicket(ctx context.Context, ticket *Ticket)
	GetAllTicketOnDriver(ctx context.Context, idDriver int) []Ticket
	GetAllTicketOnCustomer(ctx context.Context, idCustomer int) []Ticket
	GetAllTicketOnBus(ctx context.Context, idBus int) []Ticket
	GetAllTicketOnAgency(ctx context.Context, idBus int) []Ticket
	GetTotalPriceAllTicket(ctx context.Context, response *response.AllTicketPrice)
	GetTotalPriceTicketFromSpecificAgency(ctx context.Context, response *response.AllTicketPriceSpecificAgency)
	GetTotalPriceTicketFromSpecificDriver(ctx context.Context, response *response.AllTicketPriceSpecificDriver)
	UpdateTicketToPaid(ctx context.Context, externalId string, paymentId string)
	IsSeatBooked(ctx context.Context, scheduleId int, seatNumber int) bool
}
