package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/buemura/minibank/transaction-service/internal/config"
	"github.com/buemura/minibank/transaction-service/internal/database"
	transactiongrpc "github.com/buemura/minibank/transaction-service/internal/grpc"
	"github.com/buemura/minibank/transaction-service/internal/repository/postgres"
	"github.com/buemura/minibank/transaction-service/internal/service"
	accountpb "github.com/buemura/minibank/transaction-service/proto/account/v1"
	pb "github.com/buemura/minibank/transaction-service/proto/transaction/v1"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.Load()

	if err := database.RunMigrations(cfg.DatabaseURL(), "migrations"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}
	logger.Info("migrations completed")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}
	logger.Info("connected to database")

	accountConn, err := grpc.Dial(cfg.AccountServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to connect to account service", zap.Error(err))
	}
	defer accountConn.Close()

	accountClient := accountpb.NewAccountServiceClient(accountConn)

	transactionRepo := postgres.NewTransactionRepository(pool)
	idempotencyRepo := postgres.NewIdempotencyRepository(pool)

	transactionService := service.NewTransactionService(transactionRepo, idempotencyRepo, accountClient)
	transactionServer := transactiongrpc.NewTransactionServer(transactionService)

	grpcServer := grpc.NewServer()
	pb.RegisterTransactionServiceServer(grpcServer, transactionServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	grpcServer.GracefulStop()
	logger.Info("server stopped")
}
