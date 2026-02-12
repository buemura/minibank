package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/buemura/minibank/api-gtw/docs"
	"github.com/buemura/minibank/api-gtw/internal/config"
	"github.com/buemura/minibank/api-gtw/internal/handlers"
	"github.com/buemura/minibank/api-gtw/internal/middleware"
	accountpb "github.com/buemura/minibank/api-gtw/proto/account/v1"
	authpb "github.com/buemura/minibank/api-gtw/proto/auth/v1"
	transactionpb "github.com/buemura/minibank/api-gtw/proto/transaction/v1"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.Load()

	authConn, err := grpc.Dial(cfg.AuthServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to connect to auth service", zap.Error(err))
	}
	defer authConn.Close()

	accountConn, err := grpc.Dial(cfg.AccountServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to connect to account service", zap.Error(err))
	}
	defer accountConn.Close()

	transactionConn, err := grpc.Dial(cfg.TransactionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to connect to transaction service", zap.Error(err))
	}
	defer transactionConn.Close()

	authClient := authpb.NewAuthServiceClient(authConn)
	accountClient := accountpb.NewAccountServiceClient(accountConn)
	transactionClient := transactionpb.NewTransactionServiceClient(transactionConn)

	authHandler := handlers.NewAuthHandler(authClient)
	accountHandler := handlers.NewAccountHandler(accountClient, authClient)
	transactionHandler := handlers.NewTransactionHandler(transactionClient)

	authMiddleware := middleware.NewAuthMiddleware(authClient)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/docs", docs.ScalarHandler())
	router.GET("/docs/openapi.yaml", docs.SpecHandler())

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", authMiddleware.Authenticate(), authHandler.Me)
		}

		accounts := v1.Group("/accounts")
		accounts.Use(authMiddleware.Authenticate())
		{
			accounts.GET("", accountHandler.GetAccounts)
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/lookup", accountHandler.LookupAccountByNumber)
			accounts.GET("/:id", accountHandler.GetAccount)
			accounts.GET("/:id/balance", accountHandler.GetBalance)
			accounts.POST("/:id/transfers", transactionHandler.Transfer)
			accounts.POST("/:id/deposits", transactionHandler.Deposit)
			accounts.GET("/:id/transactions", transactionHandler.GetTransactionHistory)
			accounts.GET("/:id/statement", transactionHandler.GetStatement)
		}

		transactions := v1.Group("/transactions")
		transactions.Use(authMiddleware.Authenticate())
		{
			transactions.GET("/:id", transactionHandler.GetTransaction)
		}
	}

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Info("starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped")
}
