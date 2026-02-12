package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buemura/minibank/api-gtw/internal/mocks"
	transactionpb "github.com/buemura/minibank/api-gtw/proto/transaction/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoTransaction() *transactionpb.Transaction {
	return &transactionpb.Transaction{
		Id:                       "tx-1",
		IdempotencyKey:           "key-1",
		Type:                     transactionpb.TransactionType_TRANSACTION_TYPE_TRANSFER,
		Status:                   transactionpb.TransactionStatus_TRANSACTION_STATUS_COMPLETED,
		SourceAccountId:          "src-1",
		DestinationAccountId:     "dest-1",
		DestinationAccountNumber: "9876543210",
		Amount:                   "100.00",
		Currency:                 "BRL",
		Description:              "test transfer",
		CreatedAt:                timestamppb.Now(),
		ProcessedAt:              timestamppb.Now(),
	}
}

func TestTransfer_Handler_Success(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Transfer", mock.Anything, mock.Anything).Return(&transactionpb.TransferResponse{
		Success:     true,
		Transaction: protoTransaction(),
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"idempotency_key":            "key-1",
		"destination_account_number": "9876543210",
		"amount":                     "100.00",
		"description":                "test",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/transfers", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Transfer(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
	txClient.AssertExpectations(t)
}

func TestTransfer_Handler_InvalidBody(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/transfers", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Transfer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	code, _ := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "INVALID_ARGUMENT", code)
}

func TestTransfer_Handler_gRPCError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Transfer", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))

	body, _ := json.Marshal(map[string]string{
		"idempotency_key":            "key-1",
		"destination_account_number": "9876543210",
		"amount":                     "100.00",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/transfers", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Transfer(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	txClient.AssertExpectations(t)
}

func TestTransfer_Handler_BusinessError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Transfer", mock.Anything, mock.Anything).Return(&transactionpb.TransferResponse{
		Success:      false,
		ErrorCode:    "INSUFFICIENT_FUNDS",
		ErrorMessage: "Not enough balance",
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"idempotency_key":            "key-1",
		"destination_account_number": "9876543210",
		"amount":                     "100.00",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/transfers", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Transfer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	code, msg := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "INSUFFICIENT_FUNDS", code)
	assert.Equal(t, "Not enough balance", msg)
	txClient.AssertExpectations(t)
}

func TestDeposit_Handler_Success(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	depositTx := protoTransaction()
	depositTx.Type = transactionpb.TransactionType_TRANSACTION_TYPE_DEPOSIT

	txClient.On("Deposit", mock.Anything, mock.Anything).Return(&transactionpb.DepositResponse{
		Success:     true,
		Transaction: depositTx,
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"idempotency_key": "key-1",
		"amount":          "500.00",
		"description":     "test deposit",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/deposits", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Deposit(c)

	assert.Equal(t, http.StatusOK, w.Code)
	txClient.AssertExpectations(t)
}

func TestDeposit_Handler_InvalidBody(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/deposits", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Deposit(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	code, _ := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "INVALID_ARGUMENT", code)
}

func TestDeposit_Handler_gRPCError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Deposit", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	body, _ := json.Marshal(map[string]string{
		"idempotency_key": "key-1",
		"amount":          "500.00",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/deposits", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Deposit(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	txClient.AssertExpectations(t)
}

func TestDeposit_Handler_BusinessError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Deposit", mock.Anything, mock.Anything).Return(&transactionpb.DepositResponse{
		Success:      false,
		ErrorCode:    "DEPOSIT_FAILED",
		ErrorMessage: "Failed to deposit",
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"idempotency_key": "key-1",
		"amount":          "500.00",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/deposits", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Deposit(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	code, msg := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "DEPOSIT_FAILED", code)
	assert.Equal(t, "Failed to deposit", msg)
	txClient.AssertExpectations(t)
}

