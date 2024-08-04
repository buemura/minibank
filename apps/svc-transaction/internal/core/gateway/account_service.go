package gateway

import "github.com/buemura/minibank/svc-transaction/internal/core/domain/account"

type AccountService interface {
	GetAccount(id string) (*account.Account, error)
	UpdateBalance(id string, newBalance int) (*account.Account, error)
}
