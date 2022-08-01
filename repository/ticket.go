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
	row, err := tx.QueryContext(ctx, "SELECT ticket_id,ticket.schedule_id,customer_id,ticket.date FROM ticket "+filter)

	fmt.Println("CHECK SQL FILTER", filter)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)
	}

	return listEntityTicket
}

func (repo *TicketRepositoryImplementation) AddTicket(tx *sql.Tx, ctx context.Context, ticket *entity.Ticket) {
	defer helper.ShouldRollback(tx)
	res, err := tx.ExecContext(ctx, "INSERT INTO ticket(schedule_id,customer_id) VALUES (?,?)", &ticket.ScheduleId, &ticket.CustomerId)
	helper.PanicIfError(err)
	ticketId, err := res.LastInsertId()
	helper.PanicIfError(err)
	ticket.TicketId = int(ticketId)
}
func (repo *TicketRepositoryImplementation) GetOneTicket(tx *sql.Tx, ctx context.Context, entityTicket *entity.Ticket) {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE ticket_id = ?", entityTicket.TicketId)

	helper.PanicIfError(err)
	defer row.Close()

	if row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date)
		helper.PanicIfError(err)
		return
	}

	panic(exception.NewNotFoundError(fmt.Sprintf("TICKET ID %d NOT FOUND", entityTicket.TicketId)))
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
	row, err := tx.QueryContext(ctx, "SELECT ticket.ticket_id,ticket.schedule_id,customer_id,ticket.date FROM ticket LEFT JOIN schedule ON ticket.schedule_id=schedule.schedule_id WHERE driver_id = ?", idDriver)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date)
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
		err := row.Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket

}
func (repo *TicketRepositoryImplementation) GetAllTicketOnBus(tx *sql.Tx, ctx context.Context, idBus int) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT ticket.ticket_id,ticket.schedule_id,customer_id,ticket.date FROM ticket LEFT JOIN schedule ON ticket.schedule_id=schedule.schedule_id WHERE bus_id = ?", idBus)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket
}

func (repo *TicketRepositoryImplementation) GetAllTicketOnAgency(tx *sql.Tx, ctx context.Context, agencyId int) []entity.Ticket {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT ticket.ticket_id,ticket.schedule_id,customer_id,ticket.date FROM ticket LEFT JOIN schedule ON ticket.schedule_id=schedule.schedule_id WHERE agency_id = ? ", agencyId)

	helper.PanicIfError(err)
	defer row.Close()

	entityTicket := entity.Ticket{}
	listEntityTicket := []entity.Ticket{}

	for row.Next() {
		err := row.Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date)
		helper.PanicIfError(err)
		listEntityTicket = append(listEntityTicket, entityTicket)

	}

	return listEntityTicket
}

func (repo *TicketRepositoryImplementation) GetTotalPriceAllTicket(tx *sql.Tx, ctx context.Context, response *response.AllTicketPrice) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT COUNT(ticket.schedule_id) as ticket_count ,SUM(price) as total_price  FROM ticket LEFT JOIN schedule ON ticket.schedule_id = schedule.schedule_id")
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&response.TicketCount, &response.TotalPrice)

		helper.PanicIfError(err)

	}

	if response.TotalPrice == nil {
		errMsg := "Not Found Single Data From Ticket"

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) GetTotalPriceTicketFromSpecificAgency(tx *sql.Tx, ctx context.Context, response *response.AllTicketPriceSpecificAgency) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT COUNT(ticket.schedule_id) as ticket_count ,SUM(price) as total_price  FROM ticket LEFT JOIN schedule ON ticket.schedule_id = schedule.schedule_id GROUP BY schedule.from_agency_id HAVING schedule.from_agency_id = ?", response.Agency.AgencyId)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {

		err := rows.Scan(&response.TicketCount, &response.TotalPrice)
		helper.PanicIfError(err)

	}

	if response.TotalPrice == nil {
		errMsg := fmt.Sprintf(`Not Found Single Ticket From "%s" Agency With Id = %d`, response.Agency.Name, response.Agency.AgencyId)

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) GetTotalPriceTicketFromSpecificDriver(tx *sql.Tx, ctx context.Context, response *response.AllTicketPriceSpecificDriver) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT COUNT(ticket.schedule_id) as ticket_count ,SUM(price) as total_price FROM ticket LEFT JOIN schedule ON ticket.schedule_id = schedule.schedule_id GROUP BY schedule.driver_id HAVING schedule.driver_id = ?", response.Driver.DriverId)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&response.TicketCount, &response.TotalPrice)
		helper.PanicIfError(err)

	}

	if response.TotalPrice == nil {
		errMsg := fmt.Sprintf(`Not Found Single Ticket From Driver Name "%s"  With Id = %d`, response.Driver.Name, response.Driver.DriverId)

		panic(exception.NewNotFoundError(errMsg))
	}

}
