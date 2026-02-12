package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/buemura/minibank/packages/logger"
	"github.com/buemura/minibank/svc-transaction/internal/domain"
	"github.com/buemura/minibank/svc-transaction/internal/service"
	pb "github.com/buemura/minibank/svc-transaction/proto/transaction/v1"
)

type TransactionServer struct {
	pb.UnimplementedTransactionServiceServer
	transactionService *service.TransactionService
}

func NewTransactionServer(transactionService *service.TransactionService) *TransactionServer {
	return &TransactionServer{transactionService: transactionService}
}

func (s *TransactionServer) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	logger.Info("Transfer request received",
		zap.String("idempotency_key", req.IdempotencyKey),
		zap.String("source_account_id", req.SourceAccountId),
		zap.String("destination_account_number", req.DestinationAccountNumber),
		zap.String("amount", req.Amount))

	result, err := s.transactionService.Transfer(ctx, service.TransferInput{
		IdempotencyKey:           req.IdempotencyKey,
		SourceAccountID:          req.SourceAccountId,
		DestinationAccountNumber: req.DestinationAccountNumber,
		Amount:                   req.Amount,
		Description:              req.Description,
	})
	if err != nil {
		logger.Error("Transfer failed", zap.String("idempotency_key", req.IdempotencyKey), zap.Error(err))
		return nil, status.Error(codes.Internal, "transfer failed")
	}

	resp := &pb.TransferResponse{
		Success:      result.Success,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
	}

	if result.Transaction != nil {
		resp.Transaction = domainTransactionToProto(result.Transaction)
		logger.Info("Transfer completed", zap.String("transaction_id", result.Transaction.ID), zap.Bool("success", result.Success))
	} else if !result.Success {
		logger.Warn("Transfer not successful", zap.String("error_code", result.ErrorCode), zap.String("error_message", result.ErrorMessage))
	}

	return resp, nil
}

func (s *TransactionServer) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	logger.Info("Deposit request received",
		zap.String("idempotency_key", req.IdempotencyKey),
		zap.String("account_id", req.AccountId),
		zap.String("amount", req.Amount))

	result, err := s.transactionService.Deposit(ctx, service.DepositInput{
		IdempotencyKey: req.IdempotencyKey,
		AccountID:      req.AccountId,
		Amount:         req.Amount,
		Description:    req.Description,
	})
	if err != nil {
		logger.Error("Deposit failed", zap.String("idempotency_key", req.IdempotencyKey), zap.Error(err))
		return nil, status.Error(codes.Internal, "deposit failed")
	}

	resp := &pb.DepositResponse{
		Success:      result.Success,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
	}

	if result.Transaction != nil {
		resp.Transaction = domainTransactionToProto(result.Transaction)
		logger.Info("Deposit completed", zap.String("transaction_id", result.Transaction.ID), zap.Bool("success", result.Success))
	} else if !result.Success {
		logger.Warn("Deposit not successful", zap.String("error_code", result.ErrorCode), zap.String("error_message", result.ErrorMessage))
	}

	return resp, nil
}

func (s *TransactionServer) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	logger.Info("Withdraw request received",
		zap.String("idempotency_key", req.IdempotencyKey),
		zap.String("account_id", req.AccountId),
		zap.String("amount", req.Amount))

	result, err := s.transactionService.Withdraw(ctx, service.WithdrawInput{
		IdempotencyKey: req.IdempotencyKey,
		AccountID:      req.AccountId,
		Amount:         req.Amount,
		Description:    req.Description,
	})
	if err != nil {
		logger.Error("Withdraw failed", zap.String("idempotency_key", req.IdempotencyKey), zap.Error(err))
		return nil, status.Error(codes.Internal, "withdrawal failed")
	}

	resp := &pb.WithdrawResponse{
		Success:      result.Success,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
	}

	if result.Transaction != nil {
		resp.Transaction = domainTransactionToProto(result.Transaction)
		if result.FeeTransaction != nil {
			resp.FeeTransaction = domainTransactionToProto(result.FeeTransaction)
		}
		logger.Info("Withdraw completed", zap.String("transaction_id", result.Transaction.ID), zap.Bool("success", result.Success))
	} else if !result.Success {
		logger.Warn("Withdraw not successful", zap.String("error_code", result.ErrorCode), zap.String("error_message", result.ErrorMessage))
	}

	return resp, nil
}

func (s *TransactionServer) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	logger.Debug("GetTransaction request received", zap.String("transaction_id", req.TransactionId))

	tx, err := s.transactionService.GetTransaction(ctx, req.TransactionId)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			logger.Warn("GetTransaction: transaction not found", zap.String("transaction_id", req.TransactionId))
			return nil, status.Error(codes.NotFound, "transaction not found")
		}
		logger.Error("GetTransaction failed", zap.String("transaction_id", req.TransactionId), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get transaction")
	}

	return &pb.GetTransactionResponse{
		Transaction: domainTransactionToProto(tx),
	}, nil
}

