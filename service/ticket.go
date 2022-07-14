package service

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type TicketServiceInterface interface {
	GetAllTicket(ctx context.Context) []response.Ticket
	AddTicket(ctx context.Context, ticket *request.Ticket)
	GetOneTicket(ctx context.Context, ticketId int) response.Ticket
	DeleteTicket(ctx context.Context, ticketId int) response.Ticket
	GetAllTicketOnDriver(ctx context.Context, idDriver int) response.AllTicketOnDriver
	GetAllTicketOnCustomer(ctx context.Context, idCustomer int) response.AllTicketOnCustomer
	GetAllTicketOnBus(ctx context.Context, idBus int) response.AllTicketOnBus
	GetAllTicketOnAgency(ctx context.Context, idAgency int) response.AllTicketOnAgency
	UpdateArrivedTicket(ctx context.Context, idTicket int, arrived bool) response.Ticket
}

type TicketServiceImplementation struct {
	Db           *sql.DB
	RepoBus      repository.BusRepositoryInterface
	RepoCustomer repository.CustomerRepositoryInterface
	RepoDriver   repository.DriverRepositoryInterface
	RepoTicket   repository.TicketRepositoryInterface
	RepoAgency   repository.AgencyRepositoryInterface
}

func NewTicketService(
	db *sql.DB,
	repoBus repository.BusRepositoryInterface,
	repoCustomer repository.CustomerRepositoryInterface,
	repoDriver repository.DriverRepositoryInterface,
	repoTicket repository.TicketRepositoryInterface,
	repoAgency repository.AgencyRepositoryInterface,
) TicketServiceInterface {

	return &TicketServiceImplementation{Db: db, RepoBus: repoBus, RepoCustomer: repoCustomer, RepoDriver: repoDriver, RepoTicket: repoTicket, RepoAgency: repoAgency}
}

func (service *TicketServiceImplementation) GetAllTicket(ctx context.Context) []response.Ticket {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	listAllTicket := service.RepoTicket.GetAllTicket(tx, ctx)

	responseListTicket := []response.Ticket{}

	for _, val := range listAllTicket {
		responseListTicket = append(responseListTicket, helper.TicketEntityToResponse(&val))
	}

	return responseListTicket
}
func (service *TicketServiceImplementation) AddTicket(ctx context.Context, ticket *request.Ticket) {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	ticketEntity := helper.TicketRequestToEntity(ticket)
	busEntity := entity.Bus{BusId: ticket.BusId, AgencyId: ticket.AgecnyId}
	driverEntity := entity.Driver{DriverId: ticket.DriverId, AgencyId: ticket.AgecnyId}
	agencyEntity := entity.Agency{AgencyId: ticket.AgecnyId}
	customerEntity := entity.Customer{CustomerId: ticket.CustomerId}
	chanErr := make(chan error, 1)

	service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)
	go func() {
		defer func() {
			tempErr := recover()

			if tempErr != nil {
				chanErr <- tempErr.(error)
			}

			close(chanErr)
		}()

		service.RepoBus.GetOneBus(ctx, tx, &busEntity)
		service.RepoDriver.GetOneDriverOnSpecificAgency(tx, ctx, &driverEntity)
		service.RepoCustomer.GetOneCustomer(ctx, tx, &customerEntity)
	}()

	helper.PanicIfError(<-chanErr)
	service.RepoTicket.AddTicket(tx, ctx, &ticketEntity)

}
func (service *TicketServiceImplementation) GetOneTicket(ctx context.Context, ticketId int) response.Ticket {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)
	ticketEntity := entity.Ticket{TicketId: ticketId}
	service.RepoTicket.GetOneTicket(tx, ctx, &ticketEntity)

	return helper.TicketEntityToResponse(&ticketEntity)

}
func (service *TicketServiceImplementation) DeleteTicket(ctx context.Context, ticketId int) response.Ticket {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)
	ticketEntity := entity.Ticket{TicketId: ticketId}

	service.RepoTicket.GetOneTicket(tx, ctx, &ticketEntity)
	service.RepoTicket.DeleteTicket(tx, ctx, &ticketEntity)

	return helper.TicketEntityToResponse(&ticketEntity)
}
func (service *TicketServiceImplementation) GetAllTicketOnDriver(ctx context.Context, idDriver int) response.AllTicketOnDriver {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	driverEntity := entity.Driver{DriverId: idDriver}

	service.RepoDriver.GetOneDriverOnSpecificAgency(tx, ctx, &driverEntity)
	listTicket := service.RepoTicket.GetAllTicketOnDriver(tx, ctx, idDriver)
	responseListTicket := []response.TicketNoDriver{}
	for _, val := range listTicket {
		responseListTicket = append(responseListTicket, helper.TicketEntityToResponseTicketNoDriver(&val))
	}

	responseDriver := helper.DriverEntityToResponse(&driverEntity)
	responseAllTicket := response.AllTicketOnDriver{
		Driver: &responseDriver,
		Ticket: &responseListTicket,
	}

	return responseAllTicket
}

