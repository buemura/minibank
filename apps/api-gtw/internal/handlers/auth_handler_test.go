package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buemura/minibank/api-gtw/internal/mocks"
	authpb "github.com/buemura/minibank/api-gtw/proto/auth/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func protoUser() *authpb.User {
	return &authpb.User{
		Id:        "user-1",
		Email:     "test@example.com",
		FullName:  "Test User",
		Phone:     "+5511999999999",
		Status:    "active",
		CreatedAt: timestamppb.Now(),
	}
}

func TestRegister_Success(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("Register", mock.Anything, mock.Anything).Return(&authpb.RegisterResponse{
		User:         protoUser(),
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    900,
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"email":     "test@example.com",
		"password":  "password123",
		"full_name": "Test User",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "access-token", resp["access_token"])
	assert.Equal(t, "refresh-token", resp["refresh_token"])
	authClient.AssertExpectations(t)
}

func TestRegister_InvalidBody(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	body, _ := json.Marshal(map[string]string{
		"email": "invalid-email",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_gRPCError(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("Register", mock.Anything, mock.Anything).Return(nil, errors.New("user already exists"))

	body, _ := json.Marshal(map[string]string{
		"email":     "test@example.com",
		"password":  "password123",
		"full_name": "Test User",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	authClient.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("Login", mock.Anything, mock.Anything).Return(&authpb.LoginResponse{
		User:         protoUser(),
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    900,
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	authClient.AssertExpectations(t)
}

func TestLogin_InvalidBody(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Failure(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("invalid credentials"))

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "wrong-password",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	authClient.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("RefreshToken", mock.Anything, mock.Anything).Return(&authpb.RefreshTokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresIn:    900,
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "old-refresh-token",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Refresh(c)

	assert.Equal(t, http.StatusOK, w.Code)
	authClient.AssertExpectations(t)
}

func TestRefresh_InvalidBody(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Refresh(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefresh_Failure(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("RefreshToken", mock.Anything, mock.Anything).Return(nil, errors.New("invalid token"))

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "invalid-token",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Refresh(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	authClient.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("Logout", mock.Anything, mock.Anything).Return(&authpb.LogoutResponse{Success: true}, nil)

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "token-to-revoke",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/logout", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	authClient.AssertExpectations(t)
}

func TestLogout_Failure(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("Logout", mock.Anything, mock.Anything).Return(nil, errors.New("failed"))

	body, _ := json.Marshal(map[string]string{
		"refresh_token": "token-to-revoke",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/logout", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Logout(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	authClient.AssertExpectations(t)
}

func TestMe_Success(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("GetUserProfile", mock.Anything, mock.Anything).Return(&authpb.GetUserProfileResponse{
		User: protoUser(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/me", nil)
	c.Set("user_id", "user-1")

	handler.Me(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user-1", resp["id"])
	assert.Equal(t, "test@example.com", resp["email"])
	authClient.AssertExpectations(t)
}

func TestMe_Failure(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("GetUserProfile", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/me", nil)
	c.Set("user_id", "user-1")

	handler.Me(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	authClient.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("ChangePassword", mock.Anything, mock.Anything).Return(&authpb.ChangePasswordResponse{
		Success: true,
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"current_password": "oldpassword123",
		"new_password":     "newpassword456",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/change-password", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-1")

	handler.ChangePassword(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
	authClient.AssertExpectations(t)
}

func TestChangePassword_InvalidBody(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	body, _ := json.Marshal(map[string]string{
		"current_password": "old",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/change-password", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-1")

	handler.ChangePassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangePassword_gRPCError(t *testing.T) {
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAuthHandler(authClient)

	authClient.On("ChangePassword", mock.Anything, mock.Anything).Return(nil, status.Error(codes.InvalidArgument, "current password is incorrect"))

	body, _ := json.Marshal(map[string]string{
		"current_password": "wrongpassword1",
		"new_password":     "newpassword456",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/change-password", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-1")

	handler.ChangePassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	authClient.AssertExpectations(t)
}
