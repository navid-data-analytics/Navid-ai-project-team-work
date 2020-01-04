package metrics_test

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	grpcmetrics "github.com/callstats-io/go-common/grpc/metrics"
	"github.com/callstats-io/go-common/metrics"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	hwprotos "google.golang.org/grpc/examples/helloworld/helloworld"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testListener net.Listener
	testServer   *grpc.Server
	testConn     *grpc.ClientConn
	testClient   hwprotos.GreeterClient

	grpcPort = os.Getenv("GRPC_PORT")
)

// greeterServer is used to implement helloworld.GreeterServer.
type greeterServer struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeterServer) SayHello(ctx context.Context, in *hwprotos.HelloRequest) (*hwprotos.HelloReply, error) {
	return &hwprotos.HelloReply{Message: "Hello " + in.Name}, nil
}

var _ = BeforeSuite(func() {
	testServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpcmetrics.UnaryServerInterceptor),
	)
	hwprotos.RegisterGreeterServer(testServer, &greeterServer{})

	var err error
	testListener, err = net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	Expect(err).To(BeNil())

	started := make(chan bool)
	go func() {
		defer GinkgoRecover()
		time.AfterFunc(time.Millisecond*10, func() {
			started <- true
		})
		// ignore gRPC error when closing connection
		if err := testServer.Serve(testListener); !strings.Contains(err.Error(), "use of closed network connection") {
			Expect(err).To(BeNil())
		}
	}()
	<-started

	testConn, err = grpc.Dial("localhost:"+grpcPort, grpc.WithInsecure())
	Expect(err).To(BeNil())
	testClient = hwprotos.NewGreeterClient(testConn)
})

var _ = AfterSuite(func() {
	testServer.Stop()
	testConn.Close()
})

func gatherMetrics(metricType dto.MetricType, metricName, serviceName, method string) *dto.Metric {
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	Expect(err).To(BeNil())
	for _, metricFamily := range metricFamilies { // []*dto.MetricFamily
		if metricFamily.GetType() == metricType && metricFamily.GetName() == metricName {
			for _, metric := range metricFamily.GetMetric() { // []*dto.Metric
				wantedLabels := 2 // "service_name" and "method"

				for _, label := range metric.GetLabel() { // []*dto.LabelPair
					if label.GetName() == grpcmetrics.LabelServiceName && label.GetValue() == serviceName {
						wantedLabels--
					}
					if label.GetName() == grpcmetrics.LabelMethod && label.GetValue() == method {
						wantedLabels--
					}
				}

				if wantedLabels == 0 {
					// found both labels we were looking for
					return metric
				}
			}
		}
	}
	return nil
}

func requestCount(serviceName, method string) int {
	metric := gatherMetrics(dto.MetricType_COUNTER, "grpc_"+metrics.MetricRequestCount, serviceName, method)
	Expect(metric).ToNot(BeNil())
	return int(metric.GetCounter().GetValue())
}

func responseTimes(serviceName, method string, whichBucket int) (float64, int) {
	metric := gatherMetrics(dto.MetricType_HISTOGRAM, "grpc_"+metrics.MetricResponseTime, serviceName, method)
	Expect(metric).ToNot(BeNil())
	histogram := metric.GetHistogram()
	Expect(histogram).ToNot(BeNil())
	bucket := histogram.GetBucket()[whichBucket] // get the nth bucket
	Expect(bucket).ToNot(BeNil())
	upperBound, count := bucket.GetUpperBound(), int(bucket.GetCumulativeCount())
	return upperBound, count
}

// ===== TEST SETUP =====
func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Metrics Black Box Suite")
}
