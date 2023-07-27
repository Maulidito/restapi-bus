package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
)

type CustomerRepositoryImplementation struct {
	conn *sql.DB
}

func NewCustomerRepository(conn *sql.DB) entity.CustomerRepositoryInterface {
	return &CustomerRepositoryImplementation{conn: conn}
}

func (repo *CustomerRepositoryImplementation) GetAllCustomer(ctx context.Context, filter *request.CustomerFilter) []entity.Customer {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
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
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	res, err := tx.ExecContext(ctx, "Insert Into customer( name , phone_number, email ) Values (?,?,?)", customer.Name, customer.PhoneNumber, customer.Email)
	helper.PanicIfError(err)
	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	customer.CustomerId = int(id)

}

func (repo *CustomerRepositoryImplementation) GetOneCustomer(ctx context.Context, customer *entity.Customer) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT name, phone_number,email FROM customer where customer_id = ?", customer.CustomerId).
		Scan(&customer.Name, &customer.PhoneNumber, &customer.Email)

	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("ID Customer %d Not Found", customer.CustomerId)))
	}

}
func (repo *CustomerRepositoryImplementation) DeleteOneCustomer(ctx context.Context, customer *entity.Customer) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	_, err = tx.ExecContext(ctx, "DELETE FROM customer WHERE customer_id = ?", customer.CustomerId)

	helper.PanicIfError(err)

}
