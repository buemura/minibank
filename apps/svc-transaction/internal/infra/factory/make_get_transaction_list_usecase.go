package factory

import (
	"github.com/buemura/minibank/packages/cache"
	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/core/usecase"
	"github.com/buemura/minibank/svc-transaction/internal/infra/database"
)

func MakeGetTransactionListUsecase() *usecase.GetTransactionList {
	trxRepo := database.NewSqlTransactionRepository()
	cacheRepo := cache.NewRedisCacheRepository(config.REDIS_URL, config.REDIS_PASSWORD)
	return usecase.NewGetTransactionList(cacheRepo, trxRepo)
}
