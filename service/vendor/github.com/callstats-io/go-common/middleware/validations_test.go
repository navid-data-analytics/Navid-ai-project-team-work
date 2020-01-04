package middleware_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/testutil"
	"github.com/dimfeld/httptreemux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequireValidConfID", func() {
	It("should be valid if the confID is shorter than or equal to "+strconv.Itoa(middleware.MaxConfIDLength)+" characters", func() {
		// expect to succeed with exact length confID
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		// expect to succeed with equal length confID
		params := map[string]string{
			"confID": strings.Repeat("a", middleware.MaxConfIDLength),
		}
		req = req.WithContext(context.WithValue(req.Context(), httptreemux.ParamsContextKey, params))
		called := false
		middleware.RequireValidConfID("confID")(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})(writer, req)
		Expect(called).To(BeTrue())
	})

	It("should be invalid if the confID is longer than "+strconv.Itoa(middleware.MaxConfIDLength)+" characters", func() {
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())

		// expect to fail with one char too long confID
		params := map[string]string{
			"confID": strings.Repeat("a", middleware.MaxConfIDLength+1),
		}
		req = req.WithContext(context.WithValue(req.Context(), httptreemux.ParamsContextKey, params))

		writer := &testutil.FakeResponseWriter{}
		middleware.RequireValidConfID("confID")(func(w http.ResponseWriter, r *http.Request) {
			// Expect this to never be called, if this is called then the RequireValidConfID is broken
			Fail("Expected the handler not to be called in RequireValidConfID")
		})(writer, req)

		Expect(writer.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(writer.BodyData).To(Equal(middleware.RespInvalidConfID))
	})
})

var _ = Describe("RequireValidAppID", func() {
	It("should succeed if appID is a stringified integer", func() {
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())
		writer := &testutil.FakeResponseWriter{}

		params := map[string]string{
			"appID": "12345",
		}
		req = req.WithContext(context.WithValue(req.Context(), httptreemux.ParamsContextKey, params))
		called := false
		middleware.RequireValidAppID("appID")(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})(writer, req)
		Expect(called).To(BeTrue())
	})

	It("should fail if appID is not a stringified integer", func() {
		req, err := http.NewRequest("POST", "", nil)
		Expect(err).To(BeNil())

		for _, appID := range []string{"abc1234", "1234abc", "123abc345"} {
			params := map[string]string{
				"appID": appID,
			}
			req = req.WithContext(context.WithValue(req.Context(), httptreemux.ParamsContextKey, params))

			writer := &testutil.FakeResponseWriter{}
			middleware.RequireValidAppID("appID")(func(w http.ResponseWriter, r *http.Request) {
				// Expect this to never be called, if this is called then the RequireValidAppID is broken
				Fail("Expected the handler not to be called in RequireValidAppID")
			})(writer, req)

			Expect(writer.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(writer.BodyData).To(Equal(middleware.RespInvalidAppID))
		}
	})
})

var _ = Describe("RequireValidUcID", func() {
	Context("Success", func() {
		sharedSuccessTestCase := func(ucID string) {
			req, err := http.NewRequest("POST", "", nil)
			Expect(err).To(BeNil())
			writer := &testutil.FakeResponseWriter{}

			params := map[string]string{
				"ucID": ucID,
			}
			req = req.WithContext(context.WithValue(req.Context(), httptreemux.ParamsContextKey, params))
			called := false
			middleware.RequireValidUcID("ucID")(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})(writer, req)
			Expect(called).To(BeTrue())
		}

		It("should pass if ucID is a stringified integer", func() {
			sharedSuccessTestCase("12345")
		})

		It("should pass if ucID is a stringified integer with 'k' prefix", func() {
			sharedSuccessTestCase("k12345")
		})

		It("should pass if ucID is a stringified integer with 'K' prefix", func() {
			sharedSuccessTestCase("K12345")
		})

		It("should pass if ucID is a stringified integer with 'p' prefix", func() {
			sharedSuccessTestCase("p12345")
		})

		It("should pass if ucID is a stringified integer with 'P' prefix", func() {
			sharedSuccessTestCase("P12345")
		})
	})
	Context("Failure", func() {
		It("should fail if ucID is not a stringified integer (with p, P, k or K prefix)", func() {
			req, err := http.NewRequest("POST", "", nil)
			Expect(err).To(BeNil())

			for _, ucID := range []string{"abc1234", "1234abc", "123abc345"} {
				params := map[string]string{
					"ucID": ucID,
				}
				req = req.WithContext(context.WithValue(req.Context(), httptreemux.ParamsContextKey, params))

				writer := &testutil.FakeResponseWriter{}
				middleware.RequireValidUcID("ucID")(func(w http.ResponseWriter, r *http.Request) {
					// Expect this to never be called, if this is called then the RequireValidUcID is broken
					Fail("Expected the handler not to be called in RequireValidUcID")
				})(writer, req)

				Expect(writer.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(writer.BodyData).To(Equal(middleware.RespInvalidUcID))
			}
		})
	})
})
