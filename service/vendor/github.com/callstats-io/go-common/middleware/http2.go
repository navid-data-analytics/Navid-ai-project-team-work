package middleware

import (
	"net/http"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/response"
)

// RequireHTTP2 checks that the request was made with http or responds with http.StatusHTTPVersionNotSupported
func RequireHTTP2() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor != 2 {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware/http2").Error(LogMsgErrInvalidHTTPVersion)
				response.InvalidHTTPVersion(w)
				return
			}
			h(w, r)
		}
	}
}
