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
		TicketId:       ticket.TicketId,
		AgencyId:       ticket.AgencyId,
		BusId:          ticket.BusId,
		DriverId:       ticket.DriverId,
		CustomerId:     ticket.CustomerId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Date:           ticket.Date,
		Arrived:        ticket.Arrived,
	}
}

func AgencyRequestToEntity(agency *request.Agency) entity.Agency {
	return entity.Agency{

		Name:  agency.Name,
		Place: agency.Place,
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
		AgencyId:       ticket.AgecnyId,
		BusId:          ticket.BusId,
		DriverId:       ticket.DriverId,
		CustomerId:     ticket.CustomerId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Arrived:        ticket.Arrived,
	}
}

func TicketEntityToResponseTicketNoBus(ticket *entity.Ticket) response.TicketNoBus {
	return response.TicketNoBus{
		TicketId:       ticket.AgencyId,
		AgencyId:       ticket.AgencyId,
		DriverId:       ticket.DriverId,
		CustomerId:     ticket.CustomerId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Date:           ticket.Date,
		Arrived:        ticket.Arrived,
	}
}

func TicketEntityToResponseTicketNoDriver(ticket *entity.Ticket) response.TicketNoDriver {
	return response.TicketNoDriver{
		TicketId:       ticket.TicketId,
		BusId:          ticket.BusId,
		AgencyId:       ticket.AgencyId,
		CustomerId:     ticket.CustomerId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Arrived:        ticket.Arrived,
		Date:           ticket.Date,
	}
}

func TicketEntityToResponseTicketNoAgency(ticket *entity.Ticket) response.TicketNoAgency {
	return response.TicketNoAgency{
		TicketId:       ticket.TicketId,
		BusId:          ticket.BusId,
		DriverId:       ticket.DriverId,
		CustomerId:     ticket.CustomerId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Arrived:        ticket.Arrived,
		Date:           ticket.Date,
	}
}

func TicketEntityToResponseTicketNoCustomer(ticket *entity.Ticket) response.TicketNoCustomer {
	return response.TicketNoCustomer{
		TicketId:       ticket.TicketId,
		AgencyId:       ticket.AgencyId,
		BusId:          ticket.BusId,
		DriverId:       ticket.DriverId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Arrived:        ticket.Arrived,
		Date:           ticket.Date,
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
