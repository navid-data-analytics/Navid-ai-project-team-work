package metrics_test

import (
	grpccommon "github.com/callstats-io/go-common/grpc"
	grpcmetrics "github.com/callstats-io/go-common/grpc/metrics"
	"github.com/callstats-io/go-common/metrics"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	hwprotos "google.golang.org/grpc/examples/helloworld/helloworld"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testCtx = context.Background()
)

var _ = Describe("gRPC metrics", func() {
	BeforeEach(func() {
		metrics.ResetRegistry()
		err := grpcmetrics.Register(testCtx)
		Expect(err).To(BeNil())
	})

	Context("Register", func() {
		for metricName, metricType := range map[string]dto.MetricType{
			"grpc_request_count":         dto.MetricType_COUNTER,
			"grpc_response_time_seconds": dto.MetricType_HISTOGRAM,
		} {
			It("should initialize "+metricName, func() {
				_, err := testClient.SayHello(testCtx, &hwprotos.HelloRequest{Name: "olleH"})
				Expect(grpc.Code(err)).To(Equal(codes.OK))

				metricFamilies, err := prometheus.DefaultGatherer.Gather()
				Expect(err).To(BeNil())
				Expect(len(metricFamilies)).To(Equal(2))
				for _, metricFamily := range metricFamilies {
					if metricFamily.GetName() == metricName && metricFamily.GetType() == metricType {
						return
					}
				}
				Fail("Metric " + metricName + " wasn't registered")
			})
		}

		It("should not fail if called the second time", func() {
			err := grpcmetrics.Register(testCtx)
			Expect(err).To(BeNil())
		})
	})

	Context("ChainUnaryServerInterceptor", func() {
		var (
			timesIntercepted int
			wasHandlerCalled bool

			countingInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				timesIntercepted++
				return handler(ctx, req)
			}
			handlerFunc = func(ctx context.Context, req interface{}) (interface{}, error) {
				wasHandlerCalled = true
				return "foo", nil
			}
		)

		BeforeEach(func() {
			timesIntercepted = 0
			wasHandlerCalled = false
		})

		It("should 'chain' one interceptor", func() {
			chain := grpccommon.ChainUnaryServerInterceptors(
				countingInterceptor,
			)

			_, err := chain(testCtx, &hwprotos.HelloRequest{Name: "olleH"}, &grpc.UnaryServerInfo{FullMethod: "ChainedService.Handler"}, handlerFunc)
			Expect(err).To(BeNil())
			Expect(wasHandlerCalled).To(BeTrue())
			Expect(timesIntercepted).To(Equal(1))
		})

		It("should chain two interceptors", func() {
			chain := grpccommon.ChainUnaryServerInterceptors(
				countingInterceptor,
				countingInterceptor,
			)

			_, err := chain(testCtx, &hwprotos.HelloRequest{Name: "olleH"}, &grpc.UnaryServerInfo{FullMethod: "ChainedService.Handler"}, handlerFunc)
			Expect(err).To(BeNil())
			Expect(wasHandlerCalled).To(BeTrue())
			Expect(timesIntercepted).To(Equal(2))
		})
	})

	Context("UnaryServerInterceptor", func() {
		It("should increment the request counter", func() {
			_, err := testClient.SayHello(testCtx, &hwprotos.HelloRequest{Name: "olleH"})
			Expect(grpc.Code(err)).To(Equal(codes.OK))
			Expect(requestCount("helloworld.Greeter", "SayHello")).To(Equal(1))
		})

		It("should record response time", func() {
			_, err := testClient.SayHello(testCtx, &hwprotos.HelloRequest{Name: "olleH"})
			Expect(grpc.Code(err)).To(Equal(codes.OK))
			whichBucket := 0 // we assume that the response time fell into the first bucket
			upperBound, count := responseTimes("helloworld.Greeter", "SayHello", whichBucket)
			Expect(upperBound).To(Equal(0.005))
			Expect(count).To(Equal(1))
		})
	})
})
