package grpc

import (
	"context"
	"log"
	"time"

	"github.com/buemura/minibank/api-gateway/config"
	"github.com/buemura/minibank/api-gateway/internal/core/domain/common"
	"github.com/buemura/minibank/api-gateway/internal/core/domain/transaction"
	"github.com/buemura/minibank/packages/gen/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcTransactionService struct{}

func NewGrpcTransactionService() *GrpcTransactionService {
	return &GrpcTransactionService{}
}

func (*GrpcTransactionService) GetTransactions(in *transaction.GetTransactionListIn) (*transaction.PaginatedTransactionList, error) {
	conn, err := grpc.Dial(config.GRPC_HOST_SVC_TRANSACTION, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial svc-transaction server: %v", err)
	}
	defer conn.Close()

	client := protos.NewTransactionServiceClient(conn)
	log.Println("[GrpcClient][GetTransactions] - Request Get Transactions for account:", in.AccountID)

	request := &protos.GetTransactionsRequest{AccountId: in.AccountID, Page: int32(in.Page), Items: int32(in.Items)}
	res, err := client.GetTransactions(context.Background(), request)
	if err != nil {
		log.Println("[GrpcClient][GetTransactions] - Error:", err)
		return nil, err
	}

	var data []*transaction.Transaction
	for _, d := range res.Data {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.000Z", d.CreatedAt)
		updatedAt, _ := time.Parse("2006-01-02T15:04:05.000Z", d.UpdatedAt)

		data = append(data, &transaction.Transaction{
			ID:                   d.Id,
			AccountID:            d.AccountId,
			DestinationAccountID: d.DestinationAccountId,
			Amount:               int(d.Amount),
			Status:               transaction.TransactionStatus(d.Status),
			TransactionType:      transaction.TransactionType(d.TransactionType),
			CreatedAt:            createdAt,
			UpdatedAt:            updatedAt,
		})
	}

	return &transaction.PaginatedTransactionList{
		Data: data,
		Meta: common.PaginationMetaOut{
			Page:       int(res.Meta.Page),
			Items:      int(res.Meta.Items),
			TotalPages: int(res.Meta.TotalPages),
			TotalItems: int(res.Meta.TotalItems),
		},
	}, nil
}

func (*GrpcTransactionService) CreateTransaction(in *transaction.CreateTransactionIn) (*transaction.Transaction, error) {
	conn, err := grpc.Dial(config.GRPC_HOST_SVC_TRANSACTION, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial svc-transaction server: %v", err)
	}
	defer conn.Close()

	client := protos.NewTransactionServiceClient(conn)
	log.Println("[GrpcClient][CreateTransaction] - Request Get Transactions for account:", in.AccountID)

	request := &protos.CreateTransactionRequest{
		AccountId:            in.AccountID,
		Amount:               int64(in.Amount),
		DestinationAccountId: in.DestinationAccountID,
		TransactionType:      string(in.TransactionType),
	}
	res, err := client.CreateTransaction(context.Background(), request)
	if err != nil {
		log.Println("[GrpcClient][CreateTransaction] - Error:", err)
		return nil, err
	}

	createdAt, _ := time.Parse("2006-01-02T15:04:05.000Z", res.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02T15:04:05.000Z", res.UpdatedAt)

	return &transaction.Transaction{
		ID:              res.Id,
		AccountID:       res.AccountId,
		Amount:          int(res.Amount),
		Status:          transaction.TransactionStatus(res.Status),
		TransactionType: transaction.TransactionType(res.TransactionType),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}
