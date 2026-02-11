package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/buemura/minibank/svc-account/internal/domain"
	"github.com/buemura/minibank/svc-account/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type accountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) repository.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (user_id, account_number, agency, account_type, balance, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		account.UserID,
		account.AccountNumber,
		account.Agency,
		account.AccountType,
		account.Balance,
		account.Currency,
		account.Status,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `
		SELECT id, user_id, account_number, agency, account_type, balance, currency, status, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	account := &domain.Account{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.Agency,
		&account.AccountType,
		&account.Balance,
		&account.Currency,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return account, err
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	query := `
		SELECT id, user_id, account_number, agency, account_type, balance, currency, status, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		account := &domain.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.AccountNumber,
			&account.Agency,
			&account.AccountType,
			&account.Balance,
			&account.Currency,
			&account.Status,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *accountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	query := `
		SELECT id, user_id, account_number, agency, account_type, balance, currency, status, created_at, updated_at
		FROM accounts
		WHERE account_number = $1
	`

	account := &domain.Account{}
	err := r.db.QueryRow(ctx, query, accountNumber).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.Agency,
		&account.AccountType,
		&account.Balance,
		&account.Currency,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return account, err
}

func (r *accountRepository) UpdateBalance(ctx context.Context, id string, newBalance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, newBalance, id)
	return err
}

func (r *accountRepository) DebitAccount(ctx context.Context, id string, amount decimal.Decimal, txID string) (decimal.Decimal, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	defer tx.Rollback(ctx)

	var currentBalance decimal.Decimal
	query := `SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, id).Scan(&currentBalance)
	if err != nil {
		return decimal.Zero, err
	}

	newBalance := currentBalance.Sub(amount)
	if newBalance.IsNegative() {
		return decimal.Zero, ErrInsufficientFunds
	}

	updateQuery := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, newBalance, id)
	if err != nil {
		return decimal.Zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return decimal.Zero, err
	}

	return newBalance, nil
}

func (r *accountRepository) CreditAccount(ctx context.Context, id string, amount decimal.Decimal, txID string) (decimal.Decimal, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	defer tx.Rollback(ctx)

	var currentBalance decimal.Decimal
	query := `SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, id).Scan(&currentBalance)
	if err != nil {
		return decimal.Zero, err
	}

	newBalance := currentBalance.Add(amount)

	updateQuery := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, newBalance, id)
	if err != nil {
		return decimal.Zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return decimal.Zero, err
	}

	return newBalance, nil
}

func GenerateAccountNumber() string {
	return fmt.Sprintf("%010d", randomNumber())
}

func randomNumber() int64 {
	// In production, use crypto/rand for secure random numbers
	// This is simplified for demo purposes
	return 1000000000 + int64(randomInt(999999999))
}

func randomInt(max int) int {
	// Simplified random - in production use crypto/rand
	return max / 2
}
