package statement

import (
	"github.com/buemura/minibank-api-gateway/internal/core/domain/account"
	"github.com/buemura/minibank-api-gateway/internal/core/domain/transaction"
)

type GetStatementIn struct {
	AccountID string `json:"account_id"`
	Page      int    `json:"page"`
	Items     int    `json:"items"`
}

type Statement struct {
	Account      *account.Account                      `json:"account"`
	Transactions *transaction.PaginatedTransactionList `json:"transactions"`
}
