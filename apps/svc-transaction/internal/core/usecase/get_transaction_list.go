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

func (u *GetTransactionList) Execute(accountID string) ([]*transaction.Transaction, error) {
	slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Getting transaction list from cache for accountID: %s", accountID))
	trxsCache, err := u.cacheRepo.Get(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, accountID))
	if err != nil {
		return nil, err
	}
	if len(trxsCache) > 0 {
		slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Found transaction list on cache for account: %s", accountID))
		return parseCachedTransactionList(trxsCache)
	}

	slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Getting transaction list from DB for account: %s", accountID))
	txs, err := u.txRepo.FindByAccountId(accountID)
	if err != nil {
		return nil, err
	}

	trxsToString, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[GetTransactionList][Execute] - Saving  transaction list in cache for account: %s", accountID))
	err = u.cacheRepo.Set(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, accountID), string(trxsToString), 60*time.Second)
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func parseCachedTransactionList(trxsCache string) ([]*transaction.Transaction, error) {
	var trxs []*transaction.Transaction
	err := json.Unmarshal([]byte(trxsCache), &trxs)
	if err != nil {
		return nil, err
	}
	return trxs, nil
}
