package service

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/request"
	"restapi-bus/models/response"
	"restapi-bus/repository"
)

type CustomerServiceInterface interface {
	GetAllCustomer(ctx context.Context) []response.Customer
	AddCustomer(ctx context.Context, customer *request.Customer) error
	GetOneCustomer(ctx context.Context, id int) response.Customer
	DeleteOneCustomer(ctx context.Context, id int) response.Customer
}

type CustomerServiceImplemtation struct {
	Db   *sql.DB
	Repo repository.CustomerRepositoryInterface
}

func NewCustomerService(db *sql.DB, repo repository.CustomerRepositoryInterface) CustomerServiceInterface {
	return &CustomerServiceImplemtation{Db: db, Repo: repo}
}

func (service *CustomerServiceImplemtation) GetAllCustomer(ctx context.Context) []response.Customer {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	listCustomer := service.Repo.GetAllCustomer(ctx, tx)
	listCustomerResponse := []response.Customer{}

	for _, customer := range listCustomer {
		listCustomerResponse = append(listCustomerResponse, helper.CustomerEntityToResponse(&customer))

	}

	return listCustomerResponse

}
func (service *CustomerServiceImplemtation) AddCustomer(ctx context.Context, customer *request.Customer) error {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	customerEntity := helper.CustomerRequestToEntity(customer)
	err = service.Repo.AddCustomer(ctx, tx, &customerEntity)
	return err
}
func (service *CustomerServiceImplemtation) GetOneCustomer(ctx context.Context, id int) response.Customer {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)
	customerEntity := service.Repo.GetOneCustomer(ctx, tx, id)

	return helper.CustomerEntityToResponse(&customerEntity)

}
func (service *CustomerServiceImplemtation) DeleteOneCustomer(ctx context.Context, id int) response.Customer {
	tx, err := service.Db.Begin()
	helper.PanicIfError(err)

	customerEntity := service.Repo.DeleteOneCustomer(ctx, tx, id)

	return helper.CustomerEntityToResponse(&customerEntity)
}
