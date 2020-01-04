package grpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	someServiceName  = "SomeService.StreamMethod"
	parentUnaryInfo  = &grpc.UnaryServerInfo{FullMethod: someServiceName}
	parentStreamInfo = &grpc.StreamServerInfo{
		FullMethod:     someServiceName,
		IsServerStream: true,
	}
	someValue     = 1
	parentContext = context.WithValue(context.TODO(), key("parent"), someValue)
)

type key string

func TestChainUnaryServer(t *testing.T) {
	input := "input"
	output := "output"

	first := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requireContextValue(ctx, t, key("parent"), "first interceptor must know the parent context value")
		require.Equal(t, parentUnaryInfo, info, "first interceptor must know the someUnaryServerInfo")
		ctx = context.WithValue(ctx, key("first"), 1)
		return handler(ctx, req)
	}
	second := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requireContextValue(ctx, t, key("parent"), "second interceptor must know the parent context value")
		requireContextValue(ctx, t, key("first"), "second interceptor must know the first context value")
		require.Equal(t, parentUnaryInfo, info, "second interceptor must know the someUnaryServerInfo")
		ctx = context.WithValue(ctx, key("second"), 1)
		return handler(ctx, req)
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		require.EqualValues(t, input, req, "handler must get the input")
		requireContextValue(ctx, t, key("parent"), "handler must know the parent context value")
		requireContextValue(ctx, t, key("first"), "handler must know the first context value")
		requireContextValue(ctx, t, key("second"), "handler must know the second context value")
		return output, nil
	}

	chain := ChainUnaryServerInterceptors(first, second)
	out, _ := chain(parentContext, input, parentUnaryInfo, handler)
	require.EqualValues(t, output, out, "chain must return handler's output")
}

func TestChainStreamServer(t *testing.T) {
	someService := &struct{}{}
	recvMessage := "received"
	sentMessage := "sent"
	outputError := fmt.Errorf("some error")

	first := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		requireContextValue(stream.Context(), t, key("parent"), "first interceptor must know the parent context value")
		require.Equal(t, parentStreamInfo, info, "first interceptor must know the parentStreamInfo")
		require.Equal(t, someService, srv, "first interceptor must know someService")
		wrapped := WrapServerStream(stream)
		wrapped.WrappedContext = context.WithValue(stream.Context(), key("first"), 1)
		return handler(srv, wrapped)
	}
	second := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		requireContextValue(stream.Context(), t, key("parent"), "second interceptor must know the parent context value")
		requireContextValue(stream.Context(), t, key("parent"), "second interceptor must know the first context value")
		require.Equal(t, parentStreamInfo, info, "second interceptor must know the parentStreamInfo")
		require.Equal(t, someService, srv, "second interceptor must know someService")
		wrapped := WrapServerStream(stream)
		wrapped.WrappedContext = context.WithValue(stream.Context(), key("second"), 1)
		return handler(srv, wrapped)
	}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		require.Equal(t, someService, srv, "handler must know someService")
		requireContextValue(stream.Context(), t, key("parent"), "handler must know the parent context value")
		requireContextValue(stream.Context(), t, key("first"), "handler must know the first context value")
		requireContextValue(stream.Context(), t, key("second"), "handler must know the second context value")
		require.NoError(t, stream.RecvMsg(recvMessage), "handler must have access to stream messages")
		require.NoError(t, stream.SendMsg(sentMessage), "handler must be able to send stream messages")
		return outputError
	}
	fakeStream := &fakeServerStream{ctx: parentContext, recvMessage: recvMessage}
	chain := ChainStreamServer(first, second)
	err := chain(someService, fakeStream, parentStreamInfo, handler)
	require.Equal(t, outputError, err, "chain must return handler's error")
	require.Equal(t, sentMessage, fakeStream.sentMessage, "handler's sent message must propagate to stream")
}

func requireContextValue(ctx context.Context, t *testing.T, key key, msg ...interface{}) {
	val := ctx.Value(key)
	require.NotNil(t, val, msg...)
	require.Equal(t, someValue, val, msg...)
}
