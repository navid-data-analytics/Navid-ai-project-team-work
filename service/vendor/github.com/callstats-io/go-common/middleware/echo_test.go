package middleware_test

import (
	"net/http"

	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Echo", func() {
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {}

	It("should add the echo value from header", func() {
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		testValues := map[string]string{
			"X-CSIO-ECHO": "echoupper",
			"x-csio-echo": "echolower",
			"X-Csio-Echo": "echocanonical",
		}
		for header, value := range testValues {
			req.Header.Set(header, value)
			middleware.Echo()(dummyHandler)(writer, req)
			Expect(writer.HeaderData.Get(middleware.HeaderEcho)).To(Equal(value))
		}
	})
})
