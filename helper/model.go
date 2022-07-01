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
