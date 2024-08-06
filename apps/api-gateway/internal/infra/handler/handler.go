package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/buemura/minibank/api-gateway/internal/core/domain/statement"
	"github.com/buemura/minibank/api-gateway/internal/core/domain/transaction"
	"github.com/buemura/minibank/api-gateway/internal/core/usecase"
	"github.com/buemura/minibank/api-gateway/internal/infra/grpc"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, nil)
	})

	e.GET("/account/:id/statement", getStatement)
	e.POST("/account/:id/transaction", createTransaction)
}

func getStatement(c echo.Context) error {
	accID := c.Param("id")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	items, _ := strconv.Atoi(c.QueryParam("items"))

	accService := grpc.NewGrpcAccountService()
	trxService := grpc.NewGrpcTransactionService()

	uc := usecase.NewGetAccountStatement(accService, trxService)
	res, err := uc.Execute(&statement.GetStatementIn{
		AccountID: accID,
		Page:      page,
		Items:     items,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

func createTransaction(c echo.Context) error {
	body := new(transaction.CreateTransactionIn)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err := validateCreateTransactionPayload(body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	trxService := grpc.NewGrpcTransactionService()
	uc := usecase.NewCreateTransaction(trxService)
	trx, err := uc.Execute(body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, trx)
}

func validateCreateTransactionPayload(in *transaction.CreateTransactionIn) error {
	if in.AccountID == "" || in.TransactionType == "" || in.Amount == 0 {
		return errors.New("invalid arguments")
	}
	return nil
}
