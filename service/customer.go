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

type CustomerServiceInterface interface {
	GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []response.Customer
	AddCustomer(ctx context.Context, customer *request.Customer)
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

func (service *CustomerServiceImplemtation) GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []response.Customer {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	listCustomer := service.Repo.GetAllCustomer(ctx, tx, helper.RequestFilterCustomerToString(filter))
	listCustomerResponse := []response.Customer{}

	for _, customer := range listCustomer {
		listCustomerResponse = append(listCustomerResponse, helper.CustomerEntityToResponse(&customer))

	}

	return listCustomerResponse

}
func (service *CustomerServiceImplemtation) AddCustomer(ctx context.Context, customer *request.Customer) {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	customerEntity := helper.CustomerRequestToEntity(customer)
	service.Repo.AddCustomer(ctx, tx, &customerEntity)

}
func (service *CustomerServiceImplemtation) GetOneCustomer(ctx context.Context, id int) response.Customer {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	customer := entity.Customer{CustomerId: id}
	service.Repo.GetOneCustomer(ctx, tx, &customer)

	return helper.CustomerEntityToResponse(&customer)

}
func (service *CustomerServiceImplemtation) DeleteOneCustomer(ctx context.Context, id int) response.Customer {
	tx, err := service.Db.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	customer := entity.Customer{CustomerId: id}
	service.Repo.GetOneCustomer(ctx, tx, &customer)
	service.Repo.DeleteOneCustomer(ctx, tx, &customer)

	return helper.CustomerEntityToResponse(&customer)
}
