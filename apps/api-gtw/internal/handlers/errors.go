package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errorResponse(c *gin.Context, httpStatus int, code string, message string) {
	c.JSON(httpStatus, gin.H{
		"error": apiError{
			Code:    code,
			Message: message,
		},
	})
}

func handleGRPCError(c *gin.Context, err error, fallbackStatus int, fallbackMessage string) {
	st, ok := status.FromError(err)
	if !ok {
		errorResponse(c, fallbackStatus, grpcCodeToErrorCode(codes.Internal), fallbackMessage)
		return
	}

	httpStatus, errorCode := mapGRPCCode(st.Code())
	if httpStatus != 0 {
		errorResponse(c, httpStatus, errorCode, st.Message())
		return
	}

	errorResponse(c, fallbackStatus, grpcCodeToErrorCode(st.Code()), fallbackMessage)
}

func mapGRPCCode(code codes.Code) (int, string) {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest, "INVALID_ARGUMENT"
	case codes.Unauthenticated:
		return http.StatusUnauthorized, "UNAUTHENTICATED"
	case codes.PermissionDenied:
		return http.StatusForbidden, "PERMISSION_DENIED"
	case codes.NotFound:
		return http.StatusNotFound, "NOT_FOUND"
	case codes.AlreadyExists:
		return http.StatusConflict, "ALREADY_EXISTS"
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests, "RESOURCE_EXHAUSTED"
	case codes.Unavailable:
		return http.StatusServiceUnavailable, "UNAVAILABLE"
	case codes.Internal:
		return http.StatusInternalServerError, "INTERNAL"
	default:
		return 0, ""
	}
}

func grpcCodeToErrorCode(code codes.Code) string {
	_, errorCode := mapGRPCCode(code)
	if errorCode != "" {
		return errorCode
	}
	return "INTERNAL"
}
