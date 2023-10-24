package helper

import (
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

func CustomerEntityToResponse(customer *entity.Customer) response.Customer {
	return response.Customer{
		CustomerId:  customer.CustomerId,
		Name:        customer.Name,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}
}

func AgencyEntityToResponse(agency *entity.Agency) response.Agency {
	return response.Agency{
		AgencyId: agency.AgencyId,
		Name:     agency.Name,
		Place:    agency.Place,
	}
}

func DriverEntityToResponse(driver *entity.Driver) response.Driver {
	return response.Driver{
		AgencyId: driver.AgencyId,
		Name:     driver.Name,
		DriverId: driver.DriverId,
	}
}

func BusEntityToResponse(bus *entity.Bus) response.Bus {
	return response.Bus{
		BusId:       bus.BusId,
		AgencyId:    bus.AgencyId,
		NumberPlate: bus.NumberPlate,
	}
}

func TicketEntityToResponse(ticket *entity.Ticket) response.Ticket {
	return response.Ticket{
		TicketId:   ticket.TicketId,
		ScheduleId: ticket.ScheduleId,
		CustomerId: ticket.CustomerId,
		Date:       ticket.Date,
		ExternalId: ticket.ExternalId,
		PaymentId:  ticket.PaymentId,
	}
}

func AgencyRequestToEntity(agency *request.Agency) entity.Agency {
	return entity.Agency{
		Name:     agency.Name,
		Place:    agency.Place,
		Username: agency.Auth.Username,
		Password: agency.Auth.Password,
	}
}
func BusRequestToEntity(bus *request.Bus) entity.Bus {
	return entity.Bus{

		AgencyId:    bus.AgencyId,
		NumberPlate: bus.NumberPlate,
	}
}

func CustomerRequestToEntity(customer *request.Customer) entity.Customer {
	return entity.Customer{
		Email:       customer.Email,
		Name:        customer.Name,
		PhoneNumber: customer.PhoneNumber,
	}
}

func DriverRequestToEntity(driver *request.Driver) entity.Driver {
	return entity.Driver{

		AgencyId: driver.AgencyId,
		Name:     driver.Name,
	}
}

func TicketRequestToEntity(ticket *request.Ticket) entity.Ticket {
	return entity.Ticket{
		ScheduleId: ticket.ScheduleId,
		CustomerId: ticket.CustomerId,
	}
}

func ScheduleEntityToResponse(schedule *entity.Schedule) response.Schedule {
	return response.Schedule{
		ScheduleId:   schedule.ScheduleId,
		FromAgencyId: schedule.FromAgencyId,
		ToAgencyId:   schedule.ToAgencyId,
		BusId:        schedule.BusId,
		DriverId:     schedule.DriverId,
		Price:        schedule.Price,
		Date:         schedule.Date,
		Arrived:      schedule.Arrived,
	}
}

func ScheduleRequestToEntity(schedule *request.Schedule) entity.Schedule {
	return entity.Schedule{

		FromAgencyId: schedule.FromAgencyId,
		ToAgencyId:   schedule.ToAgencyId,
		BusId:        schedule.BusId,
		DriverId:     schedule.DriverId,
		Price:        schedule.Price,
		Date:         schedule.Date,
		Arrived:      schedule.Arrived,
	}
}
