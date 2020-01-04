package middleware_test

import (
	"net/http"

	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CORS", func() {
	dummyResponse := []byte("dummyresp")
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(dummyResponse)
	}

	It("should respond with values on options request", func() {
		req, err := http.NewRequest("OPTIONS", "", nil)
		req.Header.Set("Origin", "test")
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		opts := middleware.DefaultCorsOptions()
		middleware.CORS(opts)(dummyHandler)(writer, req)
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlAllowOrigin)).To(Equal("test"))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlAllowCredentials)).To(Equal(middleware.AccessControlAllowedCredentials))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlAllowMethods)).To(Equal(http.MethodPost))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlAllowHeaders)).To(Equal(middleware.AccessControlAllowedHeaders))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlExposeHeaders)).To(Equal(middleware.HeaderEcho))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlMaxAge)).To(Equal(middleware.AccessControlMaxAge))

		Expect(writer.StatusCode).To(Equal(http.StatusOK))
		Expect(writer.BodyData).To(Equal(dummyResponse))
	})
	It("should add correct header values on non-OPTIONS requests", func() {
		req, err := http.NewRequest("POST", "", nil)
		req.Header.Set("Origin", "test")
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		opts := middleware.DefaultCorsOptions()
		middleware.CORS(opts)(dummyHandler)(writer, req)
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlAllowOrigin)).To(Equal("test"))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlAllowCredentials)).To(Equal(middleware.AccessControlAllowedCredentials))
		Expect(writer.HeaderData.Get(middleware.HeaderAccessControlExposeHeaders)).To(Equal(middleware.AccessControlExposedHeaders))

		Expect(writer.StatusCode).To(Equal(http.StatusOK))
		Expect(writer.BodyData).To(Equal(dummyResponse))
	})
})
