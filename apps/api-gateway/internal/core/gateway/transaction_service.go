package gateway

import "github.com/buemura/minibank/api-gateway/internal/core/domain/transaction"

type TransactionService interface {
	GetTransactions(in *transaction.GetTransactionListIn) (*transaction.PaginatedTransactionList, error)
	CreateTransaction(trx *transaction.CreateTransactionIn) (*transaction.Transaction, error)
}
