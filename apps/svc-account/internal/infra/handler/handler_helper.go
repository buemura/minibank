package handler

import (
	"errors"

	"github.com/buemura/minibank/svc-account/internal/core/domain/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrBadRequest = errors.New("bad request")
var ErrInvalidArgument = errors.New("invalid argument")

func HandleGrpcError(err error) error {
	switch {
	case errors.Is(err, account.ErrAccountNotFound): // 5
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, account.ErrAccountAlreadyExists): // 5
		return status.Error(codes.AlreadyExists, err.Error())

	default: // pass code or 13
		if serr, ok := status.FromError(err); ok {
			return status.Error(serr.Code(), err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

}
