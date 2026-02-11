package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/buemura/minibank/svc-account/internal/domain"
	"github.com/buemura/minibank/svc-account/internal/repository"
	"github.com/buemura/minibank/svc-account/internal/repository/postgres"
	"github.com/shopspring/decimal"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrAccountBlocked    = errors.New("account is blocked")
)

type AccountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

type CreateAccountInput struct {
	UserID      string
	AccountType domain.AccountType
}

func (s *AccountService) CreateAccount(ctx context.Context, input CreateAccountInput) (*domain.Account, error) {
	accountNumber := generateAccountNumber()

	account := &domain.Account{
		UserID:        input.UserID,
		AccountNumber: accountNumber,
		Agency:        "0001",
		AccountType:   input.AccountType,
		Balance:       decimal.NewFromInt(0),
		Currency:      "BRL",
		Status:        domain.AccountStatusActive,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (s *AccountService) GetAccountsByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	return s.accountRepo.GetByUserID(ctx, userID)
}

func (s *AccountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (s *AccountService) GetBalance(ctx context.Context, accountID string) (*domain.Account, error) {
	return s.GetAccount(ctx, accountID)
}

func (s *AccountService) DebitAccount(ctx context.Context, accountID string, amount decimal.Decimal, txID string) (decimal.Decimal, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, ErrInvalidAmount
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return decimal.Zero, err
	}
	if account == nil {
		return decimal.Zero, ErrAccountNotFound
	}
	if account.Status == domain.AccountStatusBlocked {
		return decimal.Zero, ErrAccountBlocked
	}

	newBalance, err := s.accountRepo.DebitAccount(ctx, accountID, amount, txID)
	if err != nil {
		if errors.Is(err, postgres.ErrInsufficientFunds) {
			return decimal.Zero, ErrInsufficientFunds
		}
		return decimal.Zero, err
	}

	return newBalance, nil
}

func (s *AccountService) CreditAccount(ctx context.Context, accountID string, amount decimal.Decimal, txID string) (decimal.Decimal, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, ErrInvalidAmount
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return decimal.Zero, err
	}
	if account == nil {
		return decimal.Zero, ErrAccountNotFound
	}
	if account.Status == domain.AccountStatusBlocked {
		return decimal.Zero, ErrAccountBlocked
	}

	newBalance, err := s.accountRepo.CreditAccount(ctx, accountID, amount, txID)
	if err != nil {
		return decimal.Zero, err
	}

	return newBalance, nil
}

func generateAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(9000000000))
	return fmt.Sprintf("%010d", n.Int64()+1000000000)
}
