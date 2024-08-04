package grpc

import (
	"context"
	"log"
	"time"

	"github.com/buemura/minibank-api-gateway/config"
	"github.com/buemura/minibank-api-gateway/internal/core/domain/account"
	"github.com/buemura/minibank/packages/gen/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcAccountService struct{}

func NewGrpcAccountService() *GrpcAccountService {
	return &GrpcAccountService{}
}

func (*GrpcAccountService) GetAccount(id string) (*account.Account, error) {
	conn, err := grpc.Dial(config.GRPC_HOST_SVC_ACCOUNT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial svc-account server: %v", err)
	}
	defer conn.Close()

	client := protos.NewAccountServiceClient(conn)
	log.Println("[GrpcClient][GetAccount] - Request Get Account for:", id)

	request := &protos.GetAccountRequest{Id: id}
	res, err := client.GetAccount(context.Background(), request)
	if err != nil {
		log.Println("[GrpcClient][GetCustomer] - Error:", err)
		return nil, err
	}

	createdAt, _ := time.Parse("2006-01-02T15:04:05.000Z", res.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02T15:04:05.000Z", res.UpdatedAt)

	return &account.Account{
		ID:            res.Id,
		OwnerName:     res.OwnerName,
		OwnerDocument: res.OwnerDocument,
		Balance:       int(res.Balance),
		Status:        res.Status,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}
