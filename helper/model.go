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

func CustomerRequestToEntity(customer *request.Customer) entity.Customer {
	return entity.Customer{
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

func AgencyRequestToEntity(agency *request.Agency) entity.Agency {
	return entity.Agency{
		Name:  agency.Name,
		Place: agency.Place,
	}
}
