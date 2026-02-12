package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/buemura/minibank/svc-transaction/internal/domain"
	"github.com/buemura/minibank/svc-transaction/internal/repository"
	accountpb "github.com/buemura/minibank/svc-transaction/proto/account/v1"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrIdempotencyMismatch = errors.New("idempotency key already used with different parameters")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrAccountNotFound     = errors.New("account not found")
	ErrTransferFailed      = errors.New("transfer failed")
	ErrInvalidAmount       = errors.New("invalid amount")

	WithdrawalFeeRate = decimal.NewFromFloat(0.05)
)

type TransactionService struct {
	transactionRepo repository.TransactionRepository
	idempotencyRepo repository.IdempotencyRepository
	accountClient   accountpb.AccountServiceClient
}

func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	idempotencyRepo repository.IdempotencyRepository,
	accountClient accountpb.AccountServiceClient,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		idempotencyRepo: idempotencyRepo,
		accountClient:   accountClient,
	}
}

type TransferInput struct {
	IdempotencyKey           string
	SourceAccountID          string
	DestinationAccountNumber string
	Amount                   string
	Description              string
}

type TransferResult struct {
	Success      bool
	Transaction  *domain.Transaction
	ErrorCode    string
	ErrorMessage string
}

type DepositInput struct {
	IdempotencyKey string
	AccountID      string
	Amount         string
	Description    string
}

type DepositResult struct {
	Success      bool
	Transaction  *domain.Transaction
	ErrorCode    string
	ErrorMessage string
}

type WithdrawInput struct {
	IdempotencyKey string
	AccountID      string
	Amount         string
	Description    string
}

type WithdrawResult struct {
	Success        bool
	Transaction    *domain.Transaction
	FeeTransaction *domain.Transaction
	ErrorCode      string
	ErrorMessage   string
}

