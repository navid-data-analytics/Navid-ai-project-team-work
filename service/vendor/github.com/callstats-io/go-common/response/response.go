package response

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Header keys
const (
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
)

// Header values
const (
	ApplicationJSON = "application/json"
)

// Constants
const (
	MaxConfIDLength = 256
)

// Response statuses
const (
	RespStatusSuccess = "success"
	RespStatusError   = "error"
)

// Constant responses
var (
	RespInvalidHTTPVersion, _         = json.Marshal(Response{Status: RespStatusError, Msg: "Invalid HTTP major: only version 2 is supported."})
	RespInternalServerError, _        = json.Marshal(Response{Status: RespStatusError, Msg: "Internal error"})
	RespRequiredServiceUnavailable, _ = json.Marshal(Response{Status: RespStatusError, Msg: "Required service(s) temporarily unavailable"})
	RespNotFound, _                   = json.Marshal(Response{Status: RespStatusError, Msg: "Requested resource not found"})
	RespOK, _                         = json.Marshal(Response{Status: RespStatusSuccess})
)

// Response implements the data format all errors from our api will return
type Response struct {
	Status string `json:"status"` // HTTP status code
	Msg    string `json:"msg"`    // Description
}

// MarshalError returns a json response with status set to error and the errors message as msg
func MarshalError(err error) []byte {
	resp, _ := json.Marshal(&Response{
		Status: RespStatusError,
		Msg:    err.Error(),
	})
	return resp
}

// InvalidHTTPVersion responds with http 400, content type of application/json and appropriate content length with {"status":"error","msg":"Invalid HTTP major: only version 2 is supported."}
func InvalidHTTPVersion(w http.ResponseWriter) {
	writeCommonHeaders(w, len(RespInvalidHTTPVersion))
	w.WriteHeader(http.StatusHTTPVersionNotSupported)
	w.Write(RespInvalidHTTPVersion)
}

// BadRequest responds with http 400, content type of application/json and appropriate content length with the given resp bytes as payload
func BadRequest(w http.ResponseWriter, resp []byte) {
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}

// UnprocessableEntity responds with http 422, content type of application/json and appropriate content length with the given resp bytes as payload
func UnprocessableEntity(w http.ResponseWriter, resp []byte) {
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(resp)
}

// PreconditionFailed responds with http 412, content type of application/json and appropriate content length with the given resp bytes as payload
func PreconditionFailed(w http.ResponseWriter, resp []byte) {
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusPreconditionFailed)
	w.Write(resp)
}

// NotFound responds with http 404, content type of application/json and appropriate content length with the given resp bytes as payload
// If resp is nil it responds with content {"status":"error","msg":"Requested resource not found"}
func NotFound(w http.ResponseWriter, resp []byte) {
	if resp == nil {
		resp = RespNotFound
	}
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusNotFound)
	w.Write(resp)
}

// OK responds with http 200, content type of application/json and appropriate content length with the given resp bytes as payload
// If resp is nil it responds with content {"status":"success","msg":""}
func OK(w http.ResponseWriter, resp []byte) {
	if resp == nil {
		resp = RespOK
	}
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// InternalServerError responds with http 500, content type of application/json and appropriate content length
// If resp is nil it responds with content {"status":"error","msg":"Internal server error"}
func InternalServerError(w http.ResponseWriter, resp []byte) {
	if resp == nil {
		resp = RespInternalServerError
	}
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

// RequiredServiceUnavailable responds with http 503, content type of application/json and appropriate content length
// If resp is nil it responds with content {"status":"error","msg":"Required service(s) temporarily unavailable"}
func RequiredServiceUnavailable(w http.ResponseWriter, resp []byte) {
	if resp == nil {
		resp = RespRequiredServiceUnavailable
	}
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write(resp)
}

// Unauthorized responds with http 401, content type of application/json and appropriate content length with the given resp bytes as payload
func Unauthorized(w http.ResponseWriter, resp []byte) {
	writeCommonHeaders(w, len(resp))
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(resp)
}

// NoContent responds with http 204, content type of application/json and appropriate content length without any body
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeCommonHeaders(w http.ResponseWriter, contentLength int) {
	w.Header().Set(HeaderContentType, ApplicationJSON)
	if contentLength != 0 {
		w.Header().Set(HeaderContentLength, strconv.Itoa(contentLength))
	}
}
