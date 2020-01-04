package middleware

import (
	"net/http"
	"strings"

	"github.com/callstats-io/go-common/auth"
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/response"
)

// Header keys
const (
	HeaderAuthorization = "Authorization"
)

// Responses
var (
	RespEmptyAuthToken          = response.MarshalError(auth.ErrEmptyAuthToken)
	RespInvalidAuthToken        = response.MarshalError(auth.ErrInvalidAuthToken)
	RespMissingAuthHeader       = response.MarshalError(ErrMissingAuthHeader)
	RespMissingAuthBearerPrefix = response.MarshalError(ErrMissingAuthBearerPrefix)
	RespInvalidAuthAppID        = response.MarshalError(auth.ErrInvalidAuthAppID)
	RespInvalidAuthUserID       = response.MarshalError(auth.ErrInvalidAuthUserID)
	RespAuthExpiredToken        = response.MarshalError(auth.ErrAuthExpiredToken)
	RespInvalidScope            = response.MarshalError(auth.ErrInvalidScope)
	RespInvalidOriginURL        = response.MarshalError(ErrInvalidOriginURL)
)

const (
	authBearerPrefix = "Bearer "
)

// EndpointJWTAuth validates the JWT token and parses the endpoint claims to the request context.
// __NOTE__ This middleware does not verify the actual token content, such as user or appID, as getting that information is application specific.
// It expects the verification secret and algorithm to be passed as argument and returns a new middleware func.
// If the parsed token is valid and matches verification the middleware calls the next handler with context including the claims.
// Otherwise the middleware responds with HTTP 401 Unauthorized. Specifically this happens if:
// - The Authorization header could not be found
// - The Authorization header doesn't contain value of "Bearer <the jwt token>"
// - The JWT failed to parse correctly to the claims
// - The JWT token could not be verified with the secret
func EndpointJWTAuth(verifySecret []byte, signingMethod auth.SigningMethod) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get the auth token from header
			token, err := authTokenFromHeader(r)
			if err != nil {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware").Error(LogMsgErrInvalidAuthHeader, log.String(LogKeyRequestURL, r.RequestURI))
				respondUnauthorized(w, err)
				return
			}
			// parse and verify claims
			claims := &auth.EndpointClaims{}
			if err := claims.ParseAndVerify(r.Context(), verifySecret, signingMethod, token); err != nil {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware").Error(LogMsgErrInvalidAuthToken, log.String(LogKeyRequestURL, r.RequestURI))
				respondUnauthorized(w, err)
				return
			}

			// assign claims to context for further validation (e.g. appID/userID matches etc.)
			reqWithClaims := r.WithContext(auth.WithEndpointClaims(r.Context(), claims))

			// call next middleware/handler
			next(w, reqWithClaims)
		}
	}
}

// EndpointRequireJWTScopes checks that __ALL__ scopes in the required set are present in endpoint claims.
// It assumes endpoint claims have been stored in the context before calling this middleware.
// In case the claims are not found it returns HTTP 500 InternalServerError.
// In case the claims do not have all required scopes it responds with HTTP 401 Unauthorized.
func EndpointRequireJWTScopes(requiredScopes []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := auth.EndpointClaimsFromContext(r.Context())
			if claims == nil {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware").Error(LogMsgErrInvalidAuthScopeState)
				// if this is returned it means the middleware chain was not correctly set up to
				// parse and validate JWT token claims before calling this middleware.
				// See EndpointJWTAuthMiddleware above for how to do it.
				response.InternalServerError(w, nil)
				return
			}

			if err := claims.ValidateAuthScopes(requiredScopes); err != nil {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware").
					With(log.String(LogKeyRequestURL, r.RequestURI)).
					With(log.Object(LogKeyExpectedScope, requiredScopes)).
					With(log.Object(LogKeyActualScope, claims.Scope)).
					Error(LogMsgErrInvalidAuthScope, log.Error(err))
				respondUnauthorized(w, auth.ErrInvalidScope)
				return
			}

			next(w, r)
		}
	}
}

// EndpointRequireJWTScopesForHTTP1 allows request to be made with HTTP1.1 if requiredScopes is are in claim scopes
// otherwise it enforces http2 usage and returns http.StatusHTTPVersionNotSupported for http1 requests
// It assumes endpoint claims have been stored in the context before calling this middleware.
// In case the claims are not found it returns HTTP 500 InternalServerError.
// In case the claims do not have all required scopes it responds with HTTP 401 Unauthorized.
func EndpointRequireJWTScopesForHTTP1(requiredScopes []string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := auth.EndpointClaimsFromContext(r.Context())
			if claims == nil {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware").Error(LogMsgErrInvalidAuthScopeState)
				// if this is returned it means the middleware chain was not correctly set up to
				// parse and validate JWT token claims before calling this middleware.
				// See EndpointJWTAuthMiddleware above for how to do it.
				response.InternalServerError(w, nil)
				return
			}
			// require http2 if required scopes are not present in claims
			err := claims.ValidateAuthScopes(requiredScopes)
			if r.ProtoMajor != 2 && err != nil {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware").Error(LogMsgErrInvalidHTTPVersion, log.Error(err))
				response.InvalidHTTPVersion(w)
				return
			}
			h(w, r)
		}
	}
}

// EndpointJWTCORS sets the allowed origins based on the originURLs from the JWT token.
// In case there are no claims (i.e. it might not have been configured to be used at all)
// or the claims are empty/nil (might have not been added to the JWT token, go defaults to empty slice)
// any url is allwed and this echoes back the request origin header
func EndpointJWTCORS() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(HeaderOrigin)
			claims := auth.EndpointClaimsFromContext(r.Context())
			if claims != nil && origin != "" {
				if err := claims.ValidateOrigin(origin); err != nil {
					log.FromContextWithPackageName(r.Context(), "go-common/middleware").
						With(log.String(LogKeyRequestURL, r.RequestURI)).
						With(log.Object(LogKeyExpectedOrigin, claims.OriginURLs)).
						With(log.String(LogKeyActualOrigin, origin)).
						Error(LogMsgErrInvalidJWTOrigin)
					w.Header().Set(HeaderAccessControlAllowOrigin, "")
					response.PreconditionFailed(w, RespInvalidOriginURL)
					return
				}
			}
			w.Header().Set(HeaderAccessControlAllowOrigin, origin)
			next(w, r)
		}
	}
}

func authTokenFromHeader(r *http.Request) (string, error) {
	headerValue := r.Header.Get(HeaderAuthorization)

	if headerValue == "" {
		return "", ErrMissingAuthHeader
	}

	if !strings.HasPrefix(headerValue, authBearerPrefix) {
		return "", ErrMissingAuthBearerPrefix
	}

	return headerValue[len(authBearerPrefix):], nil
}

func respondUnauthorized(w http.ResponseWriter, err error) {
	switch err {
	case auth.ErrEmptyAuthToken:
		response.Unauthorized(w, RespEmptyAuthToken)
	case auth.ErrInvalidAuthToken:
		response.Unauthorized(w, RespInvalidAuthToken)
	case ErrMissingAuthHeader:
		response.Unauthorized(w, RespMissingAuthHeader)
	case ErrMissingAuthBearerPrefix:
		response.Unauthorized(w, RespMissingAuthBearerPrefix)
	case auth.ErrAuthExpiredToken:
		response.Unauthorized(w, RespAuthExpiredToken)
	case auth.ErrInvalidScope:
		response.Unauthorized(w, RespInvalidScope)
	default:
		panic("Unknown auth error") // should only happen in development
	}
}
