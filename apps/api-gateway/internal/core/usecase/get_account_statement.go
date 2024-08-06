package usecase

import (
	"github.com/buemura/minibank/api-gateway/internal/core/domain/statement"
	"github.com/buemura/minibank/api-gateway/internal/core/domain/transaction"
	"github.com/buemura/minibank/api-gateway/internal/core/gateway"
)

type GetAccountStatement struct {
	accService gateway.AccountService
	trxService gateway.TransactionService
}

func NewGetAccountStatement(accService gateway.AccountService, trxService gateway.TransactionService) *GetAccountStatement {
	return &GetAccountStatement{accService, trxService}
}

func (u *GetAccountStatement) Execute(in *statement.GetStatementIn) (*statement.Statement, error) {
	acc, err := u.accService.GetAccount(in.AccountID)
	if err != nil {
		return nil, err
	}

	trx, err := u.trxService.GetTransactions(&transaction.GetTransactionListIn{
		AccountID: in.AccountID,
		Page:      in.Page,
		Items:     in.Items,
	})
	if err != nil {
		return nil, err
	}

	return &statement.Statement{
		Account:      acc,
		Transactions: trx,
	}, nil
}
