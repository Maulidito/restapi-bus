package testing

import (
	"context"
	"fmt"
	"restapi-bus/app"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db = app.NewDatabase()
var ctx = context.Background()

func TestCustomerGetAll(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	helper.PanicIfError(err)
	repoCustomer := repository.NewCustomerRepository()

	allCustomer := repoCustomer.GetAllCustomer(ctx, tx)

	assert.Nil(t, err)
	assert.NotNil(t, allCustomer)

}

func TestCustomerAdd(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoCustomer := repository.NewCustomerRepository()
	dataCustomer := &entity.Customer{Name: "TEST123", PhoneNumber: "0832345432"}
	repoCustomer.AddCustomer(ctx, tx, dataCustomer)

	assert.Nil(t, err)
}

func TestCustomerGetOne(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoCustomer := repository.NewCustomerRepository()
	customer := entity.Customer{CustomerId: 1}
	repoCustomer.GetOneCustomer(ctx, tx, &customer)

	fmt.Println(customer)
	assert.NotNil(t, customer)
}

func TestCustomerDelete(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoCustomer := repository.NewCustomerRepository()
	customer := entity.Customer{CustomerId: 3}
	repoCustomer.DeleteOneCustomer(ctx, tx, &customer)

	assert.NotNil(t, customer)
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
