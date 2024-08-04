package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/buemura/minibank/svc-account/config"
	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
	"github.com/buemura/minibank/svc-account/internal/core/gateway"
)

type UpdateBalance struct {
	cacheRepo gateway.CacheRepository
	accRepo   account.AccountRepository
}

func NewUpdateBalance(cacheRepo gateway.CacheRepository, accRepo account.AccountRepository) *UpdateBalance {
	return &UpdateBalance{cacheRepo: cacheRepo, accRepo: accRepo}
}

func (u *UpdateBalance) Execute(in *account.UpdateBalanceIn) (*account.Account, error) {
	slog.Info(fmt.Sprintf("[UpdateBalance][Execute] - Checking if account exists: %s", in.ID))
	acc, err := u.accRepo.FindById(in.ID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, account.ErrAccountNotFound
	}

	acc.Balance = in.NewBalance
	acc.UpdatedAt = time.Now()

	slog.Info(fmt.Sprintf("[UpdateBalance][Execute] - Updating account in DB: %s", acc.ID))
	_, err = u.accRepo.Update(acc)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[UpdateBalance][Execute] - Clearing account from cache: %s", acc.ID))
	err = u.cacheRepo.Delete(fmt.Sprintf("%s:%s", config.CACHE_ACCOUNT_KEY_PREFIX, acc.ID))
	if err != nil {
		return nil, err
	}

	return acc, nil
}
