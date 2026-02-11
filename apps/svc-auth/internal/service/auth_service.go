package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/buemura/minibank/svc-auth/internal/domain"
	"github.com/buemura/minibank/svc-auth/internal/repository"
	"github.com/buemura/minibank/packages/jwt"
	"github.com/buemura/minibank/packages/logger"
	"github.com/buemura/minibank/packages/password"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token revoked")
)

type AuthService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *jwt.Manager
	refreshExpiry    time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	refreshExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		refreshExpiry:    refreshExpiry,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type AuthResult struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	logger.Debug("AuthService.Register: checking if user exists", zap.String("email", input.Email))

	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		logger.Error("AuthService.Register: failed to check existing user", zap.String("email", input.Email), zap.Error(err))
		return nil, err
	}
	if existing != nil {
		logger.Debug("AuthService.Register: user already exists", zap.String("email", input.Email))
		return nil, ErrUserAlreadyExists
	}

	logger.Debug("AuthService.Register: hashing password")
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		logger.Error("AuthService.Register: failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		FullName:     input.FullName,
		Phone:        input.Phone,
		Status:       domain.UserStatusActive,
	}

	logger.Debug("AuthService.Register: creating user in database", zap.String("email", input.Email))
	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("AuthService.Register: failed to create user in database", zap.String("email", input.Email), zap.Error(err))
		return nil, err
	}

	logger.Debug("AuthService.Register: generating tokens", zap.String("user_id", user.ID))
	result, err := s.generateTokens(ctx, user)
	if err != nil {
		logger.Error("AuthService.Register: failed to generate tokens", zap.String("user_id", user.ID), zap.Error(err))
		return nil, err
	}

	logger.Info("AuthService.Register: user registered successfully", zap.String("user_id", user.ID), zap.String("email", input.Email))
	return result, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	logger.Debug("AuthService.Login: looking up user", zap.String("email", input.Email))

	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		logger.Error("AuthService.Login: failed to get user by email", zap.String("email", input.Email), zap.Error(err))
		return nil, err
	}
	if user == nil {
		logger.Debug("AuthService.Login: user not found", zap.String("email", input.Email))
		return nil, ErrInvalidCredentials
	}

	if !password.Verify(input.Password, user.PasswordHash) {
		logger.Debug("AuthService.Login: invalid password", zap.String("email", input.Email))
		return nil, ErrInvalidCredentials
	}

	if user.Status != domain.UserStatusActive {
		logger.Debug("AuthService.Login: user not active", zap.String("email", input.Email), zap.String("status", string(user.Status)))
		return nil, ErrInvalidCredentials
	}

	logger.Debug("AuthService.Login: generating tokens", zap.String("user_id", user.ID))
	return s.generateTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	tokenHash := hashToken(refreshToken)

	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, ErrInvalidToken
	}

	if storedToken.Revoked {
		return nil, ErrTokenRevoked
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := s.refreshTokenRepo.RevokeByTokenHash(ctx, tokenHash); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) ValidateToken(ctx context.Context, accessToken string) (*jwt.Claims, error) {
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return s.refreshTokenRepo.RevokeByTokenHash(ctx, tokenHash)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *AuthService) generateTokens(ctx context.Context, user *domain.User) (*AuthResult, error) {
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := hashToken(tokenPair.RefreshToken)
	refreshTokenRecord := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(s.refreshExpiry),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
