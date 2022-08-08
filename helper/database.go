package helper

import (
	"database/sql"
)

func DoCommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err == nil {
		err := tx.Commit()

		PanicIfError(err)
		return
	}
	errRollback := tx.Rollback()
	PanicIfError(errRollback)
	panic(err)
}
