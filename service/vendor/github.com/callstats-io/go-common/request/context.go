package request

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/callstats-io/go-common/log"
)

// IDFromContext returns an requestID from the context if present or creates a new one
func IDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		return id
	}
	b := make([]byte, 16)
	// RFC 4122 compliant UUID (each byte is 2 chars hex)
	b[6] = (b[6] | 0x40) & 0x4F // 4 as 13th char
	b[8] = (b[8] | 0x80) & 0xBF // 8/9/a/b as 17th
	if _, err := rand.Read(b); err != nil {
		log.FromContextWithPackageName(ctx, "go-common/request").Error(LogErrFailedToCreateRequestID, log.Error(err))
		return ""
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// WithID returns a new context with the id stored as the request ID
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, id)
}
