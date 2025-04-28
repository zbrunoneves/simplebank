package db

import (
	"database/sql"

	"simplebank/db/sqlc"
)

type Store struct {
	*sqlc.Queries
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		Queries: sqlc.New(db),
		db:      db,
	}
}
