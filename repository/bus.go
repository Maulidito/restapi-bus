package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type BusRepositoryInterface interface {
	GetAllBus(ctx context.Context, tx *sql.Tx) []entity.Bus
	AddBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus) error
	GetOneBus(ctx context.Context, tx *sql.Tx, id int) entity.Bus
	DeleteOneBus(ctx context.Context, tx *sql.Tx, id int) entity.Bus
}

type BusRepositoryImplementation struct {
}

func NewBusRepository() BusRepositoryInterface {
	return &BusRepositoryImplementation{}
}

func (repo *BusRepositoryImplementation) GetAllBus(ctx context.Context, tx *sql.Tx) []entity.Bus {

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

func (repo *BusRepositoryImplementation) AddBus(ctx context.Context, tx *sql.Tx, bus *entity.Bus) error {
	res, err := tx.ExecContext(ctx, "Insert Into bus( agency_id , number_plate ) Values (?,?)", bus.AgencyId, bus.NumberPlate)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	id, err := res.LastInsertId()

	bus.BusId = int(id)

	return err

}

func (repo *BusRepositoryImplementation) GetOneBus(ctx context.Context, tx *sql.Tx, id int) entity.Bus {
	rows, err := tx.QueryContext(ctx, "SELECT bus_id, agency_id, number_plate FROM bus where bus_id = ?", id)

	helper.PanicIfError(err)
	defer rows.Close()

	busData := entity.Bus{}
	if rows.Next() {
		err = rows.Scan(&busData.BusId, &busData.AgencyId, &busData.NumberPlate)
		helper.PanicIfError(err)
		return busData
	}
	panic(fmt.Sprintf("ID Bus %d Not Found", id))

}
func (repo *BusRepositoryImplementation) DeleteOneBus(ctx context.Context, tx *sql.Tx, id int) entity.Bus {

	busData := repo.GetOneBus(ctx, tx, id)
	_, err := tx.ExecContext(ctx, "DELETE FROM bus WHERE bus_id = ?", id)

	if err != nil {
		tx.Rollback()
		helper.PanicIfError(err)
	}
	tx.Commit()

	return busData

}
