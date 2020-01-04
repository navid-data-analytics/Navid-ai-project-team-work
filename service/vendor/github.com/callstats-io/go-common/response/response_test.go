package response_test

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/callstats-io/go-common/response"
	"github.com/callstats-io/go-common/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalError", func() {
	It("should marshal the error to response json", func() {
		err := errors.New("test error")
		resp := response.MarshalError(err)
		Expect(resp).To(Equal([]byte("{\"status\":\"error\",\"msg\":\"" + err.Error() + "\"}")))
	})
})

var _ = Describe("InvalidHTTPVersion", func() {
	It("should write HTTP status code 505 and a msg", func() {
		w := &testutil.FakeResponseWriter{}
		response.InvalidHTTPVersion(w)
		Expect(w.StatusCode).To(Equal(http.StatusHTTPVersionNotSupported))
		Expect(w.BodyData).To(Equal(response.RespInvalidHTTPVersion))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("BadRequest", func() {
	It("should write HTTP status code 400 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.BadRequest(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("UnprocessableEntity", func() {
	It("should write HTTP status code 422 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.UnprocessableEntity(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusUnprocessableEntity))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("NotFound", func() {
	It("should write HTTP status code 404 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.NotFound(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusNotFound))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
	It("should write HTTP status code 404 and not found msg if content is nil", func() {
		w := &testutil.FakeResponseWriter{}
		response.NotFound(w, nil)
		Expect(w.StatusCode).To(Equal(http.StatusNotFound))
		Expect(w.BodyData).To(Equal(response.RespNotFound))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("OK", func() {
	It("should write HTTP status code 200 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.OK(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusOK))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
	It("should write HTTP status code 200 and success msg if content is nil", func() {
		w := &testutil.FakeResponseWriter{}
		response.OK(w, nil)
		Expect(w.StatusCode).To(Equal(http.StatusOK))
		Expect(w.BodyData).To(Equal(response.RespOK))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("InternalServerError", func() {
	It("should write HTTP status code 500 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.InternalServerError(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
	It("should write HTTP status code 500 and error msg if content is nil", func() {
		w := &testutil.FakeResponseWriter{}
		response.InternalServerError(w, nil)
		Expect(w.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(w.BodyData).To(Equal(response.RespInternalServerError))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("RequiredServiceUnavailable", func() {
	It("should write HTTP status code 503 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.RequiredServiceUnavailable(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusServiceUnavailable))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
	It("should write HTTP status code 503 and error msg if content is nil", func() {
		w := &testutil.FakeResponseWriter{}
		response.RequiredServiceUnavailable(w, nil)
		Expect(w.StatusCode).To(Equal(http.StatusServiceUnavailable))
		Expect(w.BodyData).To(Equal(response.RespRequiredServiceUnavailable))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("Unauthorized", func() {
	It("should write HTTP status code 401 and the passed in data", func() {
		w := &testutil.FakeResponseWriter{}
		expContent := []byte("expcontent")
		response.Unauthorized(w, expContent)
		Expect(w.StatusCode).To(Equal(http.StatusUnauthorized))
		Expect(w.BodyData).To(Equal(expContent))
		Expect(w.Header().Get(response.HeaderContentType)).To(Equal(response.ApplicationJSON))
		Expect(w.Header().Get(response.HeaderContentLength)).To(Equal(strconv.Itoa(len(w.BodyData))))
	})
})

var _ = Describe("NoContent", func() {
	It("should write HTTP status code 204 with no body", func() {
		w := &testutil.FakeResponseWriter{}
		response.NoContent(w)
		Expect(w.StatusCode).To(Equal(http.StatusNoContent))
	})
})
