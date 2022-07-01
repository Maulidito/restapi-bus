package repository

import (
	"context"
	"database/sql"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type CustomerRepositoryInterface interface {
	GetAllCustomer(ctx context.Context, tx *sql.Tx) []entity.Customer
	AddCustomer(ctx context.Context, tx *sql.Tx, customer *entity.Customer) error
	GetOneCustomer(ctx context.Context, tx *sql.Tx, id int) entity.Customer
	DeleteOneCustomer(ctx context.Context, tx *sql.Tx, id int) entity.Customer
}

type CustomerRepositoryImplementation struct {
}

func NewCustomerRepository() CustomerRepositoryImplementation {
	return CustomerRepositoryImplementation{}
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

func (repo *CustomerRepositoryImplementation) GetOneCustomer(ctx context.Context, tx *sql.Tx, id int) entity.Customer {
	rows, err := tx.QueryContext(ctx, "SELECT customer_id, name, phone_number FROM customer where customer_id = ?", id)

	helper.PanicIfError(err)
	defer rows.Close()

	customerData := entity.Customer{}
	if rows.Next() {
		err = rows.Scan(&customerData.CustomerId, &customerData.Name, &customerData.PhoneNumber)
	}

	helper.PanicIfError(err)
	return customerData

}
func (repo *CustomerRepositoryImplementation) DeleteOneCustomer(ctx context.Context, tx *sql.Tx, id int) entity.Customer {

	customerData := repo.GetOneCustomer(ctx, tx, id)
	_, err := tx.ExecContext(ctx, "DELETE FROM customer WHERE customer_id = ?", id)

	if err != nil {
		tx.Rollback()
		helper.PanicIfError(err)
	}
	tx.Commit()

	return customerData

}
