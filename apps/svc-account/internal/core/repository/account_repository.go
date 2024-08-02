package repository

import "github.com/buemura/minibank/svc-account/internal/core/entity"

type AccountRepository interface {
	FindById(id string) (*entity.Account, error)
}
