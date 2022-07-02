package testing

import (
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgencyGetAll(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoAgency := repository.NewAgencyRepository()

	allAgency := repoAgency.GetAllAgency(ctx, tx)
	fmt.Println(allAgency)
	assert.Nil(t, err)
	assert.NotNil(t, allAgency)

}

func TestAgencyAdd(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)

	repoAgency := repository.NewAgencyRepository()
	dataAgency := &entity.Agency{Name: "Pariwisata Raya", Place: "Jakarta"}
	err = repoAgency.AddAgency(ctx, tx, dataAgency)

	assert.Nil(t, err)
}

func TestAgencyGetOne(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoAgency := repository.NewAgencyRepository()
	AgencyData := repoAgency.GetOneAgency(ctx, tx, 3)

	fmt.Println(AgencyData)
	assert.NotNil(t, AgencyData)
}

func TestAgencyDelete(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoAgency := repository.NewAgencyRepository()
	AgencyData := repoAgency.DeleteOneAgency(ctx, tx, 3)

	fmt.Println("Deleted ", AgencyData)
	assert.NotNil(t, AgencyData)
}
