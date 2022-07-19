package helper

import (
	"database/sql"
)

func ShouldRollback(tx *sql.Tx) {
	err := recover()

	if err != nil {

		errRollback := tx.Rollback()
		PanicIfError(errRollback)

		panic(err)
	}

}

func DoCommit(tx *sql.Tx) {
	err := recover()
	if err == nil {
		err := tx.Commit()

		PanicIfError(err)
		return
	}
	panic(err)
}
