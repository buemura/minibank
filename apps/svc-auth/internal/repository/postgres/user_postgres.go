package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/buemura/minibank/svc-auth/internal/domain"
	"github.com/buemura/minibank/svc-auth/internal/repository"
	"github.com/buemura/minibank/packages/logger"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	logger.Debug("userRepository.Create: inserting user", zap.String("email", user.Email))

	query := `
		INSERT INTO users (email, password_hash, full_name, phone, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		logger.Error("userRepository.Create: failed to insert user", zap.String("email", user.Email), zap.Error(err))
		return err
	}

	logger.Debug("userRepository.Create: user inserted successfully", zap.String("user_id", user.ID))
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	logger.Debug("userRepository.GetByEmail: querying user", zap.String("email", email))

	query := `
		SELECT id, email, password_hash, full_name, phone, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		logger.Debug("userRepository.GetByEmail: user not found", zap.String("email", email))
		return nil, nil
	}

	if err != nil {
		logger.Error("userRepository.GetByEmail: query failed", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	logger.Debug("userRepository.GetByEmail: user found", zap.String("user_id", user.ID))
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, phone = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query,
		user.Email,
		user.FullName,
		user.Phone,
		user.Status,
		user.ID,
	)

	return err
}
