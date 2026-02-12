package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/buemura/minibank/packages/jwt"
	"github.com/buemura/minibank/packages/password"
	"github.com/buemura/minibank/svc-auth/internal/domain"
	"github.com/buemura/minibank/svc-auth/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestAuthService() (*AuthService, *mocks.MockUserRepository, *mocks.MockRefreshTokenRepository) {
	userRepo := new(mocks.MockUserRepository)
	refreshRepo := new(mocks.MockRefreshTokenRepository)
	jwtManager := jwt.NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(userRepo, refreshRepo, jwtManager, 7*24*time.Hour)
	return svc, userRepo, refreshRepo
}

func newActiveUser(id, email string) *domain.User {
	hash, _ := password.Hash("password123")
	return &domain.User{
		ID:           id,
		Email:        email,
		PasswordHash: hash,
		FullName:     "Test User",
		Phone:        "+5511999999999",
		Status:       domain.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestRegister_Success(t *testing.T) {
	svc, userRepo, refreshRepo := newTestAuthService()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	refreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	result, err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		Phone:    "+5511999999999",
	})

	require.NoError(t, err)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Greater(t, result.ExpiresIn, int64(0))
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	existing := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existing, nil)

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	})

	assert.ErrorIs(t, err, ErrUserAlreadyExists)
	userRepo.AssertExpectations(t)
}

func TestRegister_GetByEmailError(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("db error"))

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	})

	assert.EqualError(t, err, "db error")
	userRepo.AssertExpectations(t)
}

func TestRegister_CreateUserError(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	userRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	})

	assert.EqualError(t, err, "db error")
	userRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	svc, userRepo, refreshRepo := newTestAuthService()

	user := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	refreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	result, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.User.ID)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	userRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	user := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "wrong-password",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	userRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	user := newActiveUser("user-1", "test@example.com")
	user.Status = domain.UserStatusInactive
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	svc, userRepo, refreshRepo := newTestAuthService()

	// Generate a token pair first to get a valid refresh token
	jwtManager := jwt.NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	tokenPair, _ := jwtManager.GenerateTokenPair("user-1", "test@example.com")
	tokenHash := sha256Hash(tokenPair.RefreshToken)

	storedToken := &domain.RefreshToken{
		ID:        "rt-1",
		UserID:    "user-1",
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	user := newActiveUser("user-1", "test@example.com")

	refreshRepo.On("GetByTokenHash", mock.Anything, tokenHash).Return(storedToken, nil)
	refreshRepo.On("RevokeByTokenHash", mock.Anything, tokenHash).Return(nil)
	refreshRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
	userRepo.On("GetByID", mock.Anything, "user-1").Return(user, nil)

	result, err := svc.RefreshToken(context.Background(), tokenPair.RefreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	refreshRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	svc, _, refreshRepo := newTestAuthService()

	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, nil)

	_, err := svc.RefreshToken(context.Background(), "invalid-token")

	assert.ErrorIs(t, err, ErrInvalidToken)
	refreshRepo.AssertExpectations(t)
}

func TestRefreshToken_Revoked(t *testing.T) {
	svc, _, refreshRepo := newTestAuthService()

	storedToken := &domain.RefreshToken{
		ID:        "rt-1",
		UserID:    "user-1",
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   true,
	}
	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(storedToken, nil)

	_, err := svc.RefreshToken(context.Background(), "some-token")

	assert.ErrorIs(t, err, ErrTokenRevoked)
	refreshRepo.AssertExpectations(t)
}

func TestRefreshToken_Expired(t *testing.T) {
	svc, _, refreshRepo := newTestAuthService()

	storedToken := &domain.RefreshToken{
		ID:        "rt-1",
		UserID:    "user-1",
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Revoked:   false,
	}
	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(storedToken, nil)

	_, err := svc.RefreshToken(context.Background(), "some-token")

	assert.ErrorIs(t, err, ErrTokenExpired)
	refreshRepo.AssertExpectations(t)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	svc, userRepo, refreshRepo := newTestAuthService()

	storedToken := &domain.RefreshToken{
		ID:        "rt-1",
		UserID:    "user-1",
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(storedToken, nil)
	userRepo.On("GetByID", mock.Anything, "user-1").Return(nil, nil)

	_, err := svc.RefreshToken(context.Background(), "some-token")

	assert.ErrorIs(t, err, ErrUserNotFound)
	refreshRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestValidateToken_Success(t *testing.T) {
	svc, _, _ := newTestAuthService()

	jwtManager := jwt.NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	tokenPair, _ := jwtManager.GenerateTokenPair("user-1", "test@example.com")

	claims, err := svc.ValidateToken(context.Background(), tokenPair.AccessToken)

	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestValidateToken_Invalid(t *testing.T) {
	svc, _, _ := newTestAuthService()

	_, err := svc.ValidateToken(context.Background(), "invalid-token")

	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestValidateToken_Expired(t *testing.T) {
	svc, _, _ := newTestAuthService()

	// Create a token that's already expired
	expiredManager := jwt.NewManager("test-secret", -1*time.Hour, 7*24*time.Hour)
	tokenPair, _ := expiredManager.GenerateTokenPair("user-1", "test@example.com")

	_, err := svc.ValidateToken(context.Background(), tokenPair.AccessToken)

	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestLogout_Success(t *testing.T) {
	svc, _, refreshRepo := newTestAuthService()

	refreshRepo.On("RevokeByTokenHash", mock.Anything, mock.Anything).Return(nil)

	err := svc.Logout(context.Background(), "some-refresh-token")

	assert.NoError(t, err)
	refreshRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	user := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByID", mock.Anything, "user-1").Return(user, nil)

	result, err := svc.GetUserByID(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	userRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	svc, userRepo, _ := newTestAuthService()

	userRepo.On("GetByID", mock.Anything, "user-1").Return(nil, nil)

	_, err := svc.GetUserByID(context.Background(), "user-1")

	assert.ErrorIs(t, err, ErrUserNotFound)
	userRepo.AssertExpectations(t)
}

func sha256Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
