package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type TicketServiceImplementation struct {
	Db           *sql.DB
	RepoBus      entity.BusRepositoryInterface
	RepoCustomer entity.CustomerRepositoryInterface
	RepoDriver   entity.DriverRepositoryInterface
	RepoTicket   entity.TicketRepositoryInterface
	RepoAgency   entity.AgencyRepositoryInterface
	RepoSchedule entity.ScheduleRepositoryInterface
	RepoMq       repository.IMessageChannel
}

func NewTicketService(
	db *sql.DB,
	repoBus entity.BusRepositoryInterface,
	repoCustomer entity.CustomerRepositoryInterface,
	repoDriver entity.DriverRepositoryInterface,
	repoTicket entity.TicketRepositoryInterface,
	repoAgency entity.AgencyRepositoryInterface,
	repoSchedule entity.ScheduleRepositoryInterface,
	RepoMq repository.IMessageChannel,
) entity.TicketServiceInterface {

	return &TicketServiceImplementation{Db: db, RepoBus: repoBus, RepoCustomer: repoCustomer, RepoDriver: repoDriver, RepoTicket: repoTicket, RepoAgency: repoAgency, RepoSchedule: repoSchedule, RepoMq: RepoMq}
}

func (service *TicketServiceImplementation) GetAllTicket(ctx context.Context, filter *request.TicketFilter) []response.Ticket {

	listAllTicket := service.RepoTicket.GetAllTicket(ctx, filter)

	responseListTicket := []response.Ticket{}

	for _, val := range listAllTicket {
		responseListTicket = append(responseListTicket, helper.TicketEntityToResponse(&val))
	}

	return responseListTicket
}
func (service *TicketServiceImplementation) AddTicket(ctx context.Context, ticket *request.Ticket) {

	ticketEntity := helper.TicketRequestToEntity(ticket)

	scheduleEntity := entity.Schedule{ScheduleId: ticket.ScheduleId}
	customerEntity := entity.Customer{CustomerId: ticket.CustomerId}

	chanErr := make(chan error, 1)

	go func() {
		defer func() {
			tempErr := recover()

			if tempErr != nil {
				chanErr <- tempErr.(error)
			}

			close(chanErr)
		}()

		service.RepoSchedule.GetOneSchedule(ctx, &scheduleEntity)
		service.RepoCustomer.GetOneCustomer(ctx, &customerEntity)
	}()
	helper.PanicIfError(<-chanErr)
	// rabbitmq
	respDetailSchedule := response.DetailSchedule{ScheduleId: scheduleEntity.ScheduleId, FromAgency: response.Agency{AgencyId: scheduleEntity.FromAgencyId}, ToAgency: response.Agency{AgencyId: scheduleEntity.ToAgencyId}, Driver: response.Driver{DriverId: scheduleEntity.DriverId}, Bus: response.Bus{BusId: scheduleEntity.BusId}}

	service.RepoSchedule.GetOneDetailSchedule(ctx, &respDetailSchedule)

	service.RepoTicket.AddTicket(ctx, &ticketEntity)

	respDetailTicket := response.DetailTicket{TicketId: ticketEntity.TicketId, Customer: helper.CustomerEntityToResponse(&customerEntity), Schedule: respDetailSchedule}

	respDetailTicketByte, err := json.Marshal(respDetailTicket)
	helper.PanicIfError(err)

	service.RepoMq.PublishToEmailService(ctx, respDetailTicketByte)
	//

}
func (service *TicketServiceImplementation) GetOneTicket(ctx context.Context, ticketId int) response.Ticket {

	ticketEntity := entity.Ticket{TicketId: ticketId}
	service.RepoTicket.GetOneTicket(ctx, &ticketEntity)

	return helper.TicketEntityToResponse(&ticketEntity)

}
func (service *TicketServiceImplementation) DeleteTicket(ctx context.Context, ticketId int) response.Ticket {

	ticketEntity := entity.Ticket{TicketId: ticketId}

	service.RepoTicket.GetOneTicket(ctx, &ticketEntity)
	service.RepoTicket.DeleteTicket(ctx, &ticketEntity)

	return helper.TicketEntityToResponse(&ticketEntity)
}
func (service *TicketServiceImplementation) GetAllTicketOnDriver(ctx context.Context, idDriver int) response.AllTicketOnDriver {

	driverEntity := entity.Driver{DriverId: idDriver}

	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driverEntity)
	listTicket := service.RepoTicket.GetAllTicketOnDriver(ctx, idDriver)
	responseListTicket := []response.Ticket{}
	for _, val := range listTicket {
		responseListTicket = append(responseListTicket, helper.TicketEntityToResponse(&val))
	}

	responseDriver := helper.DriverEntityToResponse(&driverEntity)
	responseAllTicket := response.AllTicketOnDriver{
		Driver: &responseDriver,
		Ticket: &responseListTicket,
	}

	return responseAllTicket
}

