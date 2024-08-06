package hanlder

import (
	"errors"
	"net/http"

	"github.com/buemura/minibank/api-gateway/internal/core/usecase"
	"github.com/buemura/minibank/api-gateway/internal/infra/grpc"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, nil)
	})

	e.GET("/endpoints", getEndpoint)
	e.POST("/endpoints", createEndpoint)
}

func getEndpoint(c echo.Context) error {
	accService := grpc.NewGrpcAccountService()
	trxService := grpc.NewGrpcTransactionService()
	uc := usecase.NewGetAccountStatement(accService, trxService)
	res, err := uc.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

func createEndpoint(c echo.Context) error {
	body := new(dto.CreateEndpointIn)
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}
	if err := validateCreateEndpoint(body); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	er := database.NewEndpointRepositoryImpl(database.DB)
	uc := usecase.NewCreateEndpoint(er)
	res, err := uc.Execute(body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

func validateCreateEndpoint(body *dto.CreateEndpointIn) error {
	if len(body.Name) < 1 {
		return errors.New("invalid Name")
	}
	if len(body.Url) < 1 {
		return errors.New("invalid Url")
	}
	if len(body.NotifyTo) < 1 {
		return errors.New("invalid NotifyTo")
	}
	if body.CheckFrequency < 1 {
		return errors.New("invalid CheckFrequency")
	}
	return nil
}
