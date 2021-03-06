package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type CustomerRepositoryInterface interface {
	GetAllCustomer(ctx context.Context, tx *sql.Tx, filter string) []entity.Customer
	AddCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer)
	GetOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer)
	DeleteOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer)
}

type CustomerRepositoryImplementation struct {
}

func NewCustomerRepository() CustomerRepositoryInterface {
	return &CustomerRepositoryImplementation{}
}

func (repo *CustomerRepositoryImplementation) GetAllCustomer(ctx context.Context, tx *sql.Tx, filter string) []entity.Customer {
	defer helper.ShouldRollback(tx)
	fmt.Println("CHECK FILTER SQL ", filter)
	row, err := tx.QueryContext(ctx, "SELECT customer_id,name,phone_number FROM customer "+filter)
	helper.PanicIfError(err)
	defer row.Close()

	listCustomer := []entity.Customer{}

	for row.Next() {
		tempCustomer := entity.Customer{}
		err := row.Scan(&tempCustomer.CustomerId, &tempCustomer.Name, &tempCustomer.PhoneNumber)
		listCustomer = append(listCustomer, tempCustomer)
		helper.PanicIfError(err)
	}

	return listCustomer
}

func (repo *CustomerRepositoryImplementation) AddCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) {
	defer helper.ShouldRollback(tx)
	res, err := tx.ExecContext(ctx, "Insert Into customer( name , phone_number ) Values (?,?)", customer.Name, customer.PhoneNumber)
	helper.PanicIfError(err)
	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	customer.CustomerId = int(id)

}

func (repo *CustomerRepositoryImplementation) GetOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) {
	defer helper.ShouldRollback(tx)
	rows, err := tx.QueryContext(ctx, "SELECT name, phone_number FROM customer where customer_id = ?", customer.CustomerId)

	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&customer.Name, &customer.PhoneNumber)
		helper.PanicIfError(err)
		return
	}
	panic(exception.NewNotFoundError(fmt.Sprintf("ID Customer %d Not Found", customer.CustomerId)))

}
func (repo *CustomerRepositoryImplementation) DeleteOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) {
	defer helper.ShouldRollback(tx)

	_, err := tx.ExecContext(ctx, "DELETE FROM customer WHERE customer_id = ?", customer.CustomerId)

	helper.PanicIfError(err)

}
