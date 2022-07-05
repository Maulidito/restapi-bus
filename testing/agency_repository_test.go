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
	agency := entity.Agency{AgencyId: 3}
	repoAgency.GetOneAgency(ctx, tx, &agency)

	fmt.Println(agency)
	assert.NotNil(t, agency)
}

func TestAgencyDelete(t *testing.T) {
	tx, err := db.Begin()
	helper.PanicIfError(err)
	repoAgency := repository.NewAgencyRepository()
	agency := entity.Agency{AgencyId: 3}
	repoAgency.DeleteOneAgency(ctx, tx, &agency)

	fmt.Println("Deleted ", agency)
	assert.NotNil(t, agency)
}
