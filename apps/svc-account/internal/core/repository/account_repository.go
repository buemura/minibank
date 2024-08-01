package repository

import "github.com/buemura/minibank/svc-account/internal/core/entity"

type AccountRepository interface {
	GetAccountById(id string) (*entity.Account, error)
}
