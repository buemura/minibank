package usecase

import (
	"github.com/buemura/minibank/svc-account/internal/domain/account"
)

type GetAccount struct {
	repo account.AccountRepository
}

func NewGetAccount(repo account.AccountRepository) *GetAccount {
	return &GetAccount{repo: repo}
}

func (u *GetAccount) Execute(id string) (*account.Account, error) {
	acc, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
