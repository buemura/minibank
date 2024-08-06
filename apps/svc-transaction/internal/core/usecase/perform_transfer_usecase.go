package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/account"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/buemura/minibank/svc-transaction/internal/core/gateway"
)

type PerformTransfer struct {
	accService gateway.AccountService
	cacheRepo  gateway.CacheRepository
	txRepo     transaction.TransactionRepository
}

func NewPerformTransfer(
	accService gateway.AccountService,
	cacheRepo gateway.CacheRepository,
	txRepo transaction.TransactionRepository,
) *PerformTransfer {
	return &PerformTransfer{accService: accService, cacheRepo: cacheRepo, txRepo: txRepo}
}

func (u *PerformTransfer) Execute(trx *transaction.Transaction) error {
	slog.Info(fmt.Sprintf("[PerformTransfer][Execute] - Validating destination account: %s", *trx.DestinationAccountID))
	destinationAccount, err := u.accService.GetAccount(*trx.DestinationAccountID)
	if err != nil {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return err
	}

	slog.Info(fmt.Sprintf("[PerformTransfer][Execute] - Validating origin account: %s", trx.AccountID))
	originAccount, err := u.accService.GetAccount(trx.AccountID)
	if err != nil {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return err
	}

	if originAccount.Balance < trx.Amount {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return account.ErrInsufficientBalance
	}

	slog.Info("[PerformTransfer][Execute] - Updating accounts balances")
	_, err = u.accService.UpdateBalance(originAccount.ID, originAccount.Balance-trx.Amount)
	if err != nil {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return err
	}
	_, err = u.accService.UpdateBalance(destinationAccount.ID, destinationAccount.Balance+trx.Amount)
	if err != nil {
		err = u.updateTransactionFailed(trx)
		if err != nil {
			return err
		}
		return err
	}

	slog.Info(fmt.Sprintf("[PerformTransfer][Execute] - Updating transaction status: %s", trx.ID))
	trx.Status = transaction.Completed
	trx.UpdatedAt = time.Now()
	_, err = u.txRepo.Update(trx)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("[PerformTransfer][Execute] - Clearing transaction list cache for accountID: %s", trx.AccountID))
	err = u.cacheRepo.Delete(fmt.Sprintf("%s:%s", config.CACHE_TRANSACTION_LIST_KEY_PREFIX, trx.AccountID))
	if err != nil {
		return err
	}

	return nil
}

func (u *PerformTransfer) updateTransactionFailed(trx *transaction.Transaction) error {
	trx.Status = transaction.Failed
	trx.UpdatedAt = time.Now()
	_, err := u.txRepo.Update(trx)
	if err != nil {
		return err
	}
	return nil
}
