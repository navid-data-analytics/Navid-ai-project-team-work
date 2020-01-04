package auth

import "context"

type ctxKey int

// Context keys
const (
	ctxClaims ctxKey = iota
)

// EndpointClaimsFromContext returns the claims stored in context or nil if none have been stored.
func EndpointClaimsFromContext(ctx context.Context) *EndpointClaims {
	if claims, ok := ctx.Value(ctxClaims).(*EndpointClaims); ok {
		return claims
	}
	return nil
}

// WithEndpointClaims returns a new context with the passed in context as parent that contains the claims. The claims can be retrieved using EndpointClaimsFromContext.
func WithEndpointClaims(ctx context.Context, claims *EndpointClaims) context.Context {
	return context.WithValue(ctx, ctxClaims, claims)
}
