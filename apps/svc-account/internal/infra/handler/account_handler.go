package handler

import (
	"context"
	"time"

	"github.com/buemura/minibank/packages/pb"
	"github.com/buemura/minibank/svc-account/internal/application/usecase"
	"github.com/buemura/minibank/svc-account/internal/infra/cache"
	"github.com/buemura/minibank/svc-account/internal/infra/database"
	"golang.org/x/exp/slog"
)

type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
}

func (c AccountHandler) GetAccount(
	ctx context.Context,
	in *pb.GetAccountRequest,
) (*pb.Account, error) {
	slog.Info("[AccountHandler][GetAccount] - Incoming request")

	accountRepo := database.NewPgxOrderRepository()
	cacheRepo := cache.NewRedisCacheRepository()
	getAccountUC := usecase.NewGetAccount(cacheRepo, accountRepo)

	acc, err := getAccountUC.Execute(in.Id)
	if err != nil {
		slog.Error("[AccountHandler][GetAccount] - Error:", err.Error())
		return nil, HandleGrpcError(err)
	}

	return &pb.Account{
		Id:            acc.ID,
		Balance:       int32(acc.Balance),
		OwnerName:     acc.OwnerName,
		OwnerDocument: acc.OwnerDocument,
		Status:        acc.Status,
		CreatedAt:     acc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     acc.UpdatedAt.Format(time.RFC3339),
	}, nil
}
