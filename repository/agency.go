package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/exception"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type AgencyRepositoryInterface interface {
	GetAllAgency(ctx context.Context, tx *sql.Tx, filter string) []entity.Agency
	AddAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency)
	GetOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency)
	DeleteOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency)
}

type AgencyRepositoryImplementation struct {
}

func NewAgencyRepository() AgencyRepositoryInterface {
	return &AgencyRepositoryImplementation{}
}

func (repo *AgencyRepositoryImplementation) GetAllAgency(ctx context.Context, tx *sql.Tx, filter string) []entity.Agency {
	defer helper.ShouldRollback(tx)

	row, err := tx.QueryContext(ctx, "SELECT agency.agency_id,agency.name,agency.place FROM agency "+filter)
	helper.PanicIfError(err)

	defer row.Close()

	listAgency := []entity.Agency{}

	for row.Next() {
		tempAgency := entity.Agency{}
		err := row.Scan(&tempAgency.AgencyId, &tempAgency.Name, &tempAgency.Place)
		listAgency = append(listAgency, tempAgency)
		helper.PanicIfError(err)
	}

	return listAgency
}

func (repo *AgencyRepositoryImplementation) AddAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) {
	defer helper.ShouldRollback(tx)
	res, err := tx.ExecContext(ctx, "Insert Into agency( name , place ) Values (?,?)", agency.Name, agency.Place)

	helper.PanicIfError(err)

	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	agency.AgencyId = int(id)

}

func (repo *AgencyRepositoryImplementation) GetOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) {
	defer helper.ShouldRollback(tx)

	rows, err := tx.QueryContext(ctx, "SELECT name, place FROM agency where agency_id = ?", agency.AgencyId)

	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&agency.Name, &agency.Place)
		helper.PanicIfError(err)
		return
	}

	errMsg := fmt.Sprintf("ID Agency %d Not Found", agency.AgencyId)

	panic(exception.NewNotFoundError(errMsg))

}
func (repo *AgencyRepositoryImplementation) DeleteOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) {
	defer helper.ShouldRollback(tx)

	_, err := tx.ExecContext(ctx, "DELETE FROM agency WHERE agency_id = ?", agency.AgencyId)

	helper.PanicIfError(err)

}
