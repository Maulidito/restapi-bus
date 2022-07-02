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
	GetOneAgency(ctx context.Context, tx *sql.Tx, id int) entity.Agency
	DeleteOneAgency(ctx context.Context, tx *sql.Tx, id int) entity.Agency
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

func (repo *AgencyRepositoryImplementation) GetOneAgency(ctx context.Context, tx *sql.Tx, id int) entity.Agency {
	rows, err := tx.QueryContext(ctx, "SELECT agency_id, name, place FROM agency where agency_id = ?", id)

	helper.PanicIfError(err)
	defer rows.Close()

	agencyData := entity.Agency{}
	if rows.Next() {
		err = rows.Scan(&agencyData.AgencyId, &agencyData.Name, &agencyData.Place)
		helper.PanicIfError(err)
		return agencyData
	}
	panic(fmt.Sprintf("ID Agency %d Not Found", id))

}
func (repo *AgencyRepositoryImplementation) DeleteOneAgency(ctx context.Context, tx *sql.Tx, id int) entity.Agency {

	agencyData := repo.GetOneAgency(ctx, tx, id)
	_, err := tx.ExecContext(ctx, "DELETE FROM agency WHERE agency_id = ?", id)

	if err != nil {
		tx.Rollback()
		helper.PanicIfError(err)
	}
	tx.Commit()

	return agencyData

}
