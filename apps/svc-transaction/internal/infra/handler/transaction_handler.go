package handler

import (
	"context"
	"time"

	"github.com/buemura/minibank/packages/gen/protos"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"github.com/buemura/minibank/svc-transaction/internal/core/usecase"
	"github.com/buemura/minibank/svc-transaction/internal/infra/cache"
	"github.com/buemura/minibank/svc-transaction/internal/infra/database"
	"golang.org/x/exp/slog"
)

type TransactionHandler struct {
	protos.UnimplementedTransactionServiceServer
}

func (c TransactionHandler) GetTransactions(
	ctx context.Context,
	in *protos.GetTransactionsRequest,
) (*protos.GetTransactionsResponse, error) {
	slog.Info("[TransactionHandler][GetTransactions] - Incoming request")

	trxRepo := database.NewPgxTransactionRepository()
	cacheRepo := cache.NewRedisCacheRepository()
	getTrxListUC := usecase.NewGetTransactionList(cacheRepo, trxRepo)

	input := &transaction.GetTransactionListIn{
		AccountID: in.AccountId,
		Page:      int(in.Page),
		Items:     int(in.Items),
	}

	res, err := getTrxListUC.Execute(input)
	if err != nil {
		slog.Error("[TransactionHandler][GetTransactions] - Error:", err.Error())
		return nil, HandleGrpcError(err)
	}

	var trxResponse []*protos.Transaction
	for _, trx := range res.Data {
		trxResponse = append(trxResponse, &protos.Transaction{
			Id:                   trx.ID,
			AccountId:            trx.AccountID,
			DestinationAccountId: trx.DestinationAccountID,
			Amount:               int64(trx.Amount),
			Status:               string(trx.Status),
			TransactionType:      string(trx.TransactionType),
			CreatedAt:            trx.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            trx.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &protos.GetTransactionsResponse{
		Data: trxResponse,
		Meta: &protos.GetTransactionsResponse_Meta{
			Page:       int32(res.Meta.Page),
			Items:      int32(res.Meta.Items),
			TotalPages: int32(res.Meta.TotalPages),
			TotalItems: int32(res.Meta.TotalItems),
		},
	}, nil
}

func (c TransactionHandler) CreateTransaction(
	ctx context.Context,
	in *protos.CreateTransactionRequest,
) (*protos.Transaction, error) {
	slog.Info("[TransactionHandler][CreateTransaction] - Incoming request")

	trxRepo := database.NewPgxTransactionRepository()
	cacheRepo := cache.NewRedisCacheRepository()
	createTransactionUC := usecase.NewCreateTransaction(cacheRepo, trxRepo)

	createAccIn := &transaction.CreateTransactionIn{
		AccountID:            in.AccountId,
		DestinationAccountID: in.DestinationAccountId,
		Amount:               int(in.Amount),
		TransactionType:      transaction.TransactionType(in.TransactionType),
	}

	trx, err := createTransactionUC.Execute(createAccIn)
	if err != nil {
		slog.Error("[TransactionHandler][CreateAccount] - Error:", err.Error())
		return nil, HandleGrpcError(err)
	}

	return &protos.Transaction{
		Id:                   trx.ID,
		AccountId:            trx.AccountID,
		DestinationAccountId: trx.DestinationAccountID,
		Amount:               int64(trx.Amount),
		Status:               string(trx.Status),
		TransactionType:      string(trx.TransactionType),
		CreatedAt:            trx.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            trx.UpdatedAt.Format(time.RFC3339),
	}, nil
}
