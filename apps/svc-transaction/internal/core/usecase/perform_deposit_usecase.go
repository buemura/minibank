package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/buemura/minibank/svc-transaction/internal/core/gateway"
)

type PerformDeposit struct {
	accService gateway.AccountService
	cacheRepo  gateway.CacheRepository
	txRepo     transaction.TransactionRepository
}

func NewPerformDeposit(
	accService gateway.AccountService,
	cacheRepo gateway.CacheRepository,
	txRepo transaction.TransactionRepository,
) *PerformDeposit {
	return &PerformDeposit{accService: accService, cacheRepo: cacheRepo, txRepo: txRepo}
}

func (u *PerformDeposit) Execute(trx *transaction.Transaction) error {
	slog.Info(fmt.Sprintf("[PerformDeposit][Execute] - Validating account: %s", trx.AccountID))
	acc, err := u.accService.GetAccount(trx.AccountID)
	if err != nil {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return err
	}

	slog.Info("[PerformDeposit][Execute] - Updating accounts balances")
	_, err = u.accService.UpdateBalance(acc.ID, acc.Balance+trx.Amount)
	if err != nil {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return err
	}

	slog.Info(fmt.Sprintf("[PerformDeposit][Execute] - Updating transaction status: %s", trx.ID))
	trx.Status = transaction.Completed
	trx.UpdatedAt = time.Now()
	_, err = u.txRepo.Update(trx)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("[PerformDeposit][Execute] - Clearing transaction list cache for accountID: %s", trx.AccountID))
	err = u.cacheRepo.Delete(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, trx.AccountID))
	if err != nil {
		return err
	}

	return nil
}

func (u *PerformDeposit) updateTransactionFailed(trx *transaction.Transaction) error {
	trx.Status = transaction.Failed
	trx.UpdatedAt = time.Now()
	_, err := u.txRepo.Update(trx)
	if err != nil {
		return err
	}
	return nil
}
