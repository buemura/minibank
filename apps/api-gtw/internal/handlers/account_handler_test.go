package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buemura/minibank/api-gtw/internal/mocks"
	accountpb "github.com/buemura/minibank/api-gtw/proto/account/v1"
	authpb "github.com/buemura/minibank/api-gtw/proto/auth/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoAccount() *accountpb.Account {
	return &accountpb.Account{
		Id:            "acc-1",
		UserId:        "user-1",
		AccountNumber: "1234567890",
		Agency:        "0001",
		AccountType:   "checking",
		Balance:       "1000.00",
		Currency:      "BRL",
		Status:        "active",
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
	}
}

func TestCreateAccount_Success(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("CreateAccount", mock.Anything, mock.Anything).Return(&accountpb.CreateAccountResponse{
		Account: protoAccount(),
	}, nil)

	body, _ := json.Marshal(map[string]string{"account_type": "checking"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-1")

	handler.CreateAccount(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	account := resp["account"].(map[string]interface{})
	assert.Equal(t, "acc-1", account["id"])
	accountClient.AssertExpectations(t)
}

func TestCreateAccount_DefaultType(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("CreateAccount", mock.Anything, mock.MatchedBy(func(req *accountpb.CreateAccountRequest) bool {
		return req.AccountType == "checking"
	})).Return(&accountpb.CreateAccountResponse{
		Account: protoAccount(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-1")

	handler.CreateAccount(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	accountClient.AssertExpectations(t)
}

func TestCreateAccount_gRPCError(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("CreateAccount", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", "user-1")

	handler.CreateAccount(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	accountClient.AssertExpectations(t)
}

func TestGetAccounts_Success(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccountByUserId", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByUserIdResponse{
		Accounts: []*accountpb.Account{protoAccount()},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts", nil)
	c.Set("user_id", "user-1")

	handler.GetAccounts(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	accounts := resp["accounts"].([]interface{})
	assert.Len(t, accounts, 1)
	accountClient.AssertExpectations(t)
}

func TestGetAccounts_gRPCError(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccountByUserId", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts", nil)
	c.Set("user_id", "user-1")

	handler.GetAccounts(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	accountClient.AssertExpectations(t)
}

func TestGetAccount_Success(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccount", mock.Anything, mock.Anything).Return(&accountpb.GetAccountResponse{
		Account: protoAccount(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetAccount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	accountClient.AssertExpectations(t)
}

func TestGetAccount_NotFound(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccount", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetAccount(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	accountClient.AssertExpectations(t)
}

func TestGetBalance_Success(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetBalance", mock.Anything, mock.Anything).Return(&accountpb.GetBalanceResponse{
		Balance:  "1000.00",
		Currency: "BRL",
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/balance", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetBalance(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "1000.00", resp["balance"])
	assert.Equal(t, "BRL", resp["currency"])
	accountClient.AssertExpectations(t)
}

func TestGetBalance_NotFound(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetBalance", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/balance", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetBalance(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	accountClient.AssertExpectations(t)
}

func TestLookupAccountByNumber_Success(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByNumberResponse{
		Account: protoAccount(),
	}, nil)
	authClient.On("GetUserProfile", mock.Anything, mock.Anything).Return(&authpb.GetUserProfileResponse{
		User: protoUser(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/lookup?account_number=1234567890", nil)

	handler.LookupAccountByNumber(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "1234567890", resp["account_number"])
	assert.Equal(t, "0001", resp["agency"])
	assert.Equal(t, "Test User", resp["owner_name"])
	accountClient.AssertExpectations(t)
	authClient.AssertExpectations(t)
}

func TestLookupAccountByNumber_MissingParam(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/lookup", nil)

	handler.LookupAccountByNumber(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLookupAccountByNumber_NotFound(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/lookup?account_number=0000000000", nil)

	handler.LookupAccountByNumber(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	accountClient.AssertExpectations(t)
}

func TestLookupAccountByNumber_UserFetchFails(t *testing.T) {
	accountClient := new(mocks.MockAccountServiceClient)
	authClient := new(mocks.MockAuthServiceClient)
	handler := NewAccountHandler(accountClient, authClient)

	accountClient.On("GetAccountByNumber", mock.Anything, mock.Anything).Return(&accountpb.GetAccountByNumberResponse{
		Account: protoAccount(),
	}, nil)
	authClient.On("GetUserProfile", mock.Anything, mock.Anything).Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/lookup?account_number=1234567890", nil)

	handler.LookupAccountByNumber(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	accountClient.AssertExpectations(t)
	authClient.AssertExpectations(t)
}
