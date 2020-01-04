package app_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/callstats-io/go-common/app"
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/metrics"
	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/response"
	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/dimfeld/httptreemux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	var (
		dummyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.OK(w, nil)
		})

		requestURL = func(port int, https bool) string {
			envVariable := app.EnvHTTPPort
			schema := "http"
			if https {
				envVariable = app.EnvHTTPSPort
				schema = "https"
			}

			if port == 0 {
				p, err := strconv.Atoi(os.Getenv(envVariable))
				Expect(err).To(BeNil())
				port = p
			}

			return schema + "://localhost:" + strconv.Itoa(port) + "/"
		}

		testRequest = func(port int, https bool) (*http.Response, error) {
			// allow server to boot up
			time.Sleep(50 * time.Millisecond)
			return testClient.Get(requestURL(port, https))
		}

		metricsResponse = func(port int) string {
			// allow server to boot up
			time.Sleep(50 * time.Millisecond)
			req, err := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+metrics.InternalMetricsPath, nil)
			Expect(err).To(BeNil())

			resp, err := testClient.Do(req)
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			raw, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			return string(raw)
		}

		newTestApp = func(ctx context.Context, port int, https bool) *app.App {
			a := app.NewApp(ctx)
			if port != 0 {
				if https {
					a.WithHTTPSPort(port)
				} else {
					a.WithHTTPPort(port)
				}
			}
			return a
		}
	)

	var _ = Describe("ServeHTTP", func() {
		ports := []int{
			0, // uses env port
			15601,
		}
		for idx := range ports {
			port := ports[idx]
			Context("Success", func() {
				It("should start a http only app in the specified address", func() {
					logBuffer := testutil.NewLogBuffer()

					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						ctx = log.WithLogger(ctx, logBuffer.Logger())
						newTestApp(ctx, port, false).ServeHTTP(dummyHandler)

						validateResponse(testRequest(port, false))
						Expect(logBuffer.String()).To(ContainSubstring(app.LogMsgServeHTTP))
					})

					Expect(err).To(BeNil())
					validateShutdown(testRequest(port, false))
					Expect(logBuffer.String()).To(ContainSubstring(app.LogMsgShutdownHTTP))
				})
			})
		}
	})

	var _ = Describe("ServeHTTPS", func() {
		var testVaultClientWrapper *fakeVaultClient

		BeforeEach(func() {
			testVaultClientWrapper = &fakeVaultClient{Client: testVaultClient}
		})

		ports := []int{
			0, // falls back to env
			15600,
		}
		for idx := range ports {
			port := ports[idx]
			Context("Success", func() {
				It("should start a https only app in the specified address", func() {
					logBuffer := testutil.NewLogBuffer()

					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						ctx = log.WithLogger(ctx, logBuffer.Logger())
						newTestApp(ctx, port, true).
							WithVaultClient(testVaultClientWrapper).
							ServeHTTPS(dummyHandler)

						validateResponse(testRequest(port, true))
						Expect(logBuffer.String()).To(ContainSubstring(app.LogMsgServeHTTPS))
					})
					Expect(err).To(BeNil())
					validateShutdown(testRequest(port, true))
					Expect(logBuffer.String()).To(ContainSubstring(app.LogMsgShutdownHTTPS))
				})
			})
			Context("Failure", func() {
				It("should return fail to start if TLS cert fetch fails", func() {
					logBuffer := testutil.NewLogBuffer()

					testVaultClientWrapper.tlsError = errors.New("FAKEERR")
					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						ctx = log.WithLogger(ctx, logBuffer.Logger())
						Expect(func() {
							newTestApp(ctx, port, true).
								WithVaultClient(testVaultClientWrapper).
								ServeHTTPS(dummyHandler)
						}).To(Panic())

						Expect(logBuffer.String()).To(ContainSubstring(app.LogErrTLSCertFetchHTTPS))
					})
					Expect(err).To(BeNil())
					validateShutdown(testRequest(port, true))
				})
				It("should return an error if vault client is nil and could not be set up", func() {
					logBuffer := testutil.NewLogBuffer()

					// expect this test to fail as there is not vault auth credentials to use for the default vault client
					prev := os.Getenv(vault.EnvVaultAppRoleCreds)
					os.Unsetenv(vault.EnvVaultAppRoleCreds)
					defer os.Setenv(vault.EnvVaultAppRoleCreds, prev)

					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						ctx = log.WithLogger(ctx, logBuffer.Logger())
						Expect(func() {
							newTestApp(ctx, port, true).ServeHTTPS(dummyHandler)
						}).To(Panic())

						Expect(logBuffer.String()).To(ContainSubstring(app.LogErrFailedToSetupVaultAutoRenewClient))
					})
					Expect(err).To(BeNil())
					validateShutdown(testRequest(port, true))
				})
			})
		}
	})

	var _ = Describe("Metrics", func() {
		var (
			port = 80

			router *httptreemux.TreeMux
		)

		BeforeEach(func() {
			metrics.ResetRegistry()
			router = httptreemux.New()
			baseMiddleware := middleware.NewConfiguratorWithCapacity(20)
			router.UsingContext().GET("/", baseMiddleware.ToHandler(dummyHandler, middleware.Metrics("rootHandler")))
			router.UsingContext().GET("/status", baseMiddleware.ToHandler(dummyHandler, middleware.Metrics("statusHandler")))
			router.UsingContext().GET(metrics.InternalMetricsPath, metrics.PrometheusEndpointWithoutCompression())
		})

		Context("Fetch", func() {
			It("should return empty metrics", func() {
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					Expect(func() {
						newTestApp(ctx, port, false).ServeHTTP(router)
					}).ToNot(Panic())

					resp := metricsResponse(port)

					// Check that no requests have been made
					Expect(resp).ToNot(MatchRegexp("http_request_count.*handler=\"rootHandler\""))
					Expect(resp).ToNot(MatchRegexp("http_request_count.*handler=\"statusHandler\""))
				})
				Expect(err).To(BeNil())
				validateShutdown(testRequest(port, false))
			})

			It("should return request metrics", func() {
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					Expect(func() {
						newTestApp(ctx, port, false).ServeHTTP(router)
					}).ToNot(Panic())

					_, err := testRequest(port, false)
					Expect(err).To(BeNil())
					_, err = testRequest(port, false)
					Expect(err).To(BeNil())

					resp := metricsResponse(port)

					// Check that two requests have been made to root path and no requests to /status path
					Expect(resp).To(ContainSubstring(`http_request_count{handler="rootHandler"} 2`))
					Expect(resp).ToNot(MatchRegexp("http_request_count.*handler=\"statusHandler\""))
				})
				Expect(err).To(BeNil())
				validateShutdown(testRequest(port, false))
			})
		})
	})
})
