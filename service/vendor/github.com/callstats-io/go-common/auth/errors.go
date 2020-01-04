package auth

import "errors"

// Errors
var (
	ErrEmptyAuthToken       = errors.New("Empty auth JWT token")
	ErrInvalidAuthToken     = errors.New("Invalid auth JWT token")
	ErrInvalidAuthAppID     = errors.New("Token appID does not match appID in request URL")
	ErrInvalidAuthUserID    = errors.New("Token userID does not match localID in request body")
	ErrInvalidScope         = errors.New("Required scope(s) missing from auth token")
	ErrAuthExpiredToken     = errors.New("Token has expired")
	ErrInvalidSigningMethod = errors.New("Invalid signing method. Should be either ES256 or HS256")
	ErrParseOriginURL       = errors.New("Failed to parse wildcard origin")
	ErrInvalidOriginURL     = errors.New("Origin does not match any of the allowed origins")
)
