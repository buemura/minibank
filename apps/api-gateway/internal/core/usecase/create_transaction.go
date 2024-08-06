package usecase

import (
	"github.com/buemura/minibank/api-gateway/internal/core/domain/transaction"
	"github.com/buemura/minibank/api-gateway/internal/core/gateway"
)

type CreateTransaction struct {
	trxService gateway.TransactionService
}

func NewCreateTransaction(trxService gateway.TransactionService) *CreateTransaction {
	return &CreateTransaction{trxService}
}

func (u *CreateTransaction) Execute(in *transaction.CreateTransactionIn) (*transaction.Transaction, error) {
	trx, err := u.trxService.CreateTransaction(in)
	if err != nil {
		return nil, err
	}
	return trx, nil
}
