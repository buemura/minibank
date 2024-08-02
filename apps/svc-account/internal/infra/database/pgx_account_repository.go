package database

import (
	"context"
	"errors"

	"github.com/buemura/minibank/svc-account/internal/domain/account"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxOrderRepository struct {
	db *pgxpool.Pool
}

func NewPgxOrderRepository() *PgxOrderRepository {
	return &PgxOrderRepository{
		db: Conn,
	}
}

func (r *PgxOrderRepository) FindById(id string) (*account.Account, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, balance, owner_name, owner_document, status, created_at, updated_at
		FROM "accounts" 
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	acc, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[account.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, err
	}
	return acc, nil
}
