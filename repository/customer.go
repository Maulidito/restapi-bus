package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type CustomerRepositoryInterface interface {
	GetAllCustomer(ctx context.Context, tx *sql.Tx) []entity.Customer
	AddCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) error
	GetOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer)
	DeleteOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer)
}

type CustomerRepositoryImplementation struct {
}

func NewCustomerRepository() CustomerRepositoryInterface {
	return &CustomerRepositoryImplementation{}
}

func (repo *CustomerRepositoryImplementation) GetAllCustomer(ctx context.Context, tx *sql.Tx) []entity.Customer {

	row, err := tx.QueryContext(ctx, "SELECT customer_id,name,phone_number FROM customer")
	helper.PanicIfError(err)
	defer row.Close()

	listCustomer := []entity.Customer{}

	for row.Next() {
		tempCustomer := entity.Customer{}
		err := row.Scan(&tempCustomer.CustomerId, &tempCustomer.Name, &tempCustomer.PhoneNumber)
		listCustomer = append(listCustomer, tempCustomer)
		helper.PanicIfError(err)
	}

	if err = tx.Commit(); err != nil {
		helper.PanicIfError(err)
	}

	return listCustomer
}

func (repo *CustomerRepositoryImplementation) AddCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) error {
	res, err := tx.ExecContext(ctx, "Insert Into customer( name , phone_number ) Values (?,?)", customer.Name, customer.PhoneNumber)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	id, err := res.LastInsertId()

	customer.CustomerId = int(id)

	return err

}

func (repo *CustomerRepositoryImplementation) GetOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) {
	rows, err := tx.QueryContext(ctx, "SELECT name, phone_number FROM customer where customer_id = ?", customer.CustomerId)

	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&customer.Name, &customer.PhoneNumber)
		helper.PanicIfError(err)
		if err = tx.Commit(); err != nil {
			helper.PanicIfError(err)
		}
		return
	}
	panic(fmt.Sprintf("ID Customer %d Not Found", customer.CustomerId))

}
func (repo *CustomerRepositoryImplementation) DeleteOneCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) {

	repo.GetOneCustomer(ctx, tx, customer)
	_, err := tx.ExecContext(ctx, "DELETE FROM customer WHERE customer_id = ?", customer.CustomerId)

	if err != nil {
		tx.Rollback()
		helper.PanicIfError(err)
	}
	tx.Commit()

}
