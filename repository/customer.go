package repository

import (
	"context"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
)

var customerRepositorySingleton *CustomerRepositoryImplementation

type CustomerRepositoryImplementation struct {
}

func NewCustomerRepository() entity.CustomerRepositoryInterface {
	if customerRepositorySingleton == nil {
		customerRepositorySingleton = &CustomerRepositoryImplementation{}
	}
	return customerRepositorySingleton
}

func (repo *CustomerRepositoryImplementation) GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []entity.Customer {
	tx := database.GetTransactionContext(ctx)
	filterString := helper.RequestFilterCustomerToString(filter)
	row, err := tx.QueryContext(ctx, "SELECT customer_id,name,phone_number,email FROM customer "+filterString)
	helper.PanicIfError(err)
	defer row.Close()

	listCustomer := []entity.Customer{}

	for row.Next() {
		tempCustomer := entity.Customer{}
		err := row.Scan(&tempCustomer.CustomerId, &tempCustomer.Name, &tempCustomer.PhoneNumber, &tempCustomer.Email)
		listCustomer = append(listCustomer, tempCustomer)
		helper.PanicIfError(err)
	}

	return listCustomer
}

func (repo *CustomerRepositoryImplementation) AddCustomer(ctx context.Context, customer *entity.Customer) {
	tx := database.GetTransactionContext(ctx)
	res, err := tx.ExecContext(ctx, "Insert Into customer( name , phone_number, email ) Values (?,?,?)", customer.Name, customer.PhoneNumber, customer.Email)
	helper.PanicIfError(err)
	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	customer.CustomerId = int(id)

}

func (repo *CustomerRepositoryImplementation) GetOneCustomer(ctx context.Context, customer *entity.Customer) {
	tx := database.GetTransactionContext(ctx)
	err := tx.QueryRowContext(ctx, "SELECT name, phone_number,email FROM customer where customer_id = ?", customer.CustomerId).
		Scan(&customer.Name, &customer.PhoneNumber, &customer.Email)

	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("ID Customer %d Not Found", customer.CustomerId)))
	}

}
func (repo *CustomerRepositoryImplementation) DeleteOneCustomer(ctx context.Context, customer *entity.Customer) {
	tx := database.GetTransactionContext(ctx)
	_, err := tx.ExecContext(ctx, "DELETE FROM customer WHERE customer_id = ?", customer.CustomerId)

	helper.PanicIfError(err)

}
