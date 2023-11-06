package service

import (
	"context"
	"encoding/json"
	"fmt"
	"restapi-bus/constant"
	croncustom "restapi-bus/cron_custom"
	"restapi-bus/exception"
	"restapi-bus/external"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/models/web"
	"restapi-bus/repository"

	"time"
)

type TicketServiceImplementation struct {
	RepoBus      entity.BusRepositoryInterface
	RepoCustomer entity.CustomerRepositoryInterface
	RepoDriver   entity.DriverRepositoryInterface
	RepoTicket   entity.TicketRepositoryInterface
	RepoAgency   entity.AgencyRepositoryInterface
	RepoSchedule entity.ScheduleRepositoryInterface
	RepoMq       repository.IMessageChannel
	Payapi       external.InterfacePayment
	cronjob      croncustom.InterfaceCronJob
	Tx           database.TrInterface
}

func NewTicketService(

	repoBus entity.BusRepositoryInterface,
	repoCustomer entity.CustomerRepositoryInterface,
	repoDriver entity.DriverRepositoryInterface,
	repoTicket entity.TicketRepositoryInterface,
	repoAgency entity.AgencyRepositoryInterface,
	repoSchedule entity.ScheduleRepositoryInterface,
	RepoMq repository.IMessageChannel,
	payapi external.InterfacePayment,
	cronJob croncustom.InterfaceCronJob,
	tx database.TrInterface,
) entity.TicketServiceInterface {

	TicketServiceStruct := TicketServiceImplementation{

		RepoBus:      repoBus,
		RepoCustomer: repoCustomer,
		RepoDriver:   repoDriver,
		RepoTicket:   repoTicket,
		RepoAgency:   repoAgency,
		RepoSchedule: repoSchedule,
		RepoMq:       RepoMq,
		Payapi:       payapi,
		cronjob:      cronJob,
		Tx:           tx,
	}

	go TicketServiceStruct.consumeWebhookQueuePaymentSuccess()
	return &TicketServiceStruct
}

