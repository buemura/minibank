package transaction

import (
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	Pending   TransactionStatus = "PENDING"
	Completed TransactionStatus = "COMPLETED"
	Failed    TransactionStatus = "FAILED"
)

const (
	Transfer   TransactionType = "TRANSFER"
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
)

type Transaction struct {
	ID                   string
	AccountID            string
	DestinationAccountID *string
	Amount               int
	Status               TransactionStatus
	TransactionType      TransactionType
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
