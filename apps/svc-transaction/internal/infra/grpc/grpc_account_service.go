package grpc

import (
	"context"
	"log"
	"time"

	"github.com/buemura/minibank/packages/gen/protos"
	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/domain/account"
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

func (*GrpcAccountService) UpdateBalance(id string, newBalance int) (*account.Account, error) {
	conn, err := grpc.Dial(config.GRPC_HOST_SVC_ACCOUNT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial svc-account server: %v", err)
	}
	defer conn.Close()

	client := protos.NewAccountServiceClient(conn)
	log.Println("[GrpcClient][UpdateBalance] - Request Update Account for:", id)

	request := &protos.UpdateBalanceRequest{Id: id, NewBalance: int32(newBalance)}
	res, err := client.UpdateBalance(context.Background(), request)
	if err != nil {
		log.Println("[GrpcClient][UpdateBalance] - Error:", err)
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
