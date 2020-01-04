package middleware

import "net/http"

// Echo add the value passed in X-Csio-Echo to the response
func Echo() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(HeaderEcho, r.Header.Get(HeaderEcho))
			h(w, r)
		}
	}
}
