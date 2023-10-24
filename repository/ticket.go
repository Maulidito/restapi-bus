package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

var ticketRepositorySingleton *TicketRepositoryImplementation

type TicketRepositoryImplementation struct {
	conn *sql.DB
}

func NewTicketRepository(conn *sql.DB) entity.TicketRepositoryInterface {
	if ticketRepositorySingleton == nil {
		ticketRepositorySingleton = &TicketRepositoryImplementation{conn: conn}
	}
	return ticketRepositorySingleton
}

func (repo *TicketRepositoryImplementation) GetAllTicket(ctx context.Context, filter *request.TicketFilter) []entity.Ticket {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	filterString := helper.RequestFilterTicketToString(filter)
	row, err := tx.QueryContext(ctx, "SELECT ticket_id,ticket.schedule_id,customer_id,ticket.date FROM ticket "+filterString)

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

func (repo *TicketRepositoryImplementation) AddTicket(ctx context.Context, ticket *entity.Ticket) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	res, err := tx.ExecContext(ctx, "INSERT INTO ticket(schedule_id,customer_id,external_id) VALUES (?,?,?)",
		&ticket.ScheduleId, &ticket.CustomerId, &ticket.ExternalId)
	helper.PanicIfError(err)
	ticketId, err := res.LastInsertId()
	ticket.TicketId = int(ticketId)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT date FROM ticket where ticket_id = ?", ticketId).Scan(&ticket.Date)
	helper.PanicIfError(err)

}
func (repo *TicketRepositoryImplementation) GetOneTicket(ctx context.Context, entityTicket *entity.Ticket) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT * FROM ticket WHERE ticket_id = ?", entityTicket.TicketId).
		Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date, &entityTicket.ExternalId, &entityTicket.PaymentId, &entityTicket.IsPaid)

	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("TICKET ID %d NOT FOUND", entityTicket.TicketId)))
	}
}

func (repo *TicketRepositoryImplementation) IsCustomerHaveUnpaidPayment(ctx context.Context, customerId int) bool {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	row, err := tx.QueryContext(ctx, "SELECT * FROM ticket WHERE customer_id = ? AND is_paid = false", customerId)
	helper.PanicIfError(err)
	defer row.Close()

	return row.Next()

}

func (repo *TicketRepositoryImplementation) DeleteTicket(ctx context.Context, ticket *entity.Ticket) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	res, err := tx.ExecContext(ctx, "DELETE FROM ticket WHERE ticket_id = ?", ticket.TicketId)
	helper.PanicIfError(err)
	_, err = res.LastInsertId()
	helper.PanicIfError(err)

}
func (repo *TicketRepositoryImplementation) GetAllTicketOnDriver(ctx context.Context, idDriver int) []entity.Ticket {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
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
func (repo *TicketRepositoryImplementation) GetAllTicketOnCustomer(ctx context.Context, idCustomer int) []entity.Ticket {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
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
func (repo *TicketRepositoryImplementation) GetAllTicketOnBus(ctx context.Context, idBus int) []entity.Ticket {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
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

func (repo *TicketRepositoryImplementation) GetAllTicketOnAgency(ctx context.Context, agencyId int) []entity.Ticket {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	row, err := tx.QueryContext(ctx, "SELECT ticket.ticket_id,ticket.schedule_id,customer_id,ticket.date FROM ticket LEFT JOIN schedule ON ticket.schedule_id=schedule.schedule_id WHERE from_agency_id = ? ", agencyId)

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

func (repo *TicketRepositoryImplementation) GetTotalPriceAllTicket(ctx context.Context, response *response.AllTicketPrice) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT COUNT(ticket.schedule_id) as ticket_count ,SUM(price) as total_price  FROM ticket LEFT JOIN schedule ON ticket.schedule_id = schedule.schedule_id").
		Scan(&response.TicketCount, &response.TotalPrice)

	if response.TotalPrice == nil || err != nil {
		errMsg := "Not Found Single Data From Ticket"

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) GetTotalPriceTicketFromSpecificAgency(ctx context.Context, response *response.AllTicketPriceSpecificAgency) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT COUNT(ticket.schedule_id) as ticket_count ,SUM(price) as total_price  FROM ticket LEFT JOIN schedule ON ticket.schedule_id = schedule.schedule_id GROUP BY schedule.from_agency_id HAVING schedule.from_agency_id = ?", response.Agency.AgencyId).
		Scan(&response.TicketCount, &response.TotalPrice)

	if response.TotalPrice == nil || err != nil {
		errMsg := fmt.Sprintf(`Not Found Single Ticket From "%s" Agency With Id = %d`, response.Agency.Name, response.Agency.AgencyId)

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) GetTotalPriceTicketFromSpecificDriver(ctx context.Context, response *response.AllTicketPriceSpecificDriver) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT COUNT(ticket.schedule_id) as ticket_count ,SUM(price) as total_price FROM ticket LEFT JOIN schedule ON ticket.schedule_id = schedule.schedule_id GROUP BY schedule.driver_id HAVING schedule.driver_id = ?", response.Driver.DriverId).
		Scan(&response.TicketCount, &response.TotalPrice)

	if response.TotalPrice == nil || err != nil {
		errMsg := fmt.Sprintf(`Not Found Single Ticket From Driver Name "%s"  With Id = %d`, response.Driver.Name, response.Driver.DriverId)

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *TicketRepositoryImplementation) UpdateTicketToPaid(ctx context.Context, externalId string, paymentId string) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	res, err := tx.ExecContext(ctx, "UPDATE ticket SET is_paid = TRUE, payment_id=? WHERE external_id = ?", paymentId, externalId)
	helper.PanicIfError(err)
	countRowsAffected, err := res.RowsAffected()
	helper.PanicIfError(err)
	if countRowsAffected == 0 {
		helper.PanicIfError(fmt.Errorf("failed to updated ticket to paid with external_id %s and payment_id %s", externalId, paymentId))
	}

}

func (repo *TicketRepositoryImplementation) GetOneTicketbyExternalId(ctx context.Context, externalId string) entity.Ticket {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	entityTicket := entity.Ticket{}
	err = tx.QueryRowContext(ctx, "SELECT * FROM ticket WHERE external_id = ?", externalId).
		Scan(&entityTicket.TicketId, &entityTicket.ScheduleId, &entityTicket.CustomerId, &entityTicket.Date, &entityTicket.PaymentId, &entityTicket.IsPaid, &entityTicket.ExternalId)
	if err != nil {
		panic(err)
	}
	return entityTicket
}

// repo *TicketRepositoryImplementation
