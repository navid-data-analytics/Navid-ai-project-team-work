package middleware

import "errors"

// Errors
var (
	ErrMissingAuthHeader       = errors.New("Authorization-header not in request")
	ErrMissingAuthBearerPrefix = errors.New("Missing \"Bearer \"-prefix in Authorization header content")
	ErrInvalidOriginURL        = errors.New("Invalid origin URL")
)
