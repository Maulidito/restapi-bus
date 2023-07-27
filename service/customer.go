package service

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
)

type CustomerServiceImplemtation struct {
	Repo entity.CustomerRepositoryInterface
}

func NewCustomerService(repo entity.CustomerRepositoryInterface) entity.CustomerServiceInterface {
	return &CustomerServiceImplemtation{Repo: repo}
}

func (service *CustomerServiceImplemtation) GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []response.Customer {

	listCustomer := service.Repo.GetAllCustomer(ctx, filter)
	listCustomerResponse := []response.Customer{}

	for _, customer := range listCustomer {
		listCustomerResponse = append(listCustomerResponse, helper.CustomerEntityToResponse(&customer))

	}
	return listCustomerResponse

}
func (service *CustomerServiceImplemtation) AddCustomer(ctx context.Context, customer *request.Customer) {

	customerEntity := helper.CustomerRequestToEntity(customer)
	service.Repo.AddCustomer(ctx, &customerEntity)

}
func (service *CustomerServiceImplemtation) GetOneCustomer(ctx context.Context, id int) response.Customer {

	customer := entity.Customer{CustomerId: id}
	service.Repo.GetOneCustomer(ctx, &customer)

	return helper.CustomerEntityToResponse(&customer)

}
func (service *CustomerServiceImplemtation) DeleteOneCustomer(ctx context.Context, id int) response.Customer {

	customer := entity.Customer{CustomerId: id}
	service.Repo.GetOneCustomer(ctx, &customer)
	service.Repo.DeleteOneCustomer(ctx, &customer)

	return helper.CustomerEntityToResponse(&customer)
}
