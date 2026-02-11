package grpc

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/buemura/minibank/svc-account/internal/domain"
	"github.com/buemura/minibank/svc-account/internal/service"
	pb "github.com/buemura/minibank/svc-account/proto/account/v1"
	"github.com/buemura/minibank/packages/logger"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	accountService *service.AccountService
}

func NewAccountServer(accountService *service.AccountService) *AccountServer {
	return &AccountServer{accountService: accountService}
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	logger.Info("CreateAccount request received", zap.String("user_id", req.UserId), zap.String("account_type", req.AccountType))

	accountType := domain.AccountTypeChecking
	if req.AccountType == "savings" {
		accountType = domain.AccountTypeSavings
	}

	account, err := s.accountService.CreateAccount(ctx, service.CreateAccountInput{
		UserID:      req.UserId,
		AccountType: accountType,
	})
	if err != nil {
		logger.Error("CreateAccount failed", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create account")
	}

	logger.Info("Account created successfully", zap.String("account_id", account.ID), zap.String("user_id", req.UserId))
	return &pb.CreateAccountResponse{
		Account: domainAccountToProto(account),
	}, nil
}

func (s *AccountServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	logger.Debug("GetAccount request received", zap.String("account_id", req.AccountId))

	account, err := s.accountService.GetAccount(ctx, req.AccountId)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			logger.Warn("GetAccount: account not found", zap.String("account_id", req.AccountId))
			return nil, status.Error(codes.NotFound, "account not found")
		}
		logger.Error("GetAccount failed", zap.String("account_id", req.AccountId), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get account")
	}

	return &pb.GetAccountResponse{
		Account: domainAccountToProto(account),
	}, nil
}

func (s *AccountServer) GetAccountByUserId(ctx context.Context, req *pb.GetAccountByUserIdRequest) (*pb.GetAccountByUserIdResponse, error) {
	accounts, err := s.accountService.GetAccountsByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get accounts")
	}

	var protoAccounts []*pb.Account
	for _, acc := range accounts {
		protoAccounts = append(protoAccounts, domainAccountToProto(acc))
	}

	return &pb.GetAccountByUserIdResponse{
		Accounts: protoAccounts,
	}, nil
}

func (s *AccountServer) GetAccountByNumber(ctx context.Context, req *pb.GetAccountByNumberRequest) (*pb.GetAccountByNumberResponse, error) {
	account, err := s.accountService.GetAccountByNumber(ctx, req.AccountNumber)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Error(codes.Internal, "failed to get account")
	}

	return &pb.GetAccountByNumberResponse{
		Account: domainAccountToProto(account),
	}, nil
}

func (s *AccountServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	account, err := s.accountService.GetBalance(ctx, req.AccountId)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}
		return nil, status.Error(codes.Internal, "failed to get balance")
	}

	return &pb.GetBalanceResponse{
		AccountId: account.ID,
		Balance:   account.Balance.String(),
		Currency:  account.Currency,
		UpdatedAt: timestamppb.New(account.UpdatedAt),
	}, nil
}

func (s *AccountServer) DebitAccount(ctx context.Context, req *pb.DebitAccountRequest) (*pb.DebitAccountResponse, error) {
	logger.Info("DebitAccount request received", zap.String("account_id", req.AccountId), zap.String("amount", req.Amount), zap.String("transaction_id", req.TransactionId))

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		logger.Warn("DebitAccount: invalid amount format", zap.String("amount", req.Amount), zap.Error(err))
		return &pb.DebitAccountResponse{
			Success:      false,
			ErrorMessage: "invalid amount format",
		}, nil
	}

	newBalance, err := s.accountService.DebitAccount(ctx, req.AccountId, amount, req.TransactionId)
	if err != nil {
		errMsg := "failed to debit account"
		switch {
		case errors.Is(err, service.ErrInsufficientFunds):
			errMsg = "insufficient funds"
			logger.Warn("DebitAccount: insufficient funds", zap.String("account_id", req.AccountId), zap.String("amount", req.Amount))
		case errors.Is(err, service.ErrAccountNotFound):
			errMsg = "account not found"
			logger.Warn("DebitAccount: account not found", zap.String("account_id", req.AccountId))
		case errors.Is(err, service.ErrAccountBlocked):
			errMsg = "account is blocked"
			logger.Warn("DebitAccount: account blocked", zap.String("account_id", req.AccountId))
		default:
			logger.Error("DebitAccount failed", zap.String("account_id", req.AccountId), zap.Error(err))
		}
		return &pb.DebitAccountResponse{
			Success:      false,
			ErrorMessage: errMsg,
		}, nil
	}

	logger.Info("DebitAccount successful", zap.String("account_id", req.AccountId), zap.String("new_balance", newBalance.String()))
	return &pb.DebitAccountResponse{
		Success:    true,
		NewBalance: newBalance.String(),
	}, nil
}

func (s *AccountServer) CreditAccount(ctx context.Context, req *pb.CreditAccountRequest) (*pb.CreditAccountResponse, error) {
	logger.Info("CreditAccount request received", zap.String("account_id", req.AccountId), zap.String("amount", req.Amount), zap.String("transaction_id", req.TransactionId))

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		logger.Warn("CreditAccount: invalid amount format", zap.String("amount", req.Amount), zap.Error(err))
		return &pb.CreditAccountResponse{
			Success: false,
		}, nil
	}

	newBalance, err := s.accountService.CreditAccount(ctx, req.AccountId, amount, req.TransactionId)
	if err != nil {
		logger.Error("CreditAccount failed", zap.String("account_id", req.AccountId), zap.Error(err))
		return &pb.CreditAccountResponse{
			Success: false,
		}, nil
	}

	logger.Info("CreditAccount successful", zap.String("account_id", req.AccountId), zap.String("new_balance", newBalance.String()))
	return &pb.CreditAccountResponse{
		Success:    true,
		NewBalance: newBalance.String(),
	}, nil
}

func domainAccountToProto(account *domain.Account) *pb.Account {
	return &pb.Account{
		Id:            account.ID,
		UserId:        account.UserID,
		AccountNumber: account.AccountNumber,
		Agency:        account.Agency,
		AccountType:   string(account.AccountType),
		Balance:       account.Balance.String(),
		Currency:      account.Currency,
		Status:        string(account.Status),
		CreatedAt:     timestamppb.New(account.CreatedAt),
		UpdatedAt:     timestamppb.New(account.UpdatedAt),
	}
}