func (service *TicketServiceImplementation) GetAllTicketOnCustomer(ctx context.Context, idCustomer int) response.AllTicketOnCustomer {

	customerEntity := entity.Customer{CustomerId: idCustomer}
	service.RepoCustomer.GetOneCustomer(ctx, &customerEntity)

	listTicket := service.RepoTicket.GetAllTicketOnCustomer(ctx, idCustomer)

	customerResponse := helper.CustomerEntityToResponse(&customerEntity)

	responseListCustomer := []response.Ticket{}
	for _, val := range listTicket {
		responseListCustomer = append(responseListCustomer, helper.TicketEntityToResponse(&val))
	}

	responseAllTicket := response.AllTicketOnCustomer{
		Customer: &customerResponse,
		Ticket:   &responseListCustomer,
	}

	return responseAllTicket

}
func (service *TicketServiceImplementation) GetAllTicketOnBus(ctx context.Context, idBus int) response.AllTicketOnBus {

	busEntity := entity.Bus{BusId: idBus}
	service.RepoBus.GetOneBus(ctx, &busEntity)

	listTicket := service.RepoTicket.GetAllTicketOnBus(ctx, idBus)

	busResponse := helper.BusEntityToResponse(&busEntity)

	responseListBus := []response.Ticket{}
	for _, val := range listTicket {
		responseListBus = append(responseListBus, helper.TicketEntityToResponse(&val))
	}

	responseAllTicket := response.AllTicketOnBus{
		Bus:    &busResponse,
		Ticket: &responseListBus,
	}

	return responseAllTicket

}

func (service *TicketServiceImplementation) GetAllTicketOnAgency(ctx context.Context, idAgency int) response.AllTicketOnAgency {

	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, &agencyEntity)

	listTicket := service.RepoTicket.GetAllTicketOnAgency(ctx, idAgency)

	agencyResponse := helper.AgencyEntityToResponse(&agencyEntity)

	responseListAgency := []response.Ticket{}
	for _, val := range listTicket {
		responseListAgency = append(responseListAgency, helper.TicketEntityToResponse(&val))
	}

	responseAllTicket := response.AllTicketOnAgency{
		Agency: &agencyResponse,
		Ticket: &responseListAgency,
	}

	return responseAllTicket
}

func (service *TicketServiceImplementation) GetTotalPriceAllTicket(ctx context.Context) response.AllTicketPrice {

	response := response.AllTicketPrice{}

	service.RepoTicket.GetTotalPriceAllTicket(ctx, &response)
	return response

}
func (service *TicketServiceImplementation) GetTotalPriceTicketFromSpecificAgency(ctx context.Context, idAgency int) response.AllTicketPriceSpecificAgency {

	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, &agencyEntity)
	response := response.AllTicketPriceSpecificAgency{Agency: helper.AgencyEntityToResponse(&agencyEntity)}
	service.RepoTicket.GetTotalPriceTicketFromSpecificAgency(ctx, &response)
	return response

}
func (service *TicketServiceImplementation) GetTotalPriceTicketFromSpecificDriver(ctx context.Context, idDriver int) response.AllTicketPriceSpecificDriver {

	driverEntity := entity.Driver{DriverId: idDriver}
	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driverEntity)
	response := response.AllTicketPriceSpecificDriver{Driver: helper.DriverEntityToResponse(&driverEntity)}
	service.RepoTicket.GetTotalPriceTicketFromSpecificDriver(ctx, &response)
	return response
}
