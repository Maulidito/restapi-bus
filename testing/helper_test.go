package testing

import (
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestToEntitySuccessAgency(t *testing.T) {
	agency := request.Agency{Name: "TEST", Place: "TEST123"}

	agencyEntity, err := helper.RequestToEntity[*request.Agency, entity.Agency](&agency)

	assert.Nil(t, err)
	assert.Equal(t, agency.Name, agencyEntity.Name)
	assert.Equal(t, agency.Place, agencyEntity.Place)
	assert.Zero(t, agencyEntity.AgencyId)

}
func TestRequestToEntitySuccessBus(t *testing.T) {
	bus := request.Bus{AgencyId: 1, NumberPlate: "TEST123"}

	busEntity, err := helper.RequestToEntity[*request.Bus, entity.Bus](&bus)

	assert.Nil(t, err)
	assert.Equal(t, bus.AgencyId, busEntity.AgencyId)
	assert.Equal(t, bus.NumberPlate, busEntity.NumberPlate)
	assert.Zero(t, busEntity.BusId)

}

func TestRequestToEntityFailed(t *testing.T) {
	agency := request.Agency{Name: "TEST", Place: "TEST123"}

	agencyEntity, err := helper.RequestToEntity[*request.Agency, entity.Bus](&agency)

	assert.NotNil(t, err)
	assert.Equal(t, agencyEntity, entity.Bus{})

}