func TestWithdraw_Handler_Success(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	withdrawTx := protoTransaction()
	withdrawTx.Type = transactionpb.TransactionType_TRANSACTION_TYPE_WITHDRAWAL

	txClient.On("Withdraw", mock.Anything, mock.Anything).Return(&transactionpb.WithdrawResponse{
		Success:     true,
		Transaction: withdrawTx,
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"idempotency_key": "key-1",
		"amount":          "200.00",
		"description":     "test withdrawal",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/withdrawals", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Withdraw(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
	txClient.AssertExpectations(t)
}

func TestWithdraw_Handler_InvalidBody(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/withdrawals", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Withdraw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	code, _ := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "INVALID_ARGUMENT", code)
}

func TestWithdraw_Handler_gRPCError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Withdraw", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	body, _ := json.Marshal(map[string]string{
		"idempotency_key": "key-1",
		"amount":          "200.00",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/withdrawals", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Withdraw(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	txClient.AssertExpectations(t)
}

func TestWithdraw_Handler_BusinessError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("Withdraw", mock.Anything, mock.Anything).Return(&transactionpb.WithdrawResponse{
		Success:      false,
		ErrorCode:    "WITHDRAWAL_FAILED",
		ErrorMessage: "insufficient funds",
	}, nil)

	body, _ := json.Marshal(map[string]string{
		"idempotency_key": "key-1",
		"amount":          "200.00",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/accounts/acc-1/withdrawals", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.Withdraw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	code, msg := parseErrorResponse(t, w.Body.Bytes())
	assert.Equal(t, "WITHDRAWAL_FAILED", code)
	assert.Equal(t, "insufficient funds", msg)
	txClient.AssertExpectations(t)
}

func TestGetTransaction_Handler_Success(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetTransaction", mock.Anything, mock.Anything).Return(&transactionpb.GetTransactionResponse{
		Transaction: protoTransaction(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/transactions/tx-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "tx-1"}}

	handler.GetTransaction(c)

	assert.Equal(t, http.StatusOK, w.Code)
	txClient.AssertExpectations(t)
}

func TestGetTransaction_Handler_NotFound(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetTransaction", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/transactions/tx-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "tx-1"}}

	handler.GetTransaction(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	txClient.AssertExpectations(t)
}

func TestGetTransactionHistory_Handler_Success(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetTransactionHistory", mock.Anything, mock.Anything).Return(&transactionpb.GetTransactionHistoryResponse{
		Transactions: []*transactionpb.Transaction{protoTransaction()},
		TotalCount:   1,
		Page:         1,
		PageSize:     20,
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/transactions?page=1&page_size=20", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetTransactionHistory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	transactions := resp["transactions"].([]interface{})
	assert.Len(t, transactions, 1)
	txClient.AssertExpectations(t)
}

func TestGetTransactionHistory_Handler_gRPCError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetTransactionHistory", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/transactions", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetTransactionHistory(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	txClient.AssertExpectations(t)
}

func TestGetStatement_Handler_Success(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetStatement", mock.Anything, mock.Anything).Return(&transactionpb.GetStatementResponse{
		AccountId:      "acc-1",
		OpeningBalance: "900.00",
		ClosingBalance: "1000.00",
		Transactions:   []*transactionpb.Transaction{protoTransaction()},
		TotalCount:     1,
		StartDate:      timestamppb.Now(),
		EndDate:        timestamppb.Now(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/statement?start_date=2025-01-01&end_date=2025-01-31", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetStatement(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "acc-1", resp["account_id"])
	assert.Equal(t, "900.00", resp["opening_balance"])
	assert.Equal(t, "1000.00", resp["closing_balance"])
	txClient.AssertExpectations(t)
}

func TestGetStatement_Handler_DefaultDates(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetStatement", mock.Anything, mock.Anything).Return(&transactionpb.GetStatementResponse{
		AccountId:      "acc-1",
		OpeningBalance: "1000.00",
		ClosingBalance: "1000.00",
		Transactions:   []*transactionpb.Transaction{},
		TotalCount:     0,
		StartDate:      timestamppb.Now(),
		EndDate:        timestamppb.Now(),
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/statement", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetStatement(c)

	assert.Equal(t, http.StatusOK, w.Code)
	txClient.AssertExpectations(t)
}

func TestGetStatement_Handler_gRPCError(t *testing.T) {
	txClient := new(mocks.MockTransactionServiceClient)
	handler := NewTransactionHandler(txClient)

	txClient.On("GetStatement", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/accounts/acc-1/statement?start_date=2025-01-01&end_date=2025-01-31", nil)
	c.Params = gin.Params{{Key: "id", Value: "acc-1"}}

	handler.GetStatement(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	txClient.AssertExpectations(t)
}
