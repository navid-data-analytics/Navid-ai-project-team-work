package grpc

import (
	"github.com/callstats-io/go-common/log"
	context "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrUnavailable logs and wraps the given error with gRPC error code Unavailable
func ErrUnavailable(ctx context.Context, err error) error {
	log.FromContext(ctx).Error("unavailable", log.Error(err))
	return status.Error(codes.Unavailable, err.Error())
}

// ErrInvalidArgument logs and wraps the given error with gRPC error code InvalidArgument
func ErrInvalidArgument(ctx context.Context, err error) error {
	log.FromContext(ctx).Error("invalid argument", log.Error(err))
	return status.Error(codes.InvalidArgument, err.Error())
}

// ErrNotFound logs and wraps the given error with gRPC error code NotFound
func ErrNotFound(ctx context.Context, err error) error {
	log.FromContext(ctx).Error("not found", log.Error(err))
	return status.Error(codes.NotFound, err.Error())
}

// ErrFailedPrecondition logs and wraps the given error with gRPC error code FailedPrecondition
func ErrFailedPrecondition(ctx context.Context, err error) error {
	log.FromContext(ctx).Error("failed precondition", log.Error(err))
	return status.Error(codes.FailedPrecondition, err.Error())
}
