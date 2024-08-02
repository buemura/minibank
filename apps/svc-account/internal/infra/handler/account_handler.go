package handler

import (
	"context"
	"log"
	"time"

	"github.com/buemura/minibank/packages/pb"
	"github.com/buemura/minibank/svc-account/internal/application/usecase"
	"github.com/buemura/minibank/svc-account/internal/infra/database"
)

type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
}

func (c AccountHandler) GetAccount(
	ctx context.Context,
	in *pb.GetAccountRequest,
) (*pb.Account, error) {
	log.Println("[AccountHandler][GetAccount] - Incoming request")

	accountRepo := database.NewPgxOrderRepository()
	getAccountUC := usecase.NewGetAccount(accountRepo)

	acc, err := getAccountUC.Execute(in.Id)
	if err != nil {
		log.Println("[AccountHandler][GetAccount] - Error:", err.Error())
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
