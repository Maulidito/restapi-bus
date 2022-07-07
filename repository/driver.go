package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type DriverRepositoryInterface interface {
	GetAllDriver(tx *sql.Tx, ctx context.Context) []entity.Driver
	GetAllDriverOnSpecificAgency(tx *sql.Tx, ctx context.Context, agencyId int) []entity.Driver
	GetOneDriverOnSpecificAgency(tx *sql.Tx, ctx context.Context, driver *entity.Driver)
	AddDriver(tx *sql.Tx, ctx context.Context, driver *entity.Driver)
	DeleteDriver(tx *sql.Tx, ctx context.Context, driver *entity.Driver)
}

type DriverRepositoryImplementation struct {
}

func NewDiverRepository() DriverRepositoryInterface {
	return &DriverRepositoryImplementation{}
}

func (repo *DriverRepositoryImplementation) GetAllDriver(tx *sql.Tx, ctx context.Context) []entity.Driver {
	defer helper.ShouldRollback(tx)
	rows, err := tx.QueryContext(ctx, "SELECT driver_id,agency_id,name FROM driver")
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
func (repo *DriverRepositoryImplementation) GetAllDriverOnSpecificAgency(tx *sql.Tx, ctx context.Context, agencyId int) []entity.Driver {
	defer helper.ShouldRollback(tx)

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
func (repo *DriverRepositoryImplementation) GetOneDriverOnSpecificAgency(tx *sql.Tx, ctx context.Context, driver *entity.Driver) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT name FROM driver WHERE driver_id = ?", driver.DriverId)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&driver.Name)
		helper.PanicIfError(err)
		return
	}
	panic(fmt.Errorf("ERROR NOT FOUND DRIVER ID %d", driver.DriverId))

}
func (repo *DriverRepositoryImplementation) AddDriver(tx *sql.Tx, ctx context.Context, driver *entity.Driver) {
	defer helper.ShouldRollback(tx)

	res, err := tx.ExecContext(ctx, "INSERT INTO driver(agency_id,name) VALUES (? ,?)", driver.AgencyId, driver.Name)

	helper.PanicIfError(err)

	idDriver, err := res.LastInsertId()
	helper.PanicIfError(err)
	driver.DriverId = int(idDriver)

}
func (repo *DriverRepositoryImplementation) DeleteDriver(tx *sql.Tx, ctx context.Context, driver *entity.Driver) {
	defer helper.ShouldRollback(tx)

	_, err := tx.ExecContext(ctx, "DELETE FROM driver WHERE driver_id = ?", driver.DriverId)

	helper.PanicIfError(err)

}
