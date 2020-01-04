package middleware

import (
	"net/http"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/request"
)

// RequestContextLogger adds an request id based logger to the context and logs start/end of request.
// It should be set as one of the first middlewares to get maximum context
func RequestContextLogger() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			requestID := request.IDFromContext(ctx)
			ctx = request.WithID(ctx, requestID)
			logger := log.FromContextWithPackageName(ctx, "go-common/middleware/request").With(log.String(LogKeyRequestID, requestID))
			ctx = log.WithLogger(ctx, logger)
			h(w, r.WithContext(ctx))
		}
	}
}

// RequestLifetimeLogger adds an request id based logger to the context and logs start/end of request.
// It should be set as one of the first middlewares to get maximum context
func RequestLifetimeLogger() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.FromContextWithPackageName(r.Context(), "go-common/middleware/request").
				With(log.String(LogKeyRemoteIP, request.IPFromAddr(r.RemoteAddr))).
				With(log.String(LogKeyRequestURL, r.RequestURI))
			logger.Info(LogMsgStartRequest)
			h(w, r)
			logger.Info(LogMsgFinishRequest)
		}
	}
}
