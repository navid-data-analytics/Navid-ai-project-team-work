package panic

import (
	"errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

// TestUnaryInterceptorMustCatchPanic tests unary interceptor
func TestUnaryInterceptorMustCatchPanic(t *testing.T) {
	var panicRecovered interface{}
	someErr := errors.New("Panic has been catched")
	panicMessage := "Panic!"
	recovery := func(p interface{}) (err error) {
		panicRecovered = p
		return someErr
	}
	handler := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
		panic(panicMessage)
	})
	ctx := context.Background()
	interceptor := UnaryServerInterceptor(recovery)
	_, err := interceptor(ctx, nil, nil, handler)
	require.NotNil(t, panicRecovered)
	require.Equal(t, panicRecovered.(string), panicMessage)
	require.Equal(t, err, someErr)
}

// TestUnaryInterceptorMoPanic tests unary interceptor in case of no panic
func TestUnaryInterceptorMoPanic(t *testing.T) {

	response := struct {
		label string
	}{
		label: "fake response",
	}
	recovery := func(p interface{}) (err error) {
		panic("Interceptor should not run anytime")
	}
	handler := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
		return response, nil
	})
	ctx := context.Background()
	interceptor := UnaryServerInterceptor(recovery)
	handlerResp, err := interceptor(ctx, nil, nil, handler)
	require.Nil(t, err, "Error must not happens")
	require.Equal(t, response, handlerResp, "Response must be received")
}

// TestStreamInterceptorMustCatchPanic tests stream interceptor
func TestStreamInterceptorMustCatchPanic(t *testing.T) {
	var panicRecovered interface{}
	someErr := errors.New("Panic has been catched")
	panicMessage := "Panic!"
	recovery := func(p interface{}) (err error) {
		panicRecovered = p
		return someErr
	}
	handler := func(srv interface{}, ss grpc.ServerStream) error {
		panic(panicMessage)
	}
	interceptor := StreamServerInterceptor(recovery)
	err := interceptor(nil, nil, nil, handler)
	require.NotNil(t, panicRecovered)
	require.Equal(t, panicRecovered.(string), "Panic!")
	require.Equal(t, err, someErr)
}

// TestStreamInterceptorNoPanic tests stream interceptor in case of no panic
func TestStreamInterceptorNoPanic(t *testing.T) {
	handlerRunned := false
	recovery := func(p interface{}) (err error) {
		panic("Interceptor should not run anytime")
	}
	handler := func(srv interface{}, ss grpc.ServerStream) error {
		handlerRunned = true
		return nil
	}
	interceptor := StreamServerInterceptor(recovery)
	err := interceptor(nil, nil, nil, handler)
	require.Nil(t, err, "Error must not happens")
	require.Equal(t, handlerRunned, true, "Handler should be called")
}
