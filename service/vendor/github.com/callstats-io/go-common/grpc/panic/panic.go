package panic

import (
	"errors"
	"github.com/callstats-io/go-common/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"runtime/debug"
)

const (
	logKeyStacktrace = "panicStacktrace"
)

var errIntervalServerError = errors.New("Internal server error")

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (err error)

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor(recovery RecoveryHandlerFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recovery(r)
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func StreamServerInterceptor(recovery RecoveryHandlerFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recovery(r)
			}
		}()
		return handler(srv, stream)
	}
}

// LoggingRecovery provided recovery function for logging panic reason
func LoggingRecovery(ctx context.Context) RecoveryHandlerFunc {
	logger := log.FromContext(ctx)
	return func(interface{}) (err error) {
		logger.Error("Unhandled panic recovered", log.Object(logKeyStacktrace, debug.Stack()))
		return errIntervalServerError
	}
}
