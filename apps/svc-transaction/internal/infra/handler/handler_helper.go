package handler

import (
	"errors"

	"github.com/buemura/minibank/svc-transaction/internal/core/domain/transaction"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrBadRequest = errors.New("bad request")
var ErrInvalidArgument = errors.New("invalid argument")

func HandleGrpcError(err error) error {
	switch {
	case errors.Is(err, transaction.ErrTransactionNotFound): // 5
		return status.Error(codes.NotFound, err.Error())

	default: // pass code or 13
		if serr, ok := status.FromError(err); ok {
			return status.Error(serr.Code(), err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}

}
