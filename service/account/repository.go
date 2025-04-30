package account

import (
	"context"
	"database/sql"
	"os"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type CreateAccountParams struct {
	Owner    string
	Currency string
}

func (r *Repository) CreateAccount(params CreateAccountParams) (int, error) {
	qry, err := os.ReadFile("database/sql/account_create.sql")
	if err != nil {
		return -1, err
	}

	res, err := r.db.ExecContext(context.Background(), string(qry), params.Owner, params.Currency)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

type ListAccountsParams struct {
	Limit  int
	Offset int
}

func (r *Repository) ListAccounts(params ListAccountsParams) ([]Account, error) {
	qry, err := os.ReadFile("database/sql/account_list.sql")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(context.Background(), string(qry), params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	items := []Account{}
	for rows.Next() {
		var a Account
		err = rows.Scan(
			&a.ID,
			&a.Owner,
			&a.Balance,
			&a.Currency,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
