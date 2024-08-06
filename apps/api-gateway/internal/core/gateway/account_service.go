package gateway

import "github.com/buemura/minibank/api-gateway/internal/core/domain/account"

type AccountService interface {
	GetAccount(id string) (*account.Account, error)
}
