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

	agencyEntity := helper.RequestToEntity[*request.Agency, entity.Agency](&agency)

	assert.Equal(t, agency.Name, agencyEntity.Name)
	assert.Equal(t, agency.Place, agencyEntity.Place)
	assert.Zero(t, agencyEntity.AgencyId)

}
func TestRequestToEntitySuccessBus(t *testing.T) {
	bus := request.Bus{AgencyId: 1, NumberPlate: "TEST123"}

	busEntity := helper.RequestToEntity[*request.Bus, entity.Bus](&bus)

	assert.Equal(t, bus.AgencyId, busEntity.AgencyId)
	assert.Equal(t, bus.NumberPlate, busEntity.NumberPlate)
	assert.Zero(t, busEntity.BusId)

}

func TestRequestToEntityFailed(t *testing.T) {
	agency := request.Agency{Name: "TEST", Place: "TEST123"}

	agencyEntity := helper.RequestToEntity[*request.Agency, entity.Bus](&agency)

	assert.Equal(t, agencyEntity, entity.Bus{})

}

func BenchmarkRequestToEntity(b *testing.B) {

	b.Run(
		"Custom Function", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := &request.Agency{Name: "MAULIDITO DWINANDANA", Place: "DEPOK"}
				helper.RequestToEntity[*request.Agency, entity.Agency](data)
			}
		},
	)

	b.Run(
		"Normal Function", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := &request.Agency{Name: "MAULIDITO DWINANDANA", Place: "DEPOK"}
				helper.AgencyRequestToEntity(data)
			}
		},
	)

}
