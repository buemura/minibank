package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeTransfer    TransactionType = "transfer"
	TransactionTypeDeposit     TransactionType = "deposit"
	TransactionTypeWithdrawal  TransactionType = "withdrawal"
	TransactionTypeFee         TransactionType = "fee"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

type Transaction struct {
	ID                       string
	IdempotencyKey           string
	Type                     TransactionType
	Status                   TransactionStatus
	SourceAccountID          string
	SourceAccountNumber      string
	DestinationAccountID     string
	DestinationAccountNumber string
	Amount                   decimal.Decimal
	Currency                 string
	Description              string
	ReferenceID              string
	ProcessedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type IdempotencyKey struct {
	ID             string
	Key            string
	RequestHash    string
	ResponseStatus int
	ResponseBody   []byte
	TransactionID  string
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
