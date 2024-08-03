package usecase

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/buemura/minibank/svc-transaction/internal/core/gateway"
)

type GetTransactionList struct {
	cacheRepo gateway.CacheRepository
	txRepo    transaction.TransactionRepository
}

func NewGetTransactionList(cacheRepo gateway.CacheRepository, txRepo transaction.TransactionRepository) *GetTransactionList {
	return &GetTransactionList{
		cacheRepo: cacheRepo,
		txRepo:    txRepo,
	}
}

func (u *GetTransactionList) Execute(in *transaction.GetTransactionListIn) (*transaction.GetTransactionListOut, error) {
	if in.Page == 1 {
		slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Getting transaction list from cache for accountID: %s", in.AccountID))
		trxsCache, err := u.cacheRepo.Get(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, in.AccountID))
		if err != nil {
			return nil, err
		}
		if len(trxsCache) > 0 {
			slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Found transaction list on cache for account: %s", in.AccountID))
			return parseCachedTransactionList(trxsCache)
		}
	}

	slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Getting transaction list from DB for account: %s", in.AccountID))
	txs, err := u.txRepo.FindByAccountId(in)
	if err != nil {
		return nil, err
	}

	trxsToString, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}

	if in.Page > 1 {
		return txs, nil
	}

	slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Saving  transaction list in cache for account: %s", in.AccountID))
	err = u.cacheRepo.Set(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, in.AccountID), string(trxsToString), 60*time.Second)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func parseCachedTransactionList(trxsCache string) (*transaction.GetTransactionListOut, error) {
	var trxs *transaction.GetTransactionListOut
	err := json.Unmarshal([]byte(trxsCache), &trxs)
	if err != nil {
		return nil, err
	}
	return trxs, nil
}
