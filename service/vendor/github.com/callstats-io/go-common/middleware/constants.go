package middleware

import "regexp"

// Header names
const (
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderOrigin                        = "Origin"
	HeaderEcho                          = "X-Csio-Echo"
)

// Header values
const (
	AccessControlAllowedCredentials = "true"
	AccessControlAllowedHeaders     = "Accept, Authorization, Content-Type, Accept-Encoding, X-Csio-Echo"
	AccessControlExposedHeaders     = "X-Csio-Echo"
	AccessControlMaxAge             = "900" // 15 minutes
)

// Log key names
const (
	LogKeyRemoteIP        = "remoteIP"
	LogKeyRequestID       = "requestID"
	LogKeyExpectedAppID   = "expectedAppID"
	LogKeyExpectedLocalID = "expectedLocalID"
	LogKeyRequestURL      = "url"
	LogKeyExpectedScope   = "expectedScope"
	LogKeyActualScope     = "actualScope"
	LogKeyExpectedOrigin  = "expectedOrigin"
	LogKeyActualOrigin    = "actualOrigin"
	LogKeyAuthToken       = "authToken"
)

// Log messages
const (
	LogMsgErrInvalidHTTPVersion    = "Received non HTTP/2 request on HTTP/2 required endpoint"
	LogMsgErrInvalidUcID           = "Invalid ucID"
	LogMsgErrInvalidConfID         = "Invalid confID"
	LogMsgErrInvalidAppID          = "Invalid appID"
	LogMsgErrInvalidAuthHeader     = "Invalid Authorization-header"
	LogMsgErrInvalidAuthToken      = "Invalid auth JWT token"
	LogMsgErrInvalidAuthScopeState = "Invalid auth scope parsing state"
	LogMsgErrInvalidAuthScope      = "Invalid auth scope"
	LogMsgErrInvalidJWTOrigin      = "Invalid JWT based origin url"
	LogMsgErrExpiredAuthToken      = "Expired auth JWT token"

	LogMsgStartRequest  = "Start request"
	LogMsgFinishRequest = "Finish request"
)

// Validators
var (
	ValidAppIDRegExp = regexp.MustCompile("\\A\\d+\\z")
	ValidUcIDRegExp  = regexp.MustCompile("\\A[kKpPnN]{0,1}\\d+\\z")
)

const (
	// MaxConfIDLength is the maximum confID length supported
	MaxConfIDLength = 256
)
