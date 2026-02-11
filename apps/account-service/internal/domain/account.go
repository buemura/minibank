package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type AccountStatus string
type AccountType string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusBlocked  AccountStatus = "blocked"

	AccountTypeChecking AccountType = "checking"
	AccountTypeSavings  AccountType = "savings"
)

type Account struct {
	ID            string
	UserID        string
	AccountNumber string
	Agency        string
	AccountType   AccountType
	Balance       decimal.Decimal
	Currency      string
	Status        AccountStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
