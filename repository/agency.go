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

type AgencyRepositoryImplementation struct {
	conn *sql.DB
}

func NewAgencyRepository(conn *sql.DB) entity.AgencyRepositoryInterface {
	return &AgencyRepositoryImplementation{conn: conn}
}

func (repo *AgencyRepositoryImplementation) GetAllAgency(ctx context.Context, filter *request.AgencyFilter) []entity.Agency {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	filterString := helper.RequestFilterAgencyToString(filter)
	row, err := tx.QueryContext(ctx, "SELECT agency.agency_id,agency.name,agency.place FROM agency "+filterString)
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

func (repo *AgencyRepositoryImplementation) RegisterAgency(ctx context.Context, agency *entity.Agency) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	res, err := tx.ExecContext(ctx, "Insert Into agency( name , place, username, password, salt ) Values (?,?,?,?,?)", agency.Name, agency.Place, agency.Username, agency.Password, agency.Salt)

	helper.PanicIfError(err)

	id, err := res.LastInsertId()
	helper.PanicIfError(err)

	agency.AgencyId = int(id)

}

func (repo *AgencyRepositoryImplementation) GetOneAgency(ctx context.Context, agency *entity.Agency) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT name, place FROM agency where agency_id = ?", agency.AgencyId).
		Scan(
			&agency.Name,
			&agency.Place,
		)

	if err != nil {
		errMsg := fmt.Sprintf("ID Agency %d Not Found", agency.AgencyId)
		fmt.Println("CHECK GET ONE AGENCY", err)
		panic(exception.NewNotFoundError(errMsg))
	}

}
func (repo *AgencyRepositoryImplementation) DeleteOneAgency(ctx context.Context, agency *entity.Agency) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	_, err = tx.ExecContext(ctx, "DELETE FROM agency WHERE agency_id = ?", agency.AgencyId)

	helper.PanicIfError(err)

}

func (repo *AgencyRepositoryImplementation) GetOneAgencyAuth(ctx context.Context, agency *entity.Agency) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT agency_id,name,place FROM agency where username = ? AND password = ?", agency.Username, agency.Password).
		Scan(
			&agency.AgencyId,
			&agency.Name,
			&agency.Place,
		)

	if err != nil {
		errMsg := fmt.Sprintf("password %s is wrong", agency.Username)

		panic(exception.NewNotFoundError(errMsg))
	}

}

func (repo *AgencyRepositoryImplementation) GetSaltAgencyWithUsername(ctx context.Context, agencyUsername string) (saltResult string) {
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)
	err = tx.QueryRowContext(ctx, "SELECT salt place FROM agency where username = ? ", agencyUsername).Scan(&saltResult)

	if err != nil {
		errMsg := fmt.Sprintf("username %s Not Found", agencyUsername)
		panic(exception.NewNotFoundError(errMsg))
	}
	return
}

func (repo *AgencyRepositoryImplementation) IsUsenameAgencyExist(ctx context.Context, agencyUsername string) bool {

	var name_temp string
	tx, err := repo.conn.Begin()
	defer helper.DoCommitOrRollback(tx)
	helper.PanicIfError(err)

	err = tx.QueryRowContext(ctx, "SELECT name place FROM agency where username = ? ", agencyUsername).Scan(&name_temp)

	return err == nil

}
