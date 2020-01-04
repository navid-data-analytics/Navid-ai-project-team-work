package middleware

import "net/http"

// CorsOptions specifies the options CORS Request handler expects
type CorsOptions struct {
	AllowedOrigin  func(r *http.Request) string
	AllowedMethods func(r *http.Request) string
	ExposedHeaders func(r *http.Request) string
}

// DefaultCorsOptions returns options which:
// 1) Return r.Header.Get("Origin") for allowed origin
// 2) Return POST for allowed methods
// 3) Return "X-Csio-Echo" for exposed headers
func DefaultCorsOptions() *CorsOptions {
	return &CorsOptions{
		AllowedOrigin: func(r *http.Request) string {
			return r.Header.Get(HeaderOrigin)
		},
		AllowedMethods: func(r *http.Request) string {
			return http.MethodPost
		},
		ExposedHeaders: func(r *http.Request) string {
			return AccessControlExposedHeaders
		},
	}
}

// CORS handles an options request to a given path and responds with the allowed headers.
// For non-OPTIONS requests it adds the required expose/allow headers
func CORS(opts *CorsOptions) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(HeaderAccessControlAllowOrigin, opts.AllowedOrigin(r))
			w.Header().Set(HeaderAccessControlAllowCredentials, AccessControlAllowedCredentials)
			w.Header().Set(HeaderAccessControlExposeHeaders, opts.ExposedHeaders(r))

			// if this is an options request, add the remaining pre flight check CORS headers to request
			if r.Method == http.MethodOptions {
				w.Header().Set(HeaderAccessControlAllowMethods, opts.AllowedMethods(r))
				w.Header().Set(HeaderAccessControlAllowHeaders, AccessControlAllowedHeaders)
				w.Header().Set(HeaderAccessControlMaxAge, AccessControlMaxAge)
			}

			// otherwise continue the request chain and add the required CORS headers
			h(w, r)
		}
	}
}
