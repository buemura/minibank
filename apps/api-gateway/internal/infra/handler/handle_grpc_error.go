package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcToHttpError(c echo.Context, err error) error {
	if serr, ok := status.FromError(err); ok {
		switch serr.Code() {
		case codes.InvalidArgument:
			return c.NoContent(http.StatusUnprocessableEntity)

		case codes.NotFound:
			return c.NoContent(http.StatusNotFound)

		case codes.AlreadyExists:
			return c.NoContent(http.StatusUnprocessableEntity)

		case codes.PermissionDenied:
			return c.NoContent(http.StatusUnauthorized)

		case codes.FailedPrecondition:
			return c.NoContent(http.StatusUnprocessableEntity)

		case codes.Unavailable:
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": err.Error(),
			})

		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return nil
}
