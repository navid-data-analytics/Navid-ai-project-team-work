package middleware_test

import (
	"net/http"

	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/response"
	"github.com/callstats-io/go-common/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequireHTTP2", func() {
	It("should not call handler if HTTP major is 1", func() {
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		req.ProtoMajor = 1
		middleware.RequireHTTP2()(func(w http.ResponseWriter, r *http.Request) {
			// Expect this to never be called, if this is called then the RequireHTTP2 is broken
			Fail("Expected the handler not to be called in RequireHTTP2")
		})(writer, req)
		Expect(writer.BodyData).To(Equal(response.RespInvalidHTTPVersion))
		Expect(writer.StatusCode).To(Equal(http.StatusHTTPVersionNotSupported))
	})

	It("should call handler if HTTP major is 2", func() {
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		called := false
		req.ProtoMajor = 2
		middleware.RequireHTTP2()(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})(writer, req)
		Expect(called).To(BeTrue())
	})
})
