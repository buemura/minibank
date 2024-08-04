package factory

import (
	"github.com/buemura/minibank/packages/cache"
	"github.com/buemura/minibank/svc-account/config"
	"github.com/buemura/minibank/svc-account/internal/core/usecase"
	"github.com/buemura/minibank/svc-account/internal/infra/database"
)

func MakeCreateAccountUsecase() *usecase.CreateAccount {
	accountRepo := database.NewPgxAccountRepository()
	cacheRepo := cache.NewRedisCacheRepository(config.REDIS_URL, config.REDIS_PASSWORD)
	return usecase.NewCreateAccount(cacheRepo, accountRepo)
}
