package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/buemura/minibank/packages/jwt"
	"github.com/buemura/minibank/packages/password"
	"github.com/buemura/minibank/svc-auth/internal/domain"
	"github.com/buemura/minibank/svc-auth/internal/repository/mocks"
	"github.com/buemura/minibank/svc-auth/internal/service"
	pb "github.com/buemura/minibank/svc-auth/proto/auth/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestServer() (*AuthServer, *mocks.MockUserRepository, *mocks.MockRefreshTokenRepository) {
	userRepo := new(mocks.MockUserRepository)
	refreshRepo := new(mocks.MockRefreshTokenRepository)
	jwtManager := jwt.NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	svc := service.NewAuthService(userRepo, refreshRepo, jwtManager, 7*24*time.Hour)
	server := NewAuthServer(svc)
	return server, userRepo, refreshRepo
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

func TestRegister_gRPC_Success(t *testing.T) {
	server, userRepo, refreshRepo := newTestServer()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	userRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	refreshRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := server.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		Phone:    "+5511999999999",
	})

	require.NoError(t, err)
	assert.NotNil(t, resp.User)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestRegister_gRPC_AlreadyExists(t *testing.T) {
	server, userRepo, _ := newTestServer()

	existing := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existing, nil)

	_, err := server.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	userRepo.AssertExpectations(t)
}

func TestRegister_gRPC_InternalError(t *testing.T) {
	server, userRepo, _ := newTestServer()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("db error"))

	_, err := server.Register(context.Background(), &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	userRepo.AssertExpectations(t)
}

func TestLogin_gRPC_Success(t *testing.T) {
	server, userRepo, refreshRepo := newTestServer()

	user := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	refreshRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := server.Login(context.Background(), &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotNil(t, resp.User)
	assert.NotEmpty(t, resp.AccessToken)
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestLogin_gRPC_InvalidCredentials(t *testing.T) {
	server, userRepo, _ := newTestServer()

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

	_, err := server.Login(context.Background(), &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_gRPC_InvalidToken(t *testing.T) {
	server, _, refreshRepo := newTestServer()

	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, nil)

	_, err := server.RefreshToken(context.Background(), &pb.RefreshTokenRequest{
		RefreshToken: "invalid-token",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "invalid token")
	refreshRepo.AssertExpectations(t)
}

func TestRefreshToken_gRPC_Expired(t *testing.T) {
	server, _, refreshRepo := newTestServer()

	storedToken := &domain.RefreshToken{
		ID:        "rt-1",
		UserID:    "user-1",
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Revoked:   false,
	}
	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(storedToken, nil)

	_, err := server.RefreshToken(context.Background(), &pb.RefreshTokenRequest{
		RefreshToken: "some-token",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "token expired")
	refreshRepo.AssertExpectations(t)
}

func TestRefreshToken_gRPC_Revoked(t *testing.T) {
	server, _, refreshRepo := newTestServer()

	storedToken := &domain.RefreshToken{
		ID:        "rt-1",
		UserID:    "user-1",
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   true,
	}
	refreshRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(storedToken, nil)

	_, err := server.RefreshToken(context.Background(), &pb.RefreshTokenRequest{
		RefreshToken: "some-token",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "token revoked")
	refreshRepo.AssertExpectations(t)
}

func TestValidateToken_gRPC_Valid(t *testing.T) {
	server, _, _ := newTestServer()

	jwtManager := jwt.NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	tokenPair, _ := jwtManager.GenerateTokenPair("user-1", "test@example.com")

	resp, err := server.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		AccessToken: tokenPair.AccessToken,
	})

	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "user-1", resp.UserId)
	assert.Equal(t, "test@example.com", resp.Email)
}

func TestValidateToken_gRPC_Invalid(t *testing.T) {
	server, _, _ := newTestServer()

	resp, err := server.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		AccessToken: "invalid-token",
	})

	require.NoError(t, err)
	assert.False(t, resp.Valid)
}

func TestLogout_gRPC_Success(t *testing.T) {
	server, _, refreshRepo := newTestServer()

	refreshRepo.On("RevokeByTokenHash", mock.Anything, mock.Anything).Return(nil)

	resp, err := server.Logout(context.Background(), &pb.LogoutRequest{
		RefreshToken: "some-token",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	refreshRepo.AssertExpectations(t)
}

func TestLogout_gRPC_Error(t *testing.T) {
	server, _, refreshRepo := newTestServer()

	refreshRepo.On("RevokeByTokenHash", mock.Anything, mock.Anything).Return(errors.New("db error"))

	resp, err := server.Logout(context.Background(), &pb.LogoutRequest{
		RefreshToken: "some-token",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	refreshRepo.AssertExpectations(t)
}

func TestGetUserProfile_gRPC_Success(t *testing.T) {
	server, userRepo, _ := newTestServer()

	user := newActiveUser("user-1", "test@example.com")
	userRepo.On("GetByID", mock.Anything, "user-1").Return(user, nil)

	resp, err := server.GetUserProfile(context.Background(), &pb.GetUserProfileRequest{
		UserId: "user-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", resp.User.Id)
	assert.Equal(t, "test@example.com", resp.User.Email)
	userRepo.AssertExpectations(t)
}

func TestGetUserProfile_gRPC_NotFound(t *testing.T) {
	server, userRepo, _ := newTestServer()

	userRepo.On("GetByID", mock.Anything, "user-1").Return(nil, nil)

	_, err := server.GetUserProfile(context.Background(), &pb.GetUserProfileRequest{
		UserId: "user-1",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	userRepo.AssertExpectations(t)
}