func (service *TicketServiceImplementation) GetAllTicketOnCustomer(ctx context.Context, idCustomer int) response.AllTicketOnCustomer {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	customerEntity := entity.Customer{CustomerId: idCustomer}
	service.RepoCustomer.GetOneCustomer(ctx, tx, &customerEntity)

	listTicket := service.RepoTicket.GetAllTicketOnCustomer(tx, ctx, idCustomer)

	customerResponse := helper.CustomerEntityToResponse(&customerEntity)

	responseListCustomer := []response.TicketNoCustomer{}
	for _, val := range listTicket {
		responseListCustomer = append(responseListCustomer, helper.TicketEntityToResponseTicketNoCustomer(&val))
	}

	responseAllTicket := response.AllTicketOnCustomer{
		Customer: &customerResponse,
		Ticket:   &responseListCustomer,
	}

	return responseAllTicket

}
func (service *TicketServiceImplementation) GetAllTicketOnBus(ctx context.Context, idBus int) response.AllTicketOnBus {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	busEntity := entity.Bus{BusId: idBus}
	service.RepoBus.GetOneBus(ctx, tx, &busEntity)

	listTicket := service.RepoTicket.GetAllTicketOnBus(tx, ctx, idBus)

	busResponse := helper.BusEntityToResponse(&busEntity)

	responseListBus := []response.TicketNoBus{}
	for _, val := range listTicket {
		responseListBus = append(responseListBus, helper.TicketEntityToResponseTicketNoBus(&val))
	}

	responseAllTicket := response.AllTicketOnBus{
		Bus:    &busResponse,
		Ticket: &responseListBus,
	}

	return responseAllTicket

}

func (service *TicketServiceImplementation) GetAllTicketOnAgency(ctx context.Context, idAgency int) response.AllTicketOnAgency {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, tx, &agencyEntity)

	listTicket := service.RepoTicket.GetAllTicketOnAgency(tx, ctx, idAgency)

	agencyResponse := helper.AgencyEntityToResponse(&agencyEntity)

	responseListAgency := []response.TicketNoAgency{}
	for _, val := range listTicket {
		responseListAgency = append(responseListAgency, helper.TicketEntityToResponseTicketNoAgency(&val))
	}

	responseAllTicket := response.AllTicketOnAgency{
		Agency: &agencyResponse,
		Ticket: &responseListAgency,
	}

	return responseAllTicket
}

func (service *TicketServiceImplementation) UpdateArrivedTicket(ctx context.Context, idTicket int, arrived bool) response.Ticket {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	defer helper.DoCommit(tx)

	entityTicket := entity.Ticket{TicketId: idTicket}
	service.RepoTicket.GetOneTicket(tx, ctx, &entityTicket)
	entityTicket.Arrived = arrived
	service.RepoTicket.UpdateArrivedTicket(tx, ctx, &entityTicket)

	return helper.TicketEntityToResponse(&entityTicket)
}
