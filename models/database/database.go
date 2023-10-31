package database

import (
	"context"
	"database/sql"
	"fmt"
	"restapi-bus/helper"
)

type tx_key string

const key = "TXKEY"

type Tr struct {
	Db *sql.DB
}

type TrInterface interface {
	BeginTransactionWithContext(ctx context.Context) context.Context
	DoCommitOrRollbackWithContext(ctx context.Context)
	GetTransactionContext(ctx context.Context) (any, error)
}

func (tr *Tr) BeginTransactionWithContext(ctx context.Context) context.Context {
	tx, err := tr.Db.BeginTx(ctx, nil)
	helper.PanicIfError(err)
	return context.WithValue(ctx, tx_key(key), tx)
}

func (tr *Tr) DoCommitOrRollbackWithContext(ctx context.Context) {
	tx, ok := ctx.Value(tx_key(key)).(*sql.Tx)
	if !ok {
		panic(fmt.Errorf("something went wrong with context"))
	}
	err := recover()
	if err == nil {
		err := tx.Commit()

		helper.PanicIfError(err)
		return
	}
	errRollback := tx.Rollback()
	helper.PanicIfError(errRollback)
	panic(err)
}

func (tr *Tr) GetTransactionContext(ctx context.Context) (any, error) {
	value := ctx.Value(tx_key(key))
	if value == nil {
		return nil, fmt.Errorf("something went wrong with context")
	}
	return value, nil
}

func GetTxKey() tx_key {
	return tx_key(key)
}

func GetTransactionContext(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(GetTxKey()).(*sql.Tx)
	if !ok {
		panic("something went wrong with ctx")
	}
	return tx

}
