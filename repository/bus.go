package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
	"restapi-bus/models/request"
)

var busRepositorySingleton *BusRepositoryImplementation

type BusRepositoryImplementation struct {
	conn *sql.DB
}

func NewBusRepository(conn *sql.DB) entity.BusRepositoryInterface {
	if busRepositorySingleton == nil {
		busRepositorySingleton = &BusRepositoryImplementation{conn: conn}
	}
	return busRepositorySingleton
}

func (repo *BusRepositoryImplementation) GetAllBus(ctx context.Context, filter *request.BusFilter) []entity.Bus {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	filterString := helper.RequestFilterBusToString(filter)
	row, err := tx.QueryContext(ctx, "SELECT bus_id,agency_id,number_plate,total_seat FROM bus "+filterString)
	helper.PanicIfError(err)
	defer row.Close()
	listBus := []entity.Bus{}

	for row.Next() {
		tempBus := entity.Bus{}
		err := row.Scan(&tempBus.BusId, &tempBus.AgencyId, &tempBus.NumberPlate, &tempBus.TotalSeat)
		listBus = append(listBus, tempBus)
		helper.PanicIfError(err)
	}

	return listBus
}

func (repo *BusRepositoryImplementation) GetAllBusSpecificAgency(ctx context.Context, agencyId int) []entity.Bus {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	row, err := tx.QueryContext(ctx, "SELECT bus_id,agency_id,number_plate,total_seat FROM bus WHERE agency_id = ?", agencyId)

	helper.PanicIfError(err)
	defer row.Close()
	listBus := []entity.Bus{}

	for row.Next() {
		tempBus := entity.Bus{}
		err := row.Scan(&tempBus.BusId, &tempBus.AgencyId, &tempBus.NumberPlate, &tempBus.TotalSeat)
		listBus = append(listBus, tempBus)
		helper.PanicIfError(err)
	}

	return listBus
}

func (repo *BusRepositoryImplementation) AddBus(ctx context.Context, bus *entity.Bus) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	res, err := tx.ExecContext(ctx, "Insert Into bus( agency_id , number_plate, total_seat ) Values (?,?,?)", bus.AgencyId, bus.NumberPlate, bus.TotalSeat)

	helper.PanicIfError(err)

	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	bus.BusId = int(id)

}

func (repo *BusRepositoryImplementation) GetOneBus(ctx context.Context, bus *entity.Bus) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT agency_id,number_plate,total_seat FROM bus where bus_id = ?", bus.BusId).
		Scan(&bus.AgencyId, &bus.NumberPlate, &bus.TotalSeat)

	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("id bus %d not found ", bus.BusId)))
	}

}
func (repo *BusRepositoryImplementation) DeleteOneBus(ctx context.Context, bus *entity.Bus) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	repo.GetOneBus(ctx, bus)
	_, err = tx.ExecContext(ctx, "DELETE FROM bus WHERE bus_id = ?", bus.BusId)

	helper.PanicIfError(err)

}
