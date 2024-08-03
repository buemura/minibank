package usecase

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/buemura/minibank/svc-account/config"
	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
	"github.com/buemura/minibank/svc-account/internal/core/gateway"
)

type CreateAccount struct {
	cacheRepo gateway.CacheRepository
	accRepo   account.AccountRepository
}

func NewCreateAccount(cacheRepo gateway.CacheRepository, accRepo account.AccountRepository) *CreateAccount {
	return &CreateAccount{cacheRepo: cacheRepo, accRepo: accRepo}
}

func (u *CreateAccount) Execute(in *account.CreateAccountIn) (*account.Account, error) {
	slog.Info(fmt.Sprintf("[CreateAccount][Execute] - Checking if account account already exists: %s", in.OwnerDocument))
	existingAcc, err := u.accRepo.FindByOwnerDocument(in.OwnerDocument)
	if err != nil {
		return nil, err
	}
	if existingAcc != nil {
		return nil, account.ErrAccountAlreadyExists
	}

	acc, err := account.NewAccount(in)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[CreateAccount][Execute] - Saving account in db: %s", acc.ID))
	_, err = u.accRepo.Create(acc)
	if err != nil {
		return nil, err
	}

	accToString, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[CreateAccount][Execute] - Saving account in cache: %s", acc.ID))
	err = u.cacheRepo.Set(fmt.Sprintf("%s:%s", config.CACHE_ACCOUNT_KEY_PREFIX, acc.ID), string(accToString), 60*time.Second)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
