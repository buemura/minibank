package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/buemura/minibank/auth-service/internal/domain"
	"github.com/buemura/minibank/auth-service/internal/service"
	pb "github.com/buemura/minibank/auth-service/proto/auth/v1"
	"github.com/buemura/minibank/packages/logger"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthServer(authService *service.AuthService) *AuthServer {
	return &AuthServer{authService: authService}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	logger.Info("Register request received", zap.String("email", req.Email))

	result, err := s.authService.Register(ctx, service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			logger.Warn("Registration failed: user already exists", zap.String("email", req.Email))
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		logger.Error("Registration failed", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	logger.Info("User registered successfully", zap.String("email", req.Email), zap.String("user_id", result.User.ID))
	return &pb.RegisterResponse{
		User:         domainUserToProto(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logger.Info("Login request received", zap.String("email", req.Email))

	result, err := s.authService.Login(ctx, service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			logger.Warn("Login failed: invalid credentials", zap.String("email", req.Email))
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		logger.Error("Login failed", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to login")
	}

	logger.Info("User logged in successfully", zap.String("email", req.Email), zap.String("user_id", result.User.ID))
	return &pb.LoginResponse{
		User:         domainUserToProto(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	logger.Debug("RefreshToken request received")

	result, err := s.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			logger.Warn("RefreshToken failed: invalid token")
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		case errors.Is(err, service.ErrTokenExpired):
			logger.Warn("RefreshToken failed: token expired")
			return nil, status.Error(codes.Unauthenticated, "token expired")
		case errors.Is(err, service.ErrTokenRevoked):
			logger.Warn("RefreshToken failed: token revoked")
			return nil, status.Error(codes.Unauthenticated, "token revoked")
		default:
			logger.Error("RefreshToken failed", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to refresh token")
		}
	}

	logger.Debug("Token refreshed successfully")
	return &pb.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	logger.Debug("ValidateToken request received")

	claims, err := s.authService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		logger.Debug("Token validation failed", zap.Error(err))
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	logger.Debug("Token validated successfully", zap.String("user_id", claims.UserID))
	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	logger.Debug("Logout request received")

	err := s.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		logger.Error("Logout failed", zap.Error(err))
		return &pb.LogoutResponse{Success: false}, nil
	}

	logger.Debug("Logout successful")
	return &pb.LogoutResponse{Success: true}, nil
}

func (s *AuthServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	logger.Debug("GetUserProfile request received", zap.String("user_id", req.UserId))

	user, err := s.authService.GetUserByID(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			logger.Warn("GetUserProfile: user not found", zap.String("user_id", req.UserId))
			return nil, status.Error(codes.NotFound, "user not found")
		}
		logger.Error("GetUserProfile failed", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	logger.Debug("GetUserProfile successful", zap.String("user_id", req.UserId))
	return &pb.GetUserProfileResponse{
		User: domainUserToProto(user),
	}, nil
}

func domainUserToProto(user *domain.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Status:    string(user.Status),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