func (s *TransactionServer) GetTransactionHistory(ctx context.Context, req *pb.GetTransactionHistoryRequest) (*pb.GetTransactionHistoryResponse, error) {
	logger.Debug("GetTransactionHistory request received", zap.String("account_id", req.AccountId), zap.Int32("page", req.Page), zap.Int32("page_size", req.PageSize))

	transactions, totalCount, err := s.transactionService.GetTransactionHistory(ctx, req.AccountId, int(req.Page), int(req.PageSize))
	if err != nil {
		logger.Error("GetTransactionHistory failed", zap.String("account_id", req.AccountId), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get transaction history")
	}

	var protoTransactions []*pb.Transaction
	for _, tx := range transactions {
		protoTransactions = append(protoTransactions, domainTransactionToProto(tx))
	}

	logger.Debug("GetTransactionHistory successful", zap.String("account_id", req.AccountId), zap.Int("total_count", totalCount))
	return &pb.GetTransactionHistoryResponse{
		Transactions: protoTransactions,
		TotalCount:   int32(totalCount),
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

func (s *TransactionServer) GetStatement(ctx context.Context, req *pb.GetStatementRequest) (*pb.GetStatementResponse, error) {
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	logger.Debug("GetStatement request received", zap.String("account_id", req.AccountId), zap.Time("start_date", startDate), zap.Time("end_date", endDate))

	transactions, totalCount, openingBalance, closingBalance, err := s.transactionService.GetStatement(
		ctx, req.AccountId, startDate, endDate, int(req.Page), int(req.PageSize),
	)
	if err != nil {
		logger.Error("GetStatement failed", zap.String("account_id", req.AccountId), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get statement")
	}

	var protoTransactions []*pb.Transaction
	for _, tx := range transactions {
		protoTransactions = append(protoTransactions, domainTransactionToProto(tx))
	}

	logger.Debug("GetStatement successful", zap.String("account_id", req.AccountId), zap.Int("total_count", totalCount))
	return &pb.GetStatementResponse{
		AccountId:      req.AccountId,
		OpeningBalance: openingBalance,
		ClosingBalance: closingBalance,
		Transactions:   protoTransactions,
		TotalCount:     int32(totalCount),
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}, nil
}

func domainTransactionToProto(tx *domain.Transaction) *pb.Transaction {
	protoTx := &pb.Transaction{
		Id:                       tx.ID,
		IdempotencyKey:           tx.IdempotencyKey,
		Type:                     domainTypeToProto(tx.Type),
		Status:                   domainStatusToProto(tx.Status),
		SourceAccountId:          tx.SourceAccountID,
		SourceAccountNumber:      tx.SourceAccountNumber,
		DestinationAccountId:     tx.DestinationAccountID,
		DestinationAccountNumber: tx.DestinationAccountNumber,
		Amount:                   tx.Amount.String(),
		Currency:                 tx.Currency,
		Description:              tx.Description,
		CreatedAt:                timestamppb.New(tx.CreatedAt),
	}

	if tx.ProcessedAt != nil {
		protoTx.ProcessedAt = timestamppb.New(*tx.ProcessedAt)
	}

	return protoTx
}

func domainTypeToProto(t domain.TransactionType) pb.TransactionType {
	switch t {
	case domain.TransactionTypeTransfer:
		return pb.TransactionType_TRANSACTION_TYPE_TRANSFER
	case domain.TransactionTypeDeposit:
		return pb.TransactionType_TRANSACTION_TYPE_DEPOSIT
	case domain.TransactionTypeWithdrawal:
		return pb.TransactionType_TRANSACTION_TYPE_WITHDRAWAL
	case domain.TransactionTypeFee:
		return pb.TransactionType_TRANSACTION_TYPE_FEE
	default:
		return pb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}
}

func domainStatusToProto(s domain.TransactionStatus) pb.TransactionStatus {
	switch s {
	case domain.TransactionStatusPending:
		return pb.TransactionStatus_TRANSACTION_STATUS_PENDING
	case domain.TransactionStatusCompleted:
		return pb.TransactionStatus_TRANSACTION_STATUS_COMPLETED
	case domain.TransactionStatusFailed:
		return pb.TransactionStatus_TRANSACTION_STATUS_FAILED
	case domain.TransactionStatusReversed:
		return pb.TransactionStatus_TRANSACTION_STATUS_REVERSED
	default:
		return pb.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
	}
}
