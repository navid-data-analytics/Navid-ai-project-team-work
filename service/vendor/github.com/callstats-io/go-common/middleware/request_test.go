package middleware_test

import (
	"net/http"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/request"
	"github.com/callstats-io/go-common/response"
	"github.com/callstats-io/go-common/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestContextLogger", func() {
	Context("Success", func() {
		It("should add a request id to the logger", func() {
			req, err := http.NewRequest("POST", "", nil)
			Expect(err).To(BeNil())
			writer := &testutil.FakeResponseWriter{}
			expResp := []byte("RESPRCTXLOGSUCCESS")

			logBuffer := testutil.NewLogBuffer()
			req = req.WithContext(log.WithLogger(req.Context(), logBuffer.Logger()))

			middleware.RequestContextLogger()(func(w http.ResponseWriter, r *http.Request) {
				requestID := request.IDFromContext(r.Context())
				Expect(requestID).ToNot(BeEmpty())
				logger := log.FromContextWithPackageName(r.Context(), "go-common/middleware/request")
				logger.Info("ABC")
				Expect(logBuffer.String()).To(ContainSubstring(requestID))
				response.OK(w, expResp)
			})(writer, req)

			// check handler was called
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(expResp))
		})
	})
})

var _ = Describe("RequestLifetimeLogger", func() {
	Context("Success", func() {
		It("should log the remote IP for the request", func() {
			req, err := http.NewRequest("POST", "", nil)
			Expect(err).To(BeNil())
			writer := &testutil.FakeResponseWriter{}
			expResp := []byte("RESPRCTXLOGSUCCESS")

			logBuffer := testutil.NewLogBuffer()
			req = req.WithContext(log.WithLogger(req.Context(), logBuffer.Logger()))
			req.RemoteAddr = "127.0.0.1:12345"

			middleware.RequestLifetimeLogger()(func(w http.ResponseWriter, r *http.Request) {
				buf := logBuffer.String()
				ip := request.IPFromAddr(r.RemoteAddr)
				Expect(buf).To(ContainSubstring(middleware.LogMsgStartRequest))
				Expect(buf).To(ContainSubstring(middleware.LogKeyRemoteIP))
				Expect(buf).To(ContainSubstring(ip))
				response.OK(w, expResp)
			})(writer, req)

			// check handler was called
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(expResp))

			// test log messages
			Expect(logBuffer.String()).To(ContainSubstring(middleware.LogMsgFinishRequest))
		})
	})
})
