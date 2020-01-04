package auth

import jwt "github.com/dgrijalva/jwt-go"

// SigningMethod exports jwt signing method as auth.SigningMethod
type SigningMethod jwt.SigningMethod

// Signing methods
var (
	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS384 = jwt.SigningMethodHS384
	SigningMethodHS512 = jwt.SigningMethodHS512
	SigningMethodES256 = jwt.SigningMethodES256
	SigningMethodES384 = jwt.SigningMethodES384
	SigningMethodES512 = jwt.SigningMethodES512
)

// Log field keys
const (
	LogKeyAuthToken = "authToken"
)

// Log messages
const (
	LogMsgErrInvalidAuthToken = "Invalid auth JWT token"
	LogMsgErrExpiredAuthToken = "Expired auth JWT token"
)