func (s *TransactionService) Transfer(ctx context.Context, input TransferInput) (*TransferResult, error) {
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return &TransferResult{
			Success:      false,
			ErrorCode:    "INVALID_AMOUNT",
			ErrorMessage: "Invalid transfer amount",
		}, nil
	}

	requestHash := s.calculateRequestHash(input)

	existing, err := s.idempotencyRepo.GetByKey(ctx, input.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		if existing.RequestHash != requestHash {
			return &TransferResult{
				Success:      false,
				ErrorCode:    "IDEMPOTENCY_MISMATCH",
				ErrorMessage: "Idempotency key already used with different parameters",
			}, nil
		}

		if existing.TransactionID != "" {
			tx, err := s.transactionRepo.GetByID(ctx, existing.TransactionID)
			if err != nil {
				return nil, err
			}
			return &TransferResult{
				Success:     tx.Status == domain.TransactionStatusCompleted,
				Transaction: tx,
			}, nil
		}
	}

	idemKey := &domain.IdempotencyKey{
		Key:            input.IdempotencyKey,
		RequestHash:    requestHash,
		ResponseStatus: 0,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	if existing == nil {
		if err := s.idempotencyRepo.Create(ctx, idemKey); err != nil {
			return nil, err
		}
	}

	destAccountResp, err := s.accountClient.GetAccountByNumber(ctx, &accountpb.GetAccountByNumberRequest{
		AccountNumber: input.DestinationAccountNumber,
	})
	if err != nil || destAccountResp.Account == nil {
		return s.failTransfer(ctx, input.IdempotencyKey, "INVALID_DESTINATION", "Destination account not found")
	}

	sourceAccountResp, err := s.accountClient.GetAccount(ctx, &accountpb.GetAccountRequest{
		AccountId: input.SourceAccountID,
	})
	if err != nil || sourceAccountResp.Account == nil {
		return s.failTransfer(ctx, input.IdempotencyKey, "INVALID_SOURCE", "Source account not found")
	}

	txID := uuid.New().String()

	debitResp, err := s.accountClient.DebitAccount(ctx, &accountpb.DebitAccountRequest{
		AccountId:     input.SourceAccountID,
		Amount:        input.Amount,
		TransactionId: txID,
	})
	if err != nil || !debitResp.Success {
		errMsg := "Failed to debit source account"
		if debitResp != nil && debitResp.ErrorMessage != "" {
			errMsg = debitResp.ErrorMessage
		}
		return s.failTransfer(ctx, input.IdempotencyKey, "INSUFFICIENT_FUNDS", errMsg)
	}

	creditResp, err := s.accountClient.CreditAccount(ctx, &accountpb.CreditAccountRequest{
		AccountId:     destAccountResp.Account.Id,
		Amount:        input.Amount,
		TransactionId: txID,
	})
	if err != nil || !creditResp.Success {
		_, _ = s.accountClient.CreditAccount(ctx, &accountpb.CreditAccountRequest{
			AccountId:     input.SourceAccountID,
			Amount:        input.Amount,
			TransactionId: txID + "-rollback",
		})
		return s.failTransfer(ctx, input.IdempotencyKey, "TRANSFER_FAILED", "Failed to credit destination account")
	}

	now := time.Now()
	transaction := &domain.Transaction{
		IdempotencyKey:           input.IdempotencyKey,
		Type:                     domain.TransactionTypeTransfer,
		Status:                   domain.TransactionStatusCompleted,
		SourceAccountID:          input.SourceAccountID,
		SourceAccountNumber:      sourceAccountResp.Account.AccountNumber,
		DestinationAccountID:     destAccountResp.Account.Id,
		DestinationAccountNumber: input.DestinationAccountNumber,
		Amount:                   amount,
		Currency:                 "BRL",
		Description:              input.Description,
		ProcessedAt:              &now,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	responseBody, _ := json.Marshal(TransferResult{Success: true, Transaction: transaction})
	_ = s.idempotencyRepo.UpdateResponse(ctx, input.IdempotencyKey, 200, responseBody, transaction.ID)

	return &TransferResult{
		Success:     true,
		Transaction: transaction,
	}, nil
}

func (s *TransactionService) failTransfer(ctx context.Context, idempotencyKey, errorCode, errorMessage string) (*TransferResult, error) {
	result := &TransferResult{
		Success:      false,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	responseBody, _ := json.Marshal(result)
	_ = s.idempotencyRepo.UpdateResponse(ctx, idempotencyKey, 400, responseBody, "")
	return result, nil
}

func (s *TransactionService) Deposit(ctx context.Context, input DepositInput) (*DepositResult, error) {
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return &DepositResult{
			Success:      false,
			ErrorCode:    "INVALID_AMOUNT",
			ErrorMessage: "Invalid deposit amount",
		}, nil
	}

	requestHash := s.calculateDepositRequestHash(input)

	existing, err := s.idempotencyRepo.GetByKey(ctx, input.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		if existing.RequestHash != requestHash {
			return &DepositResult{
				Success:      false,
				ErrorCode:    "IDEMPOTENCY_MISMATCH",
				ErrorMessage: "Idempotency key already used with different parameters",
			}, nil
		}

		if existing.TransactionID != "" {
			tx, err := s.transactionRepo.GetByID(ctx, existing.TransactionID)
			if err != nil {
				return nil, err
			}
			return &DepositResult{
				Success:     tx.Status == domain.TransactionStatusCompleted,
				Transaction: tx,
			}, nil
		}
	}

	idemKey := &domain.IdempotencyKey{
		Key:            input.IdempotencyKey,
		RequestHash:    requestHash,
		ResponseStatus: 0,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	if existing == nil {
		if err := s.idempotencyRepo.Create(ctx, idemKey); err != nil {
			return nil, err
		}
	}

	txID := uuid.New().String()

	creditResp, err := s.accountClient.CreditAccount(ctx, &accountpb.CreditAccountRequest{
		AccountId:     input.AccountID,
		Amount:        input.Amount,
		TransactionId: txID,
	})
	if err != nil {
		return s.failDeposit(ctx, input.IdempotencyKey, "DEPOSIT_FAILED", "gRPC error: "+err.Error())
	}
	if !creditResp.Success {
		return s.failDeposit(ctx, input.IdempotencyKey, "DEPOSIT_FAILED", "Account service returned failure")
	}

	now := time.Now()
	transaction := &domain.Transaction{
		IdempotencyKey:       input.IdempotencyKey,
		Type:                 domain.TransactionTypeDeposit,
		Status:               domain.TransactionStatusCompleted,
		DestinationAccountID: input.AccountID,
		Amount:               amount,
		Currency:             "BRL",
		Description:          input.Description,
		ProcessedAt:          &now,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	responseBody, _ := json.Marshal(DepositResult{Success: true, Transaction: transaction})
	_ = s.idempotencyRepo.UpdateResponse(ctx, input.IdempotencyKey, 200, responseBody, transaction.ID)

	return &DepositResult{
		Success:     true,
		Transaction: transaction,
	}, nil
}

func (s *TransactionService) failDeposit(ctx context.Context, idempotencyKey, errorCode, errorMessage string) (*DepositResult, error) {
	result := &DepositResult{
		Success:      false,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	responseBody, _ := json.Marshal(result)
	_ = s.idempotencyRepo.UpdateResponse(ctx, idempotencyKey, 400, responseBody, "")
	return result, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, input WithdrawInput) (*WithdrawResult, error) {
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return &WithdrawResult{
			Success:      false,
			ErrorCode:    "INVALID_AMOUNT",
			ErrorMessage: "Invalid withdrawal amount",
		}, nil
	}

	requestHash := s.calculateWithdrawRequestHash(input)

	existing, err := s.idempotencyRepo.GetByKey(ctx, input.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		if existing.RequestHash != requestHash {
			return &WithdrawResult{
				Success:      false,
				ErrorCode:    "IDEMPOTENCY_MISMATCH",
				ErrorMessage: "Idempotency key already used with different parameters",
			}, nil
		}

		if existing.TransactionID != "" {
			tx, err := s.transactionRepo.GetByID(ctx, existing.TransactionID)
			if err != nil {
				return nil, err
			}
			return &WithdrawResult{
				Success:     tx.Status == domain.TransactionStatusCompleted,
				Transaction: tx,
			}, nil
		}
	}

	idemKey := &domain.IdempotencyKey{
		Key:            input.IdempotencyKey,
		RequestHash:    requestHash,
		ResponseStatus: 0,
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	if existing == nil {
		if err := s.idempotencyRepo.Create(ctx, idemKey); err != nil {
			return nil, err
		}
	}

	txID := uuid.New().String()

	debitResp, err := s.accountClient.DebitAccount(ctx, &accountpb.DebitAccountRequest{
		AccountId:     input.AccountID,
		Amount:        input.Amount,
		TransactionId: txID,
	})
	if err != nil {
		return s.failWithdraw(ctx, input.IdempotencyKey, "WITHDRAWAL_FAILED", "gRPC error: "+err.Error())
	}
	if !debitResp.Success {
		errMsg := "Account service returned failure"
		if debitResp.ErrorMessage != "" {
			errMsg = debitResp.ErrorMessage
		}
		return s.failWithdraw(ctx, input.IdempotencyKey, "WITHDRAWAL_FAILED", errMsg)
	}

	now := time.Now()
	transaction := &domain.Transaction{
		IdempotencyKey:  input.IdempotencyKey,
		Type:            domain.TransactionTypeWithdrawal,
		Status:          domain.TransactionStatusCompleted,
		SourceAccountID: input.AccountID,
		Amount:          amount,
		Currency:        "BRL",
		Description:     input.Description,
		ProcessedAt:     &now,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	feeAmount := amount.Mul(WithdrawalFeeRate)
	feeTxID := uuid.New().String()

	feeDebitResp, err := s.accountClient.DebitAccount(ctx, &accountpb.DebitAccountRequest{
		AccountId:     input.AccountID,
		Amount:        feeAmount.StringFixed(2),
		TransactionId: feeTxID,
	})
	if err != nil || !feeDebitResp.GetSuccess() {
		_, _ = s.accountClient.CreditAccount(ctx, &accountpb.CreditAccountRequest{
			AccountId:     input.AccountID,
			Amount:        input.Amount,
			TransactionId: txID + "-rollback",
		})
		return s.failWithdraw(ctx, input.IdempotencyKey, "WITHDRAWAL_FAILED", "Failed to debit withdrawal fee")
	}

	feeTransaction := &domain.Transaction{
		Type:            domain.TransactionTypeFee,
		Status:          domain.TransactionStatusCompleted,
		SourceAccountID: input.AccountID,
		Amount:          feeAmount,
		Currency:        "BRL",
		Description:     "Withdrawal fee (5%)",
		ReferenceID:     transaction.ID,
		ProcessedAt:     &now,
	}

	if err := s.transactionRepo.Create(ctx, feeTransaction); err != nil {
		return nil, err
	}

	responseBody, _ := json.Marshal(WithdrawResult{Success: true, Transaction: transaction, FeeTransaction: feeTransaction})
	_ = s.idempotencyRepo.UpdateResponse(ctx, input.IdempotencyKey, 200, responseBody, transaction.ID)

	return &WithdrawResult{
		Success:        true,
		Transaction:    transaction,
		FeeTransaction: feeTransaction,
	}, nil
}

func (s *TransactionService) failWithdraw(ctx context.Context, idempotencyKey, errorCode, errorMessage string) (*WithdrawResult, error) {
	result := &WithdrawResult{
		Success:      false,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	responseBody, _ := json.Marshal(result)
	_ = s.idempotencyRepo.UpdateResponse(ctx, idempotencyKey, 400, responseBody, "")
	return result, nil
}

func (s *TransactionService) calculateWithdrawRequestHash(input WithdrawInput) string {
	data := map[string]string{
		"account_id":  input.AccountID,
		"amount":      input.Amount,
		"description": input.Description,
	}
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

func (s *TransactionService) calculateDepositRequestHash(input DepositInput) string {
	data := map[string]string{
		"account_id":  input.AccountID,
		"amount":      input.Amount,
		"description": input.Description,
	}
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

func (s *TransactionService) GetTransaction(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	tx, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

func (s *TransactionService) GetTransactionHistory(ctx context.Context, accountID string, page, pageSize int) ([]*domain.Transaction, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.transactionRepo.GetByAccountID(ctx, accountID, page, pageSize)
}

func (s *TransactionService) GetStatement(ctx context.Context, accountID string, startDate, endDate time.Time, page, pageSize int) ([]*domain.Transaction, int, string, string, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	transactions, totalCount, err := s.transactionRepo.GetByAccountIDAndDateRange(ctx, accountID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, "", "", err
	}

	balanceResp, err := s.accountClient.GetBalance(ctx, &accountpb.GetBalanceRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, 0, "", "", err
	}

	closingBalance := balanceResp.Balance
	openingBalance := closingBalance

	return transactions, totalCount, openingBalance, closingBalance, nil
}

func (s *TransactionService) calculateRequestHash(input TransferInput) string {
	data := map[string]string{
		"source_account_id":          input.SourceAccountID,
		"destination_account_number": input.DestinationAccountNumber,
		"amount":                     input.Amount,
		"description":                input.Description,
	}
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}