func (service *TicketServiceImplementation) GetAllTicket(ctx context.Context, filter *request.TicketFilter) []response.Ticket {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	listAllTicket := service.RepoTicket.GetAllTicket(ctx, filter)

	responseListTicket := []response.Ticket{}

	for _, val := range listAllTicket {
		responseListTicket = append(responseListTicket, helper.TicketEntityToResponse(&val))
	}

	return responseListTicket
}
func (service *TicketServiceImplementation) AddTicket(ctx context.Context, ticket *request.Ticket) response.Ticket {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	ticketEntity := helper.TicketRequestToEntity(ticket)

	scheduleEntity := entity.Schedule{ScheduleId: ticket.ScheduleId}
	customerEntity := entity.Customer{CustomerId: ticket.CustomerId}
	busEntity := entity.Bus{}

	service.RepoCustomer.GetOneCustomer(ctx, &customerEntity)
	service.RepoSchedule.GetOneSchedule(ctx, &scheduleEntity)
	busEntity.BusId = scheduleEntity.BusId
	service.RepoBus.GetOneBus(ctx, &busEntity)

	if service.RepoTicket.IsCustomerHaveUnpaidPayment(ctx, customerEntity.CustomerId) {
		panic(exception.NewBadRequestError(fmt.Sprintf("Customer %s with email %s still have an Unpaid Order Payment", customerEntity.Name, customerEntity.Email)))
	}
	if ticketEntity.SeatNumber > busEntity.TotalSeat {
		panic(exception.NewBadRequestError(fmt.Sprintf("Customer Choose Seat No %d But The Bus Only Have %d", ticketEntity.SeatNumber, busEntity.TotalSeat)))
	}
	if service.RepoTicket.IsSeatBooked(ctx, ticket.ScheduleId, ticketEntity.SeatNumber) {
		panic(exception.NewBadRequestError(fmt.Sprintf("The Seat No %d Have Been Booked", ticketEntity.SeatNumber)))
	}
	respDetailSchedule := response.DetailSchedule{ScheduleId: scheduleEntity.ScheduleId, FromAgency: response.Agency{AgencyId: scheduleEntity.FromAgencyId}, ToAgency: response.Agency{AgencyId: scheduleEntity.ToAgencyId}, Driver: response.Driver{DriverId: scheduleEntity.DriverId}, Bus: response.Bus{BusId: scheduleEntity.BusId}}

	service.RepoSchedule.GetOneDetailSchedule(ctx, &respDetailSchedule)

	dataVirtualAccount := service.Payapi.MakeVirtualAccount(ctx,
		"Bus Agency Ticket Payment",
		fmt.Sprintf("FIXED-VA-%s-%s-%d", customerEntity.Email, customerEntity.PhoneNumber, time.Now().UnixNano()),
		ticket.BankCode,
		scheduleEntity.Price)
	ticketEntity.ExternalId = fmt.Sprintf("%v", dataVirtualAccount["external_id"])

	service.RepoTicket.AddTicket(ctx, &ticketEntity)
	date, err := time.Parse(time.DateTime, ticketEntity.Date)
	helper.PanicIfError(err)
	expiration_date := dataVirtualAccount["expiration_date"]
	time_expire, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", expiration_date))
	time_expire = time_expire.Local()
	helper.PanicIfError(err)
	expiration_date_string := time_expire.Format("2006-01-02")
	expiration_time_string := int(time.Since(time_expire).Abs().Minutes())
	expiration_hour_string := int(time.Since(time_expire).Abs().Hours())
	expiration_day_string := int(time.Since(time_expire).Abs().Hours() / 24)
	respTicketOrder := response.TicketOrder{
		TicketId:            ticketEntity.TicketId,
		Schedule:            respDetailSchedule,
		Customer:            helper.CustomerEntityToResponse(&customerEntity),
		Date:                fmt.Sprintf("%v", date.Format(time.DateOnly)),
		VirtualAccontNumber: fmt.Sprintf("%v", dataVirtualAccount["account_number"]),
		ExpiryDate:          expiration_date_string,
		ExpiryMinute:        expiration_time_string,
		ExpiryHour:          expiration_hour_string,
		ExpiryDay:           expiration_day_string,
		BankCode:            fmt.Sprintf("%v", dataVirtualAccount["bank_code"]),
		MerchantCode:        fmt.Sprintf("%v", dataVirtualAccount["merchant_code"]),
		SeatNumber:          ticketEntity.SeatNumber,
	}

	ticketResponse := helper.TicketEntityToResponse(&ticketEntity)

	respTicketOrderByte, err := json.Marshal(respTicketOrder)
	helper.PanicIfError(err)

	service.RepoMq.PublishToEmailServiceTopic(ctx, constant.TOPIC_PAYMENT_EMAIL, constant.QUEUE_PAYMENT, respTicketOrderByte)

	err = service.cronjob.SetCronJobOnce(
		dataVirtualAccount["external_id"].(string),
		func() {
			ctx = context.Background()
			ctx = service.Tx.BeginTransactionWithContext(ctx)
			defer service.Tx.DoCommitOrRollbackWithContext(ctx)
			service.RepoTicket.DeleteTicket(ctx, &ticketEntity)
			err := recover()
			if err != nil {

				fmt.Printf("Cant delete ticket %d, maybe it already deleted", ticketEntity.TicketId)
			}
		},
		fmt.Sprintf("%d %d * * * ", time_expire.Minute(), time_expire.Hour()),
		fmt.Sprintf("Once - Ticket %d - Time %s", ticketResponse.TicketId, time_expire.String()),
	)

	helper.PanicIfError(err)

	return ticketResponse

}
func (service *TicketServiceImplementation) GetOneTicket(ctx context.Context, ticketId int) response.Ticket {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	ticketEntity := entity.Ticket{TicketId: ticketId}
	service.RepoTicket.GetOneTicket(ctx, &ticketEntity)

	return helper.TicketEntityToResponse(&ticketEntity)

}
func (service *TicketServiceImplementation) DeleteTicket(ctx context.Context, ticketId int) response.Ticket {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	ticketEntity := entity.Ticket{TicketId: ticketId}

	service.RepoTicket.GetOneTicket(ctx, &ticketEntity)
	service.RepoTicket.DeleteTicket(ctx, &ticketEntity)

	return helper.TicketEntityToResponse(&ticketEntity)
}
func (service *TicketServiceImplementation) GetAllTicketOnDriver(ctx context.Context, idDriver int) response.AllTicketOnDriver {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
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
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
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
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
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
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
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
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	response := response.AllTicketPrice{}

	service.RepoTicket.GetTotalPriceAllTicket(ctx, &response)
	return response

}
func (service *TicketServiceImplementation) GetTotalPriceTicketFromSpecificAgency(ctx context.Context, idAgency int) response.AllTicketPriceSpecificAgency {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	agencyEntity := entity.Agency{AgencyId: idAgency}
	service.RepoAgency.GetOneAgency(ctx, &agencyEntity)
	response := response.AllTicketPriceSpecificAgency{Agency: helper.AgencyEntityToResponse(&agencyEntity)}
	service.RepoTicket.GetTotalPriceTicketFromSpecificAgency(ctx, &response)
	return response

}
func (service *TicketServiceImplementation) GetTotalPriceTicketFromSpecificDriver(ctx context.Context, idDriver int) response.AllTicketPriceSpecificDriver {
	ctx = service.Tx.BeginTransactionWithContext(ctx)
	defer service.Tx.DoCommitOrRollbackWithContext(ctx)
	driverEntity := entity.Driver{DriverId: idDriver}
	service.RepoDriver.GetOneDriverOnSpecificAgency(ctx, &driverEntity)
	response := response.AllTicketPriceSpecificDriver{Driver: helper.DriverEntityToResponse(&driverEntity)}
	service.RepoTicket.GetTotalPriceTicketFromSpecificDriver(ctx, &response)
	return response
}

