package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type TicketRepositoryInterface interface {
	GetAllTicket(tx *sql.Tx, ctx context.Context) []entity.Ticket
	AddTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	GetOneTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	DeleteTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	GetAllTicketOnDriver(tx *sql.Tx, ctx context.Context, idDriver int) []entity.Ticket
	GetAllTicketOnCustomer(tx *sql.Tx, ctx context.Context, idCustomer int) []entity.Ticket
	GetAllTicketOnBus(tx *sql.Tx, ctx context.Context, idBus int) []entity.Ticket
	GetAllTicketOnAgency(tx *sql.Tx, ctx context.Context, idBus int) []entity.Ticket
}

type TicketRepositoryImplementation struct {
}

func NewTicketRepository() TicketRepositoryInterface {
	return &TicketRepositoryImplementation{}
}

func (repo *TicketRepositoryImplementation) GetAllTicket(tx *sql.Tx, ctx context.Context) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket")

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.AgencyId, &entityTicket.BusId, &entityTicket.DriverId, &entityTicket.CustomerId, &entityTicket.DeparturePlace, &entityTicket.ArrivalPlace, &entityTicket.Price, &entityTicket.Date, &entityTicket.Arrived)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket
}

func (repo *TicketRepositoryImplementation) AddTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket) {
	defer helper.ShouldRollback(tx)
	res, err := tx.ExecContext(ctx, "INSERT INTO ticket(agency_id,bus_id,driver_id,customer_id,departure_place,arrival_place,price,arrived) VALUES (?,?,?,?,?,?,?,?)", ticket.AgencyId, ticket.BusId, ticket.DriverId, ticket.CustomerId, ticket.DeparturePlace, ticket.ArrivalPlace, ticket.Price, ticket.Arrived)
	helper.PanicIfError(err)
	ticketId, err := res.LastInsertId()
	helper.PanicIfError(err)
	ticket.TicketId = int(ticketId)
}
func (repo *TicketRepositoryImplementation) GetOneTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket) {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE ticket_id = ?", ticket.TicketId)

	helper.PanicIfError(err)
	defer row.Close()

	if row.Next() {
		err := row.Scan(&ticket.TicketId, &ticket.AgencyId, &ticket.BusId, &ticket.DriverId, &ticket.CustomerId, &ticket.DeparturePlace, &ticket.ArrivalPlace, &ticket.Price, &ticket.Date, &ticket.Arrived)
		helper.PanicIfError(err)
		return
	}

	panic(exception.NewNotFoundError(fmt.Sprintf("TICKET ID %d NOT FOUND", ticket.TicketId)))
}
func (repo *TicketRepositoryImplementation) DeleteTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket) {
	defer helper.ShouldRollback(tx)
	res, err := tx.ExecContext(ctx, "DELETE FROM ticket WHERE ticket_id = ?", ticket.TicketId)
	helper.PanicIfError(err)
	_, err = res.LastInsertId()
	helper.PanicIfError(err)

}
func (repo *TicketRepositoryImplementation) GetAllTicketOnDriver(tx *sql.Tx, ctx context.Context, idDriver int) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE driver_id = ?", idDriver)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.AgencyId, &entityTicket.BusId, &entityTicket.DriverId, &entityTicket.CustomerId, &entityTicket.DeparturePlace, &entityTicket.ArrivalPlace, &entityTicket.Price, &entityTicket.Date, &entityTicket.Arrived)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket
}
func (repo *TicketRepositoryImplementation) GetAllTicketOnCustomer(tx *sql.Tx, ctx context.Context, idCustomer int) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE customer_id = ?", idCustomer)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.AgencyId, &entityTicket.BusId, &entityTicket.DriverId, &entityTicket.CustomerId, &entityTicket.DeparturePlace, &entityTicket.ArrivalPlace, &entityTicket.Price, &entityTicket.Date, &entityTicket.Arrived)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket

}
func (repo *TicketRepositoryImplementation) GetAllTicketOnBus(tx *sql.Tx, ctx context.Context, idBus int) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE bus_id = ?", idBus)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.AgencyId, &entityTicket.BusId, &entityTicket.DriverId, &entityTicket.CustomerId, &entityTicket.DeparturePlace, &entityTicket.ArrivalPlace, &entityTicket.Price, &entityTicket.Date, &entityTicket.Arrived)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket
}

func (repo *TicketRepositoryImplementation) GetAllTicketOnAgency(tx *sql.Tx, ctx context.Context, agencyId int) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE agency_id =? ", agencyId)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.AgencyId, &entityTicket.BusId, &entityTicket.DriverId, &entityTicket.CustomerId, &entityTicket.DeparturePlace, &entityTicket.ArrivalPlace, &entityTicket.Price, &entityTicket.Date, &entityTicket.Arrived)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket
}
