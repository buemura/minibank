package usecase

import (
	"github.com/buemura/minibank/svc-account/internal/core/entity"
	"github.com/buemura/minibank/svc-account/internal/core/repository"
)

type GetAccount struct {
	repo repository.AccountRepository
}

func NewGetAccount(repo repository.AccountRepository) *GetAccount {
	return &GetAccount{repo: repo}
}

func (u *GetAccount) Execute(id string) (*entity.Account, error) {
	acc, err := u.repo.GetAccountById(id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
