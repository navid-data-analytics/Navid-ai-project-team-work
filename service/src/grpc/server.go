package grpc

import (
	"net"

	context "golang.org/x/net/context"

	"github.com/callstats-io/ai-decision/service/gen/protos"
	grpc_utils "github.com/callstats-io/go-common/grpc"
	"github.com/callstats-io/go-common/grpc/metrics"
	"github.com/callstats-io/go-common/grpc/panic"
	"google.golang.org/grpc"
)

// Server thin wrapper for grpcServer
type Server struct {
	grpcServer *grpc.Server
}

// NewServer builds new Server
func NewServer(ctx context.Context, msrv protos.AIDecisionMessageServiceServer, ssrv protos.AIDecisionStateServiceServer) (*Server, error) {
	s := &Server{}
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_utils.ChainUnaryServerInterceptors(
			metrics.UnaryServerInterceptor,
			panic.UnaryServerInterceptor(panic.LoggingRecovery(ctx)),
		)),
		grpc.StreamInterceptor(
			grpc_utils.ChainStreamServer(
				metrics.StreamServerInterceptor,
				panic.StreamServerInterceptor(panic.LoggingRecovery(ctx))),
		),
	)
	err := metrics.Register(ctx)
	if err != nil {
		return nil, err
	}

	protos.RegisterAIDecisionMessageServiceServer(s.grpcServer, msrv)
	protos.RegisterAIDecisionStateServiceServer(s.grpcServer, ssrv)
	return s, nil
}

// Serve starts serving grpc requests
func (s *Server) Serve(ctx context.Context, listener net.Listener) error {
	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()
	return s.grpcServer.Serve(listener)
}
