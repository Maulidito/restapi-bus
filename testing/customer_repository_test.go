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
	fmt.Println(allCustomer)
	assert.Nil(t, err)
	assert.NotNil(t, allCustomer)

}

func TestCustomerAdd(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoCustomer := repository.NewCustomerRepository()
	dataCustomer := &entity.Customer{Name: "TEST123", PhoneNumber: "0832345432"}
	err = repoCustomer.AddCustomer(ctx, tx, dataCustomer)

	assert.Nil(t, err)
}

func TestCustomerGetOne(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoCustomer := repository.NewCustomerRepository()
	customerData := repoCustomer.GetOneCustomer(ctx, tx, 3)

	fmt.Println(customerData)
	assert.NotNil(t, customerData)
}

func TestCustomerDelete(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoCustomer := repository.NewCustomerRepository()
	customerData := repoCustomer.DeleteOneCustomer(ctx, tx, 3)

	fmt.Println("Deleted ", customerData)
	assert.NotNil(t, customerData)
}
