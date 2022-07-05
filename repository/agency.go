package repository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/helper"
	"restapi-bus/models/entity"
)

type AgencyRepositoryInterface interface {
	GetAllAgency(ctx context.Context, tx *sql.Tx) []entity.Agency
	AddAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) error
	GetOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency)
	DeleteOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency)
}

type AgencyRepositoryImplementation struct {
}

func NewAgencyRepository() AgencyRepositoryInterface {
	return &AgencyRepositoryImplementation{}
}

func (repo *AgencyRepositoryImplementation) GetAllAgency(ctx context.Context, tx *sql.Tx) []entity.Agency {

	row, err := tx.QueryContext(ctx, "SELECT agency_id,name,place FROM agency")
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

func (repo *AgencyRepositoryImplementation) AddAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) error {
	res, err := tx.ExecContext(ctx, "Insert Into agency( name , place ) Values (?,?)", agency.Name, agency.Place)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	id, err := res.LastInsertId()

	agency.AgencyId = int(id)

	return err

}

func (repo *AgencyRepositoryImplementation) GetOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) {
	rows, err := tx.QueryContext(ctx, "SELECT name, place FROM agency where agency_id = ?", agency.AgencyId)

	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&agency.Name, &agency.Place)
		helper.PanicIfError(err)
		return
	}
	panic(fmt.Sprintf("ID Agency %d Not Found", agency.AgencyId))

}
func (repo *AgencyRepositoryImplementation) DeleteOneAgency(ctx context.Context, tx *sql.Tx, agency *entity.Agency) {

	repo.GetOneAgency(ctx, tx, agency)
	_, err := tx.ExecContext(ctx, "DELETE FROM agency WHERE agency_id = ?", agency.AgencyId)

	if err != nil {
		tx.Rollback()
		helper.PanicIfError(err)
	}
	tx.Commit()

}
