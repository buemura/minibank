package transaction

import "github.com/buemura/minibank-api-gateway/internal/core/domain/common"

type GetTransactionListIn struct {
	AccountID string `json:"account_id"`
	Page      int    `json:"page"`
	Items     int    `json:"items"`
}

type PaginatedTransactionList struct {
	Data []*Transaction           `json:"data"`
	Meta common.PaginationMetaOut `json:"meta"`
}

type GetTransactionListOut struct {
	PaginatedTransactionList
}

type CreateTransactionIn struct {
	AccountID            string          `json:"account_id"`
	DestinationAccountID *string         `json:"destination_account_id,omitempty"`
	Amount               int             `json:"amount"`
	TransactionType      TransactionType `json:"transaction_type"`
}
