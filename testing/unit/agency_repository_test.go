package testing

import (
	"context"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
	"restapi-bus/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAgencyGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"agency_id", "name", "place"}).
		AddRow(1, "Agency Jakarta", "Jakarta").
		AddRow(2, "Agency Depok", "Depok").
		AddRow(3, "Agency Bandung", "Bandung")
	mock.ExpectQuery("SELECT agency.agency_id,agency.name,agency.place FROM agency").WillReturnRows(rows)
	mock.ExpectCommit()
	helper.PanicIfError(err)

	repoAgency := repository.NewAgencyRepository(db)

	allAgency := repoAgency.GetAllAgency(context.Background(), &request.AgencyFilter{})

	assert.Nil(t, err)
	assert.NotNil(t, allAgency)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error unfulfilled expectations %s", err)
	}

}

func TestAgencyAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.PanicIfError(err)
	dataAgency := &entity.Agency{Name: "Pariwisata Raya", Place: "Jakarta",
		Username: "bambang", Password: "password", Salt: "1231232"}

	mock.ExpectBegin()
	mock.ExpectExec("Insert Into agency").
		WithArgs(dataAgency.Name, dataAgency.Place, dataAgency.Username, dataAgency.Password, dataAgency.Salt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repoAgency := repository.NewAgencyRepository(db)

	repoAgency.RegisterAgency(context.Background(), dataAgency)

	assert.Nil(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error unfulfilled expectations %s", err)
	}
}

func TestAgencyGetOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	rows := sqlmock.NewRows([]string{"name", "place"}).
		AddRow("Agency Jakarta", "Jakarta")
	helper.PanicIfError(err)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT name, place FROM agency where agency_id = ?").WithArgs(1).WillReturnRows(rows)
	mock.ExpectCommit()
	repoAgency := repository.NewAgencyRepository(db)
	agency := entity.Agency{AgencyId: 1}
	repoAgency.GetOneAgency(context.Background(), &agency)

	assert.NotNil(t, agency)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error unfulfilled expectations %s", err)
	}

}

func TestAgencyDelete(t *testing.T) {
	// tx, err := db.Begin()
	// helper.PanicIfError(err)
	// repoAgency := repository.NewAgencyRepository()
	// agency := entity.Agency{AgencyId: 3}
	// repoAgency.DeleteOneAgency(ctx, tx, &agency)

	// fmt.Println("Deleted ", agency)
	// assert.NotNil(t, agency)
}
