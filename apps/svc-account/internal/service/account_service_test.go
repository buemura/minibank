package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/buemura/minibank/svc-account/internal/domain"
	"github.com/buemura/minibank/svc-account/internal/repository/mocks"
	"github.com/buemura/minibank/svc-account/internal/repository/postgres"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newActiveAccount(id, userID string) *domain.Account {
	return &domain.Account{
		ID:            id,
		UserID:        userID,
		AccountNumber: "1234567890",
		Agency:        "0001",
		AccountType:   domain.AccountTypeChecking,
		Balance:       decimal.NewFromInt(1000),
		Currency:      "BRL",
		Status:        domain.AccountStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func TestCreateAccount_Success(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	account, err := svc.CreateAccount(context.Background(), CreateAccountInput{
		UserID:      "user-1",
		AccountType: domain.AccountTypeChecking,
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", account.UserID)
	assert.Equal(t, "0001", account.Agency)
	assert.Equal(t, "BRL", account.Currency)
	assert.Equal(t, domain.AccountStatusActive, account.Status)
	assert.True(t, account.Balance.Equal(decimal.Zero))
	assert.Len(t, account.AccountNumber, 10)
	repo.AssertExpectations(t)
}

func TestCreateAccount_RepoError(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	_, err := svc.CreateAccount(context.Background(), CreateAccountInput{
		UserID:      "user-1",
		AccountType: domain.AccountTypeChecking,
	})

	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}

func TestGetAccount_Success(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	expected := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(expected, nil)

	account, err := svc.GetAccount(context.Background(), "acc-1")

	require.NoError(t, err)
	assert.Equal(t, "acc-1", account.ID)
	repo.AssertExpectations(t)
}

func TestGetAccount_NotFound(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, nil)

	_, err := svc.GetAccount(context.Background(), "acc-1")

	assert.ErrorIs(t, err, ErrAccountNotFound)
	repo.AssertExpectations(t)
}

func TestGetAccount_RepoError(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, errors.New("db error"))

	_, err := svc.GetAccount(context.Background(), "acc-1")

	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}

func TestGetAccountsByUserID_Success(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	accounts := []*domain.Account{
		newActiveAccount("acc-1", "user-1"),
		newActiveAccount("acc-2", "user-1"),
	}
	repo.On("GetByUserID", mock.Anything, "user-1").Return(accounts, nil)

	result, err := svc.GetAccountsByUserID(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestGetAccountByNumber_Success(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	expected := newActiveAccount("acc-1", "user-1")
	repo.On("GetByAccountNumber", mock.Anything, "1234567890").Return(expected, nil)

	account, err := svc.GetAccountByNumber(context.Background(), "1234567890")

	require.NoError(t, err)
	assert.Equal(t, "acc-1", account.ID)
	repo.AssertExpectations(t)
}

func TestGetAccountByNumber_NotFound(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("GetByAccountNumber", mock.Anything, "0000000000").Return(nil, nil)

	_, err := svc.GetAccountByNumber(context.Background(), "0000000000")

	assert.ErrorIs(t, err, ErrAccountNotFound)
	repo.AssertExpectations(t)
}

func TestGetBalance_DelegatesToGetAccount(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	expected := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(expected, nil)

	account, err := svc.GetBalance(context.Background(), "acc-1")

	require.NoError(t, err)
	assert.Equal(t, expected.Balance, account.Balance)
	repo.AssertExpectations(t)
}

func TestDebitAccount_Success(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("DebitAccount", mock.Anything, "acc-1", decimal.NewFromInt(100), "tx-1").
		Return(decimal.NewFromInt(900), nil)

	newBalance, err := svc.DebitAccount(context.Background(), "acc-1", decimal.NewFromInt(100), "tx-1")

	require.NoError(t, err)
	assert.True(t, newBalance.Equal(decimal.NewFromInt(900)))
	repo.AssertExpectations(t)
}

func TestDebitAccount_InvalidAmount_Zero(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	_, err := svc.DebitAccount(context.Background(), "acc-1", decimal.Zero, "tx-1")

	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestDebitAccount_InvalidAmount_Negative(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	_, err := svc.DebitAccount(context.Background(), "acc-1", decimal.NewFromInt(-50), "tx-1")

	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestDebitAccount_AccountNotFound(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, nil)

	_, err := svc.DebitAccount(context.Background(), "acc-1", decimal.NewFromInt(100), "tx-1")

	assert.ErrorIs(t, err, ErrAccountNotFound)
	repo.AssertExpectations(t)
}

func TestDebitAccount_AccountBlocked(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	account.Status = domain.AccountStatusBlocked
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)

	_, err := svc.DebitAccount(context.Background(), "acc-1", decimal.NewFromInt(100), "tx-1")

	assert.ErrorIs(t, err, ErrAccountBlocked)
	repo.AssertExpectations(t)
}

func TestDebitAccount_InsufficientFunds(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("DebitAccount", mock.Anything, "acc-1", decimal.NewFromInt(5000), "tx-1").
		Return(decimal.Zero, postgres.ErrInsufficientFunds)

	_, err := svc.DebitAccount(context.Background(), "acc-1", decimal.NewFromInt(5000), "tx-1")

	assert.ErrorIs(t, err, ErrInsufficientFunds)
	repo.AssertExpectations(t)
}

func TestDebitAccount_RepoError(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("DebitAccount", mock.Anything, "acc-1", decimal.NewFromInt(100), "tx-1").
		Return(decimal.Zero, errors.New("db error"))

	_, err := svc.DebitAccount(context.Background(), "acc-1", decimal.NewFromInt(100), "tx-1")

	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}

func TestCreditAccount_Success(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("CreditAccount", mock.Anything, "acc-1", decimal.NewFromInt(200), "tx-1").
		Return(decimal.NewFromInt(1200), nil)

	newBalance, err := svc.CreditAccount(context.Background(), "acc-1", decimal.NewFromInt(200), "tx-1")

	require.NoError(t, err)
	assert.True(t, newBalance.Equal(decimal.NewFromInt(1200)))
	repo.AssertExpectations(t)
}

func TestCreditAccount_InvalidAmount(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	_, err := svc.CreditAccount(context.Background(), "acc-1", decimal.NewFromInt(-10), "tx-1")

	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestCreditAccount_AccountNotFound(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	repo.On("GetByID", mock.Anything, "acc-1").Return(nil, nil)

	_, err := svc.CreditAccount(context.Background(), "acc-1", decimal.NewFromInt(200), "tx-1")

	assert.ErrorIs(t, err, ErrAccountNotFound)
	repo.AssertExpectations(t)
}

func TestCreditAccount_AccountBlocked(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	account.Status = domain.AccountStatusBlocked
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)

	_, err := svc.CreditAccount(context.Background(), "acc-1", decimal.NewFromInt(200), "tx-1")

	assert.ErrorIs(t, err, ErrAccountBlocked)
	repo.AssertExpectations(t)
}

func TestCreditAccount_RepoError(t *testing.T) {
	repo := new(mocks.MockAccountRepository)
	svc := NewAccountService(repo)

	account := newActiveAccount("acc-1", "user-1")
	repo.On("GetByID", mock.Anything, "acc-1").Return(account, nil)
	repo.On("CreditAccount", mock.Anything, "acc-1", decimal.NewFromInt(200), "tx-1").
		Return(decimal.Zero, errors.New("db error"))

	_, err := svc.CreditAccount(context.Background(), "acc-1", decimal.NewFromInt(200), "tx-1")

	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}
