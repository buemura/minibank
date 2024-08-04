package usecase

import (
	"fmt"
	"log/slog"

	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/buemura/minibank/svc-transaction/internal/core/gateway"
)

type CreateTransaction struct {
	cacheRepo gateway.CacheRepository
	trxRepo   transaction.TransactionRepository
}

func NewCreateTransaction(cacheRepo gateway.CacheRepository, trxRepo transaction.TransactionRepository) *CreateTransaction {
	return &CreateTransaction{cacheRepo: cacheRepo, trxRepo: trxRepo}
}

func (u *CreateTransaction) Execute(in *transaction.CreateTransactionIn) (*transaction.Transaction, error) {
	trx, err := transaction.NewTransaction(in)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[CreateTransaction][Execute] - Saving transaction in DB: %s", trx.ID))
	_, err = u.trxRepo.Create(trx)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[CreateTransaction][Execute] - Clearing transaction list cache for accountID: %s", trx.ID))
	err = u.cacheRepo.Delete(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, trx.AccountID))
	if err != nil {
		return nil, err
	}

	return trx, nil
}
