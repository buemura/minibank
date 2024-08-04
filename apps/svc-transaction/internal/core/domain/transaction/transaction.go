package transaction

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/lucsky/cuid"
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

func NewTransaction(in *CreateTransactionIn) (*Transaction, error) {
	if err := validate(in); err != nil {
		return nil, err
	}

	cuid, err := cuid.NewCrypto(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:                   cuid,
		AccountID:            in.AccountID,
		DestinationAccountID: in.DestinationAccountID,
		Amount:               in.Amount,
		Status:               Pending,
		TransactionType:      in.TransactionType,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

func validate(in *CreateTransactionIn) error {
	if in.AccountID == "" {
		return errors.New("invalid account id")
	}

	if in.TransactionType == "" {
		return errors.New("invalid transaction type")
	}

	if in.Amount <= 0 {
		return errors.New("invalid amount")
	}

	return nil
}
