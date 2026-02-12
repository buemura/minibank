package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/buemura/minibank/packages/metrics"
	"github.com/buemura/minibank/packages/tracing"
	"github.com/buemura/minibank/svc-auth/internal/config"
	"github.com/buemura/minibank/svc-auth/internal/database"
	authgrpc "github.com/buemura/minibank/svc-auth/internal/grpc"
	"github.com/buemura/minibank/svc-auth/internal/repository/postgres"
	"github.com/buemura/minibank/svc-auth/internal/service"
	pb "github.com/buemura/minibank/svc-auth/proto/auth/v1"
	"github.com/buemura/minibank/packages/jwt"
	"github.com/buemura/minibank/packages/logger"
)

func main() {
	logger.Init("svc-auth", true)
	defer logger.Sync()

	cfg := config.Load()

	if err := database.RunMigrations(cfg.DatabaseURL(), "migrations"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}
	logger.Info("migrations completed")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdownTracer, err := tracing.Init(ctx, "svc-auth", cfg.OTLPEndpoint)
	if err != nil {
		logger.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			logger.Error("failed to shutdown tracer", zap.Error(err))
		}
	}()

	shutdownMeter, err := metrics.Init(ctx)
	if err != nil {
		logger.Fatal("failed to initialize metrics", zap.Error(err))
	}
	defer func() {
		if err := shutdownMeter(context.Background()); err != nil {
			logger.Error("failed to shutdown meter provider", zap.Error(err))
		}
	}()

	go func() {
		logger.Info("starting metrics server", zap.String("port", cfg.MetricsPort))
		if err := metrics.ServeMetrics(cfg.MetricsPort); err != nil {
			logger.Error("metrics server error", zap.Error(err))
		}
	}()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}
	logger.Info("connected to database")

	userRepo := postgres.NewUserRepository(pool)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(pool)

	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtManager, cfg.RefreshTokenExpiry)

	authServer := authgrpc.NewAuthServer(authService)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

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