func (service *TicketServiceImplementation) consumeWebhookQueuePaymentSuccess() {

	ctx := context.Background()
	messages := service.RepoMq.ConsumeQueue(ctx, constant.CONSUMER_PAYMENT_WEBHOOK, constant.QUEUE_PAYMENT_WEBHOOK)
	paymentSuccess := web.PaymentSuccess{}

	go func() {
		for msg := range messages {
			json.Unmarshal(msg.Body, &paymentSuccess)
			service.cronjob.StopCronJob(paymentSuccess.ExternalID)
			service.RepoTicket.UpdateTicketToPaid(ctx, paymentSuccess.ExternalID, paymentSuccess.PaymentID)
			Ticket := service.RepoTicket.GetOneTicketbyExternalId(ctx, paymentSuccess.ExternalID)
			err := recover()
			if err != nil {
				fmt.Println(err.(error).Error())
				continue
			}
			Schedule := response.DetailSchedule{ScheduleId: Ticket.ScheduleId}
			Customer := entity.Customer{CustomerId: Ticket.CustomerId}
			service.RepoSchedule.GetOneDetailSchedule(ctx, &Schedule)
			service.RepoCustomer.GetOneCustomer(ctx, &Customer)

			respDetailTicket := response.DetailTicket{
				TicketId:   Ticket.TicketId,
				Schedule:   Schedule,
				Customer:   helper.CustomerEntityToResponse(&Customer),
				Date:       Ticket.Date,
				PaymentId:  Ticket.PaymentId.String,
				ExternalId: Ticket.ExternalId,
				IsPaid:     Ticket.IsPaid,
				SeatNumber: Ticket.SeatNumber,
			}
			respDetailTicketByte, errorJson := json.Marshal(respDetailTicket)
			helper.PanicIfError(errorJson)
			fmt.Printf("SUCCESS PAYMENT WITH Payment ID %s", paymentSuccess.PaymentID)
			service.RepoMq.PublishToEmailServiceTopic(ctx, constant.TOPIC_TICKET_EMAIL, constant.QUEUE_TICKET, respDetailTicketByte)

		}
	}()
}
