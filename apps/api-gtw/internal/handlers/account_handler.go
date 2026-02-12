package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/buemura/minibank/api-gtw/internal/cache"
	accountpb "github.com/buemura/minibank/api-gtw/proto/account/v1"
	authpb "github.com/buemura/minibank/api-gtw/proto/auth/v1"
	"github.com/gin-gonic/gin"
)

const (
	accountCacheTTL = 30 * time.Second
	lookupCacheTTL  = 2 * time.Minute
)

type AccountHandler struct {
	accountClient accountpb.AccountServiceClient
	authClient    authpb.AuthServiceClient
	cache         cache.Cache
}

func NewAccountHandler(accountClient accountpb.AccountServiceClient, authClient authpb.AuthServiceClient, cache cache.Cache) *AccountHandler {
	return &AccountHandler{accountClient: accountClient, authClient: authClient, cache: cache}
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
		handleGRPCError(c, err, http.StatusInternalServerError, "Failed to create account")
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
		handleGRPCError(c, err, http.StatusInternalServerError, "Failed to get accounts")
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
	cacheKey := "account:id:" + accountID

	// Try cache first
	if cached, err := h.cache.Get(cacheKey); err == nil && cached != "" {
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	resp, err := h.accountClient.GetAccount(c.Request.Context(), &accountpb.GetAccountRequest{
		AccountId: accountID,
	})

	if err != nil {
		handleGRPCError(c, err, http.StatusNotFound, "Account not found")
		return
	}

	result := gin.H{"account": mapProtoAccountToResponse(resp.Account)}
	jsonBytes, _ := json.Marshal(result)
	_ = h.cache.Set(cacheKey, string(jsonBytes), accountCacheTTL)

	c.JSON(http.StatusOK, result)
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	accountID := c.Param("id")

	resp, err := h.accountClient.GetBalance(c.Request.Context(), &accountpb.GetBalanceRequest{
		AccountId: accountID,
	})

	if err != nil {
		handleGRPCError(c, err, http.StatusNotFound, "Account not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":  resp.Balance,
		"currency": resp.Currency,
	})
}

func (h *AccountHandler) LookupAccountByNumber(c *gin.Context) {
	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		errorResponse(c, http.StatusBadRequest, "INVALID_ARGUMENT", "account_number query parameter is required")
		return
	}

	cacheKey := "account:number:" + accountNumber

	// Try cache first
	if cached, err := h.cache.Get(cacheKey); err == nil && cached != "" {
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	accountResp, err := h.accountClient.GetAccountByNumber(c.Request.Context(), &accountpb.GetAccountByNumberRequest{
		AccountNumber: accountNumber,
	})
	if err != nil {
		handleGRPCError(c, err, http.StatusNotFound, "Account not found")
		return
	}

	userResp, err := h.authClient.GetUserProfile(c.Request.Context(), &authpb.GetUserProfileRequest{
		UserId: accountResp.Account.UserId,
	})
	if err != nil {
		handleGRPCError(c, err, http.StatusInternalServerError, "Failed to fetch account owner")
		return
	}

	result := gin.H{
		"account_number": accountResp.Account.AccountNumber,
		"agency":         accountResp.Account.Agency,
		"owner_name":     userResp.User.FullName,
	}
	jsonBytes, _ := json.Marshal(result)
	_ = h.cache.Set(cacheKey, string(jsonBytes), lookupCacheTTL)

	c.JSON(http.StatusOK, result)
}

func mapProtoAccountToResponse(account *accountpb.Account) map[string]any {
	return map[string]any{
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
