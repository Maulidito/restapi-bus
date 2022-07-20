package testing

import (
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllBus(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repo := repository.NewBusRepository()
	listBus := repo.GetAllBus(ctx, tx, "")

	fmt.Println(listBus)
	assert.NotEmpty(t, listBus)

}

func TestGetAllBusOnSpecificAgency(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repo := repository.NewBusRepository()
	bus := repo.GetAllBusSpecificAgency(ctx, tx, 1)

	fmt.Println(bus)
	assert.NotEmpty(t, bus)

}

func TestGetOneBusOnSpecificAgency(t *testing.T) {
	bus := entity.Bus{
		BusId:    2,
		AgencyId: 1,
	}
	defer func() {
		err := recover()
		fmt.Println(bus, err)

		assert.NotEmpty(t, err)
		assert.Empty(t, bus)
	}()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repo := repository.NewBusRepository()

	repo.GetOneBus(ctx, tx, &bus)

}

func TestAddBus(t *testing.T) {
	bus := entity.Bus{
		NumberPlate: "T 8923 XCF",
		AgencyId:    2,
	}
	defer func() {
		err := recover()
		fmt.Println(bus, err)

		assert.NotEmpty(t, err)
		assert.Empty(t, bus)
	}()

	tx, err := db.Begin()
	helper.PanicIfError(err)
	repo := repository.NewBusRepository()

	repo.AddBus(ctx, tx, &bus)

}
