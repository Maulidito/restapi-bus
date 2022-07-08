package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type BusRepositoryInterface interface {
	GetAllBus(ctx context.Context, tx *sql.Tx) []entity.Bus
	AddBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus)
	GetOneBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus)
	DeleteOneBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus)
	GetAllBusSpecificAgency(ctx context.Context, tx *sql.Tx, agencyId int) []entity.Bus
}

type BusRepositoryImplementation struct {
}

func NewBusRepository() BusRepositoryInterface {
	return &BusRepositoryImplementation{}
}

func (repo *BusRepositoryImplementation) GetAllBus(ctx context.Context, tx *sql.Tx) []entity.Bus {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT bus_id,agency_id,number_plate FROM bus")
	helper.PanicIfError(err)
	defer row.Close()
	listBus := []entity.Bus{}

	for row.Next() {
		tempBus := entity.Bus{}
		err := row.Scan(&tempBus.BusId, &tempBus.AgencyId, &tempBus.NumberPlate)
		listBus = append(listBus, tempBus)
		helper.PanicIfError(err)
	}

	return listBus
}

func (repo *BusRepositoryImplementation) GetAllBusSpecificAgency(ctx context.Context, tx *sql.Tx, agencyId int) []entity.Bus {
	defer helper.ShouldRollback(tx)
	row, err := tx.QueryContext(ctx, "SELECT bus_id,agency_id,number_plate FROM bus WHERE agency_id = ?", agencyId)
	helper.PanicIfError(err)
	defer row.Close()
	listBus := []entity.Bus{}

	for row.Next() {
		tempBus := entity.Bus{}
		err := row.Scan(&tempBus.BusId, &tempBus.AgencyId, &tempBus.NumberPlate)
		listBus = append(listBus, tempBus)
		helper.PanicIfError(err)
	}

	return listBus
}

func (repo *BusRepositoryImplementation) AddBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus) {
	defer helper.ShouldRollback(tx)
	res, err := tx.ExecContext(ctx, "Insert Into bus( agency_id , number_plate ) Values (?,?)", bus.AgencyId, bus.NumberPlate)

	helper.PanicIfError(err)

	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	bus.BusId = int(id)

}

func (repo *BusRepositoryImplementation) GetOneBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus) {
	defer helper.ShouldRollback(tx)
	rows, err := tx.QueryContext(ctx, "SELECT agency_id,number_plate FROM bus where bus_id = ?", bus.BusId)

	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&bus.AgencyId, &bus.NumberPlate)
		helper.PanicIfError(err)
		return
	}

	panic(exception.NewNotFoundError(fmt.Sprintf("id bus %d not found ", bus.BusId)))

}
func (repo *BusRepositoryImplementation) DeleteOneBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus) {
	defer helper.ShouldRollback(tx)
	repo.GetOneBus(ctx, tx, bus)
	_, err := tx.ExecContext(ctx, "DELETE FROM bus WHERE bus_id = ?", bus.BusId)

	helper.PanicIfError(err)

}
