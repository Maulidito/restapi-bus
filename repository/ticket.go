package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/response"
)

type TicketRepositoryInterface interface {
	GetAllTicket(tx *sql.Tx, ctx context.Context, filter string) []entity.Ticket
	AddTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	GetOneTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	DeleteTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	GetAllTicketOnDriver(tx *sql.Tx, ctx context.Context, idDriver int) []entity.Ticket
	GetAllTicketOnCustomer(tx *sql.Tx, ctx context.Context, idCustomer int) []entity.Ticket
	GetAllTicketOnBus(tx *sql.Tx, ctx context.Context, idBus int) []entity.Ticket
	GetAllTicketOnAgency(tx *sql.Tx, ctx context.Context, idBus int) []entity.Ticket
	UpdateArrivedTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket)
	GetTotalPriceAllTicket(tx *sql.Tx, ctx context.Context, response *response.AllTicketPrice)
	GetTotalPriceTicketFromSpecificAgency(tx *sql.Tx, ctx context.Context, response *response.AllTicketPriceSpecificAgency)
	GetTotalPriceTicketFromSpecificDriver(tx *sql.Tx, ctx context.Context, response *response.AllTicketPriceSpecificDriver)
}

type TicketRepositoryImplementation struct {
}

func NewTicketRepository() TicketRepositoryInterface {
	return &TicketRepositoryImplementation{}
}

func (repo *TicketRepositoryImplementation) GetAllTicket(tx *sql.Tx, ctx context.Context, filter string) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket "+filter)

	fmt.Println("CHECK SQL FILTER", filter)

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

func (repo *TicketRepositoryImplementation) UpdateArrivedTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket) {
	defer helper.ShouldRollback(tx)
	_, err := tx.ExecContext(ctx, "UPDATE ticket SET arrived = ? WHERE ticket_id =? ", ticket.Arrived, ticket.TicketId)
	helper.PanicIfError(err)
}

func (repo *TicketRepositoryImplementation) GetTotalPriceAllTicket(tx *sql.Tx, ctx context.Context, response *response.AllTicketPrice) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT SUM(price) as TotalPrice ,COUNT(price) as TicketCount FROM ticket ")
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&response.TotalPrice, &response.TicketCount)

		helper.PanicIfError(err)

	}

	if response.TotalPrice == nil {
		errMsg := "Not Found Single Data From Ticket"

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) GetTotalPriceTicketFromSpecificAgency(tx *sql.Tx, ctx context.Context, response *response.AllTicketPriceSpecificAgency) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT SUM(price) as TotalPrice ,COUNT(price) as TicketCount FROM ticket WHERE agency_id = ?", response.Agency.AgencyId)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {

		err := rows.Scan(&response.TotalPrice, &response.TicketCount)
		helper.PanicIfError(err)

	}
	fmt.Println("CHCEK TOTAL PRICE", response.TotalPrice)

	if response.TotalPrice == nil {
		errMsg := fmt.Sprintf(`Not Found Single Ticket From "%s" Agency With Id = %d`, response.Agency.Name, response.Agency.AgencyId)

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) GetTotalPriceTicketFromSpecificDriver(tx *sql.Tx, ctx context.Context, response *response.AllTicketPriceSpecificDriver) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT SUM(price) as TotalPrice ,COUNT(price) as TicketCount FROM ticket WHERE driver_id = ?", response.Driver.DriverId)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&response.TotalPrice, &response.TicketCount)
		helper.PanicIfError(err)

	}

	if response.TotalPrice == nil {
		errMsg := fmt.Sprintf(`Not Found Single Ticket From Driver Name "%s"  With Id = %d`, response.Driver.Name, response.Driver.DriverId)

		panic(exception.NewNotFoundError(errMsg))
	}

}
