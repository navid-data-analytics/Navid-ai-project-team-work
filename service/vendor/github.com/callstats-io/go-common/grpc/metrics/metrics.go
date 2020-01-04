package metrics

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strings"
	"time"
)

// UnaryServerInterceptor intercepts gRPC calls to gather metrics.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	serviceName, method := splitMethodName(info.FullMethod)

	requestCounter.WithLabelValues(serviceName, method).Inc()

	start := time.Now()
	resp, err := handler(ctx, req)

	statusCode := grpc.Code(err).String()
	responseTimer.WithLabelValues(serviceName, method, statusCode).Observe(float64(time.Since(start)) / float64(time.Second))

	return resp, err
}

func splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}

// StreamServerInterceptor intercepts gRPC stream calls to gather metrics.
func StreamServerInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	serviceName, method := splitMethodName(info.FullMethod)
	streamType := "server"
	if info.IsClientStream {
		streamType = "client"
	}
	start := time.Now()
	err = handler(srv, stream)
	statusCode := grpc.Code(err).String()
	streamRequestsTimer.WithLabelValues(serviceName, method, statusCode, streamType).Observe(float64(time.Since(start)) / float64(time.Second))
	return err
}
