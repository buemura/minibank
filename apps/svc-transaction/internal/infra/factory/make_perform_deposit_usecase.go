package factory

import (
	"github.com/buemura/minibank/packages/cache"
	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/usecase"
	"github.com/buemura/minibank/svc-transaction/internal/infra/database"
	"github.com/buemura/minibank/svc-transaction/internal/infra/grpc"
)

func MakePerformDepositUsecase() *usecase.PerformDeposit {
	accService := grpc.NewGrpcAccountService()
	trxRepo := database.NewSqlTransactionRepository()
	cacheRepo := cache.NewRedisCacheRepository(config.REDIS_URL, config.REDIS_PASSWORD)
	return usecase.NewPerformDeposit(accService, cacheRepo, trxRepo)
}
