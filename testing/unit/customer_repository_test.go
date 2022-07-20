package testing

import (
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var db, mock, rows = func() (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows) {
	db, mock, err := sqlmock.New()
	helper.PanicIfError(err)

	rows := sqlmock.NewRows([]string{"customer_id", "name", "phone_number"})
	rows.AddRow(1, "dito", "08522346789")
	rows.AddRow(2, "Ambay", "0855253469")
	rows.AddRow(3, "Zigay", "08534236880")
	rows.AddRow(4, "Cikay", "08526781230")
	rows.AddRow(5, "Koko", "0834223380")

	return db, mock, rows
}()

var ctx = context.Background()

func TestGetCustomer(t *testing.T) {

}
func TestCustomerGetAll(t *testing.T) {

	mock.ExpectBegin()

	mock.ExpectQuery("SELECT customer_id,name,phone_number FROM customer").WillReturnRows(rows)
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoCustomer := repository.NewCustomerRepository()

	allCustomer := repoCustomer.GetAllCustomer(ctx, tx, "")
	fmt.Println(rows)

	assert.Nil(t, err)
	assert.NotNil(t, allCustomer)

}

func TestCustomerAdd(t *testing.T) {
	customer := &entity.Customer{Name: "TEST123", PhoneNumber: "0832345432"}

	mock.ExpectBegin()
	mock.ExpectExec("Insert Into customer").WithArgs(customer.Name, customer.PhoneNumber).WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoCustomer := repository.NewCustomerRepository()

	repoCustomer.AddCustomer(ctx, tx, customer)

	fmt.Println(customer)
	assert.NotZero(t, customer.CustomerId)
}

func TestCustomerGetOne(t *testing.T) {
	customer := entity.Customer{CustomerId: 1}

	twoRows := sqlmock.NewRows([]string{"name", "phone_number"}).AddRow("DITO", "085234343").AddRow("GITA", "02343435")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT name, phone_number FROM customer ").WithArgs(customer.CustomerId).WillReturnRows(twoRows)

	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoCustomer := repository.NewCustomerRepository()

	repoCustomer.GetOneCustomer(ctx, tx, &customer)
	fmt.Println(customer)
	assert.NotEmpty(t, customer.Name)
}

func TestCustomerDelete(t *testing.T) {
	customer := entity.Customer{CustomerId: 3}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM customer").WithArgs(customer.CustomerId).WillReturnResult(sqlmock.NewResult(3, 1))
	mock.ExpectCommit()
	tx, err := db.Begin()

	helper.PanicIfError(err)
	repoCustomer := repository.NewCustomerRepository()

	repoCustomer.DeleteOneCustomer(ctx, tx, &customer)

	assert.Nil(t, err)
}

func BenchmarkDatabaseTransaction(b *testing.B) {
	fmt.Println("IN MAX OPEN CONN", db.Stats().MaxOpenConnections)

	go func() {

		for i := 0; i < b.N; i++ {
			fmt.Println("IN MAX Idle CLosed", db.Stats().MaxIdleClosed)
			fmt.Println("IN MAX Life Time", db.Stats().MaxLifetimeClosed)
			fmt.Println("IN USE", db.Stats().InUse)
			fmt.Println("IN OpenCon", db.Stats().OpenConnections)
			fmt.Println("IN Idle", db.Stats().Idle)
			TestCustomerGetAll(&testing.T{})
			TestCustomerGetOne(&testing.T{})
		}
	}()

}
