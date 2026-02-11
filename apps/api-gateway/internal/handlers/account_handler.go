package handlers

import (
	"net/http"

	accountpb "github.com/buemura/minibank/api-gateway/proto/account/v1"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountClient accountpb.AccountServiceClient
}

func NewAccountHandler(accountClient accountpb.AccountServiceClient) *AccountHandler {
	return &AccountHandler{accountClient: accountClient}
}

type CreateAccountRequest struct {
	AccountType string `json:"account_type"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.AccountType = "checking"
	}

	if req.AccountType == "" {
		req.AccountType = "checking"
	}

	resp, err := h.accountClient.CreateAccount(c.Request.Context(), &accountpb.CreateAccountRequest{
		UserId:      userID,
		AccountType: req.AccountType,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"account": mapProtoAccountToResponse(resp.Account),
	})
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	userID := c.GetString("user_id")

	resp, err := h.accountClient.GetAccountByUserId(c.Request.Context(), &accountpb.GetAccountByUserIdRequest{
		UserId: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get accounts"})
		return
	}

	accounts := make([]map[string]interface{}, 0)
	for _, acc := range resp.Accounts {
		accounts = append(accounts, mapProtoAccountToResponse(acc))
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountID := c.Param("id")

	resp, err := h.accountClient.GetAccount(c.Request.Context(), &accountpb.GetAccountRequest{
		AccountId: accountID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account": mapProtoAccountToResponse(resp.Account),
	})
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	accountID := c.Param("id")

	resp, err := h.accountClient.GetBalance(c.Request.Context(), &accountpb.GetBalanceRequest{
		AccountId: accountID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":  resp.Balance,
		"currency": resp.Currency,
	})
}

func mapProtoAccountToResponse(account *accountpb.Account) map[string]interface{} {
	return map[string]interface{}{
		"id":             account.Id,
		"user_id":        account.UserId,
		"account_number": account.AccountNumber,
		"agency":         account.Agency,
		"account_type":   account.AccountType,
		"balance":        account.Balance,
		"currency":       account.Currency,
		"status":         account.Status,
		"created_at":     account.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":     account.UpdatedAt.AsTime().Format("2006-01-02T15:04:05Z07:00"),
	}
}
