package helper

import (
	"errors"
	"fmt"
	"log"
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

		BusId:          ticket.BusId,
		DriverId:       ticket.DriverId,
		CustomerId:     ticket.CustomerId,
		DeparturePlace: ticket.DeparturePlace,
		ArrivalPlace:   ticket.ArrivalPlace,
		Price:          ticket.Price,
		Date:           ticket.Date,
	}
}

func RequestToEntity[REQ interface {
	*request.Bus | *request.Agency | *request.Customer | *request.Driver | *request.Ticket
},
	ENT interface {
		entity.Bus | entity.Agency | entity.Customer | entity.Driver | entity.Ticket
	},
](requestInput REQ) (dataReturn ENT) {

	defer func() {
		data := recover()
		if data != nil {
			err := errors.New(fmt.Sprint(data))
			log.Fatal(err)
		}

	}()

	dataCustomer, isCustomer := any(requestInput).(*request.Customer)
	dataAgency, isAgency := any(requestInput).(*request.Agency)
	dataDriver, isDriver := any(requestInput).(*request.Driver)
	dataTicket, isTicket := any(requestInput).(*request.Ticket)
	dataBus, isBus := any(requestInput).(*request.Bus)

	switch {
	case isCustomer:
		{
			dataReturn = any(entity.Customer{
				Name:        dataCustomer.Name,
				PhoneNumber: dataCustomer.PhoneNumber,
			}).(ENT)
		}
	case isAgency:
		{
			dataReturn = any(entity.Agency{
				Name:  dataAgency.Name,
				Place: dataAgency.Place,
			}).(ENT)

		}
	case isDriver:
		{
			dataReturn = any(entity.Driver{
				AgencyId: dataDriver.AgencyId,
				Name:     dataDriver.Name,
			}).(ENT)

		}
	case isTicket:
		{
			dataReturn = any(entity.Ticket{
				BusId:          dataTicket.BusId,
				DriverId:       dataTicket.DriverId,
				CustomerId:     dataTicket.CustomerId,
				DeparturePlace: dataTicket.DeparturePlace,
				ArrivalPlace:   dataTicket.DeparturePlace,
				Price:          dataTicket.Price,
				Date:           dataTicket.Date,
			}).(ENT)

		}
	case isBus:
		{
			dataReturn = any(entity.Bus{
				AgencyId:    dataBus.AgencyId,
				NumberPlate: dataBus.NumberPlate,
			}).(ENT)

		}

	}

	return dataReturn

}
