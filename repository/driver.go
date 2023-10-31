package repository

import (
	"context"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/database"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
)

var driverRepositorySingleton *DriverRepositoryImplementation

type DriverRepositoryImplementation struct {
}

func NewDiverRepository() entity.DriverRepositoryInterface {
	if driverRepositorySingleton == nil {
		driverRepositorySingleton = &DriverRepositoryImplementation{}
	}
	return driverRepositorySingleton
}

func (repo *DriverRepositoryImplementation) GetAllDriver(ctx context.Context, filter *request.DriverFilter) []entity.Driver {
	tx := database.GetTransactionContext(ctx)
	filterString := helper.RequestFilterDriverToString(filter)
	rows, err := tx.QueryContext(ctx, "SELECT driver_id,agency_id,name FROM driver "+filterString)
	helper.PanicIfError(err)

	listDriver := []entity.Driver{}
	for rows.Next() {
		setDriver := entity.Driver{}
		err = rows.Scan(&setDriver.DriverId, &setDriver.AgencyId, &setDriver.Name)
		helper.PanicIfError(err)
		listDriver = append(listDriver, setDriver)
	}

	return listDriver

}
func (repo *DriverRepositoryImplementation) GetAllDriverOnSpecificAgency(ctx context.Context, agencyId int) []entity.Driver {
	tx := database.GetTransactionContext(ctx)

	rows, err := tx.QueryContext(ctx, "SELECT driver_id,agency_id,name FROM driver WHERE agency_id = ?", agencyId)
	helper.PanicIfError(err)

	listDriver := []entity.Driver{}
	for rows.Next() {
		setDriver := entity.Driver{}
		err = rows.Scan(&setDriver.DriverId, &setDriver.AgencyId, &setDriver.Name)
		helper.PanicIfError(err)
		listDriver = append(listDriver, setDriver)
	}

	return listDriver

}
func (repo *DriverRepositoryImplementation) GetOneDriverOnSpecificAgency(ctx context.Context, driver *entity.Driver) {
	tx := database.GetTransactionContext(ctx)
	err := tx.QueryRowContext(ctx, "SELECT agency_id,name FROM driver WHERE driver_id = ?", driver.DriverId).
		Scan(&driver.AgencyId, &driver.Name)

	if err != nil {

		panic(exception.NewNotFoundError(fmt.Sprintf("ERROR NOT FOUND DRIVER ID %d", driver.DriverId)))
	}

}
func (repo *DriverRepositoryImplementation) AddDriver(ctx context.Context, driver *entity.Driver) {
	tx := database.GetTransactionContext(ctx)
	res, err := tx.ExecContext(ctx, "INSERT INTO driver(agency_id,name) VALUES (? ,?)", driver.AgencyId, driver.Name)

	helper.PanicIfError(err)

	idDriver, err := res.LastInsertId()
	helper.PanicIfError(err)
	driver.DriverId = int(idDriver)

}
func (repo *DriverRepositoryImplementation) DeleteDriver(ctx context.Context, driver *entity.Driver) {
	tx := database.GetTransactionContext(ctx)
	_, err := tx.ExecContext(ctx, "DELETE FROM driver WHERE driver_id = ?", driver.DriverId)

	helper.PanicIfError(err)

}
