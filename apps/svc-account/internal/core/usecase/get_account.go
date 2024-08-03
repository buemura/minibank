package usecase

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
	"github.com/buemura/minibank/svc-account/internal/core/gateway"
)

type GetAccount struct {
	cacheRepo gateway.CacheRepository
	accRepo   account.AccountRepository
}

func NewGetAccount(cacheRepo gateway.CacheRepository, accRepo account.AccountRepository) *GetAccount {
	return &GetAccount{cacheRepo: cacheRepo, accRepo: accRepo}
}

func (u *GetAccount) Execute(id string) (*account.Account, error) {
	slog.Info(fmt.Sprintf("[GetAccount][Execute] - Getting account from cache: %s", id))
	accCache, err := u.cacheRepo.Get(fmt.Sprintf("account:%s", id))
	if err != nil {
		return nil, err
	}
	if len(accCache) > 0 {
		slog.Info(fmt.Sprintf("[GetAccount][Execute] - Found account on cache: %s", id))
		return parseCachedAccount(accCache)
	}

	slog.Info(fmt.Sprintf("[GetAccount][Execute] - Getting account from DB: %s", id))
	acc, err := u.accRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	accToString, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[GetAccount][Execute] - Saving account in cache: %s", id))
	err = u.cacheRepo.Set(fmt.Sprintf("account:%s", id), string(accToString), 60*time.Second)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func parseCachedAccount(accCache string) (*account.Account, error) {
	var acc *account.Account
	err := json.Unmarshal([]byte(accCache), &acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
