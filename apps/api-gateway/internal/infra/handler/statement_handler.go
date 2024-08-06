package handler

import (
	"net/http"
	"strconv"

	"github.com/buemura/minibank/api-gateway/internal/core/domain/statement"
	"github.com/buemura/minibank/api-gateway/internal/core/usecase"
	"github.com/buemura/minibank/api-gateway/internal/infra/grpc"
	"github.com/labstack/echo/v4"
)

func SetupStatementRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, nil)
	})

	e.GET("/account/:id/statement", getEndpoint)
}

func getEndpoint(c echo.Context) error {
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
