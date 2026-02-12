package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	errorResponse(c, http.StatusBadRequest, "INVALID_ARGUMENT", "invalid email format")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "INVALID_ARGUMENT", body["error"]["code"])
	assert.Equal(t, "invalid email format", body["error"]["message"])
}

func TestHandleGRPCError_MappedCode(t *testing.T) {
	tests := []struct {
		name           string
		grpcCode       codes.Code
		grpcMsg        string
		expectedStatus int
		expectedCode   string
	}{
		{"NotFound", codes.NotFound, "account not found", http.StatusNotFound, "NOT_FOUND"},
		{"InvalidArgument", codes.InvalidArgument, "invalid amount", http.StatusBadRequest, "INVALID_ARGUMENT"},
		{"Unauthenticated", codes.Unauthenticated, "invalid credentials", http.StatusUnauthorized, "UNAUTHENTICATED"},
		{"AlreadyExists", codes.AlreadyExists, "user already exists", http.StatusConflict, "ALREADY_EXISTS"},
		{"Internal", codes.Internal, "internal error", http.StatusInternalServerError, "INTERNAL"},
		{"Unavailable", codes.Unavailable, "service unavailable", http.StatusServiceUnavailable, "UNAVAILABLE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			grpcErr := status.Error(tt.grpcCode, tt.grpcMsg)
			handleGRPCError(c, grpcErr, http.StatusInternalServerError, "fallback")

			assert.Equal(t, tt.expectedStatus, w.Code)

			var body map[string]map[string]string
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
			assert.Equal(t, tt.expectedCode, body["error"]["code"])
			assert.Equal(t, tt.grpcMsg, body["error"]["message"])
		})
	}
}

func TestHandleGRPCError_UnmappedCode(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	grpcErr := status.Error(codes.DataLoss, "data loss occurred")
	handleGRPCError(c, grpcErr, http.StatusBadGateway, "Something went wrong")

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var body map[string]map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Something went wrong", body["error"]["message"])
}

func TestHandleGRPCError_NonGRPCError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handleGRPCError(c, errors.New("plain error"), http.StatusInternalServerError, "Operation failed")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "INTERNAL", body["error"]["code"])
	assert.Equal(t, "Operation failed", body["error"]["message"])
}
