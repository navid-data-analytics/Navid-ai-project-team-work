package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/response"
	"github.com/dimfeld/httptreemux"
)

// Constant responses
var (
	RespInvalidAppID, _  = json.Marshal(response.Response{Status: response.RespStatusError, Msg: "Invalid appID, must be a number"})
	RespInvalidUcID, _   = json.Marshal(response.Response{Status: response.RespStatusError, Msg: "Invalid ucID, must be a number"})
	RespInvalidConfID, _ = json.Marshal(response.Response{Status: response.RespStatusError, Msg: "Invalid confID: max length is " + strconv.Itoa(MaxConfIDLength) + " characters in encoded form"})
)

// RequireValidConfID validates that the confID is max MaxConfIDLength characters in url encoded form
func RequireValidConfID(paramName string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params := httptreemux.ContextParams(r.Context())
			if len(url.QueryEscape(params[paramName])) > MaxConfIDLength {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware/validations").Error(LogMsgErrInvalidConfID, log.String(paramName, url.QueryEscape(params[paramName])))
				response.BadRequest(w, RespInvalidConfID)
				return
			}
			h(w, r)
		}
	}
}

// RequireValidAppID validates that the appID is a number
func RequireValidAppID(paramName string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params := httptreemux.ContextParams(r.Context())
			if !ValidAppIDRegExp.MatchString(params[paramName]) {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware/validations").Error(LogMsgErrInvalidAppID, log.String(paramName, url.QueryEscape(params[paramName])))
				response.BadRequest(w, RespInvalidAppID)
				return
			}
			h(w, r)
		}
	}
}

// RequireValidUcID validates that the ucID is a number
func RequireValidUcID(paramName string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params := httptreemux.ContextParams(r.Context())
			if !ValidUcIDRegExp.MatchString(params[paramName]) {
				log.FromContextWithPackageName(r.Context(), "go-common/middleware/validations").Error(LogMsgErrInvalidUcID, log.String(paramName, url.QueryEscape(params[paramName])))
				response.BadRequest(w, RespInvalidUcID)
				return
			}
			h(w, r)
		}
	}
}
