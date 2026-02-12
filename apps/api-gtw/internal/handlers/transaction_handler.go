package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"

	transactionpb "github.com/buemura/minibank/api-gtw/proto/transaction/v1"
)

type TransactionHandler struct {
	transactionClient transactionpb.TransactionServiceClient
}

func NewTransactionHandler(transactionClient transactionpb.TransactionServiceClient) *TransactionHandler {
	return &TransactionHandler{transactionClient: transactionClient}
}

type TransferRequest struct {
	IdempotencyKey           string `json:"idempotency_key" binding:"required"`
	DestinationAccountNumber string `json:"destination_account_number" binding:"required"`
	Amount                   string `json:"amount" binding:"required"`
	Description              string `json:"description"`
}

type DepositRequest struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	Amount         string `json:"amount" binding:"required"`
	Description    string `json:"description"`
}

type WithdrawRequest struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	Amount         string `json:"amount" binding:"required"`
	Description    string `json:"description"`
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	accountID := c.Param("id")

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.transactionClient.Transfer(c.Request.Context(), &transactionpb.TransferRequest{
		IdempotencyKey:           req.IdempotencyKey,
		SourceAccountId:          accountID,
		DestinationAccountNumber: req.DestinationAccountNumber,
		Amount:                   req.Amount,
		Description:              req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error_code":    resp.ErrorCode,
			"error_message": resp.ErrorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"transaction": mapProtoTransactionToResponse(resp.Transaction),
	})
}

func (h *TransactionHandler) Deposit(c *gin.Context) {
	accountID := c.Param("id")

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.transactionClient.Deposit(c.Request.Context(), &transactionpb.DepositRequest{
		IdempotencyKey: req.IdempotencyKey,
		AccountId:      accountID,
		Amount:         req.Amount,
		Description:    req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deposit failed"})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error_code":    resp.ErrorCode,
			"error_message": resp.ErrorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"transaction": mapProtoTransactionToResponse(resp.Transaction),
	})
}

func (h *TransactionHandler) Withdraw(c *gin.Context) {
	accountID := c.Param("id")

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.transactionClient.Withdraw(c.Request.Context(), &transactionpb.WithdrawRequest{
		IdempotencyKey: req.IdempotencyKey,
		AccountId:      accountID,
		Amount:         req.Amount,
		Description:    req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Withdrawal failed"})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":       false,
			"error_code":    resp.ErrorCode,
			"error_message": resp.ErrorMessage,
		})
		return
	}

	response := gin.H{
		"success":     true,
		"transaction": mapProtoTransactionToResponse(resp.Transaction),
	}
	if resp.FeeTransaction != nil {
		response["fee_transaction"] = mapProtoTransactionToResponse(resp.FeeTransaction)
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	resp, err := h.transactionClient.GetTransaction(c.Request.Context(), &transactionpb.GetTransactionRequest{
		TransactionId: transactionID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": mapProtoTransactionToResponse(resp.Transaction),
	})
}

func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	accountID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := h.transactionClient.GetTransactionHistory(c.Request.Context(), &transactionpb.GetTransactionHistoryRequest{
		AccountId: accountID,
		Page:      int32(page),
		PageSize:  int32(pageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction history"})
		return
	}

	transactions := make([]map[string]interface{}, 0)
	for _, tx := range resp.Transactions {
		transactions = append(transactions, mapProtoTransactionToResponse(tx))
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total_count":  resp.TotalCount,
		"page":         resp.Page,
		"page_size":    resp.PageSize,
	})
}

func (h *TransactionHandler) GetStatement(c *gin.Context) {
	accountID := c.Param("id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	resp, err := h.transactionClient.GetStatement(c.Request.Context(), &transactionpb.GetStatementRequest{
		AccountId: accountID,
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
		Page:      int32(page),
		PageSize:  int32(pageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statement"})
		return
	}

	transactions := make([]map[string]interface{}, 0)
	for _, tx := range resp.Transactions {
		transactions = append(transactions, mapProtoTransactionToResponse(tx))
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":      resp.AccountId,
		"opening_balance": resp.OpeningBalance,
		"closing_balance": resp.ClosingBalance,
		"transactions":    transactions,
		"total_count":     resp.TotalCount,
		"start_date":      resp.StartDate.AsTime().Format("2006-01-02"),
		"end_date":        resp.EndDate.AsTime().Format("2006-01-02"),
	})
}

func mapProtoTransactionToResponse(tx *transactionpb.Transaction) map[string]interface{} {
	result := map[string]interface{}{
		"id":                         tx.Id,
		"idempotency_key":            tx.IdempotencyKey,
		"type":                       mapTransactionType(tx.Type),
		"status":                     mapTransactionStatus(tx.Status),
		"source_account_id":          tx.SourceAccountId,
		"source_account_number":      tx.SourceAccountNumber,
		"destination_account_id":     tx.DestinationAccountId,
		"destination_account_number": tx.DestinationAccountNumber,
		"amount":                     tx.Amount,
		"currency":                   tx.Currency,
		"description":                tx.Description,
		"created_at":                 tx.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z07:00"),
	}

	if tx.ProcessedAt != nil {
		result["processed_at"] = tx.ProcessedAt.AsTime().Format("2006-01-02T15:04:05Z07:00")
	}

	return result
}

func mapTransactionType(t transactionpb.TransactionType) string {
	switch t {
	case transactionpb.TransactionType_TRANSACTION_TYPE_TRANSFER:
		return "TRANSFER"
	case transactionpb.TransactionType_TRANSACTION_TYPE_DEPOSIT:
		return "DEPOSIT"
	case transactionpb.TransactionType_TRANSACTION_TYPE_WITHDRAWAL:
		return "WITHDRAWAL"
	case transactionpb.TransactionType_TRANSACTION_TYPE_FEE:
		return "FEE"
	default:
		return "UNKNOWN"
	}
}

func mapTransactionStatus(s transactionpb.TransactionStatus) string {
	switch s {
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_PENDING:
		return "PENDING"
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_COMPLETED:
		return "COMPLETED"
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_FAILED:
		return "FAILED"
	case transactionpb.TransactionStatus_TRANSACTION_STATUS_REVERSED:
		return "REVERSED"
	default:
		return "UNKNOWN"
	}
}
