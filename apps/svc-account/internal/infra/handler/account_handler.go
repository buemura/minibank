package handler

import (
	"context"
	"time"

	"github.com/buemura/minibank/packages/gen/protos"
	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
	"github.com/buemura/minibank/svc-account/internal/core/usecase"
	"github.com/buemura/minibank/svc-account/internal/infra/cache"
	"github.com/buemura/minibank/svc-account/internal/infra/database"
	"golang.org/x/exp/slog"
)

type AccountHandler struct {
	protos.UnimplementedAccountServiceServer
}

func (c AccountHandler) GetAccount(
	ctx context.Context,
	in *protos.GetAccountRequest,
) (*protos.Account, error) {
	slog.Info("[AccountHandler][GetAccount] - Incoming request")

	accountRepo := database.NewPgxAccountRepository()
	cacheRepo := cache.NewRedisCacheRepository()
	getAccountUC := usecase.NewGetAccount(cacheRepo, accountRepo)

	acc, err := getAccountUC.Execute(in.Id)
	if err != nil {
		slog.Error("[AccountHandler][GetAccount] - Error:", err.Error())
		return nil, HandleGrpcError(err)
	}

	return &protos.Account{
		Id:            acc.ID,
		Balance:       int32(acc.Balance),
		OwnerName:     acc.OwnerName,
		OwnerDocument: acc.OwnerDocument,
		Status:        acc.Status,
		CreatedAt:     acc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     acc.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (c AccountHandler) CreateAccount(
	ctx context.Context,
	in *protos.CreateAccountRequest,
) (*protos.Account, error) {
	slog.Info("[AccountHandler][CreateAccount] - Incoming request")

	accountRepo := database.NewPgxAccountRepository()
	cacheRepo := cache.NewRedisCacheRepository()
	createAccountUC := usecase.NewCreateAccount(cacheRepo, accountRepo)

	createAccIn := &account.CreateAccountIn{
		OwnerName:     in.OwnerName,
		OwnerDocument: in.OwnerDocument,
	}

	acc, err := createAccountUC.Execute(createAccIn)
	if err != nil {
		slog.Error("[AccountHandler][CreateAccount] - Error:", err.Error())
		return nil, HandleGrpcError(err)
	}

	return &protos.Account{
		Id:            acc.ID,
		Balance:       int32(acc.Balance),
		OwnerName:     acc.OwnerName,
		OwnerDocument: acc.OwnerDocument,
		Status:        acc.Status,
		CreatedAt:     acc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     acc.UpdatedAt.Format(time.RFC3339),
	}, nil
}
