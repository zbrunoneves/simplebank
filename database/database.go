package database

import (
	"context"
	"database/sql"
	"errors"
)

func ExecTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}
