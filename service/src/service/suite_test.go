package service_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/callstats-io/ai-decision/service/gen/protos"
	"github.com/callstats-io/ai-decision/service/src/flowdock"
	sgrpc "github.com/callstats-io/ai-decision/service/src/grpc"
	"github.com/callstats-io/ai-decision/service/src/service"
	"github.com/callstats-io/ai-decision/service/src/storage/mocks"
	"google.golang.org/grpc"
)

var (
	testCtx, testCtxCancel = context.WithCancel(context.Background())
	testServer             *sgrpc.Server
	testClientConn         *grpc.ClientConn
	testMessageClient      protos.AIDecisionMessageServiceClient
	testStateClient        protos.AIDecisionStateServiceClient
	mockStorage            *mocks.Storage
)

func mustBeNil(err error) {
	if err != nil {
		panic(err)
	}
}

func suiteSetup() {
	mockStorage = mocks.NewMockedStorage()
	flowdockClient := flowdock.NewClient("")
	aiDecisionMessageService, err := service.NewAIDecisionMessageService(mockStorage, flowdockClient)
	mustBeNil(err)
	aiDecisionStateService, err := service.NewAIDecisionStateService(mockStorage)
	mustBeNil(err)
	testServer, err = sgrpc.NewServer(testCtx, aiDecisionMessageService, aiDecisionStateService)
	mustBeNil(err)

	testServerListener, err := net.Listen("tcp", fmt.Sprintf("localhost:0"))
	mustBeNil(err)

	go func() {
		defer func() { mustBeNil(testServerListener.Close()) }()
		// both errors are because of stopping, one from net listener and one from grpc itself
		if err := testServer.Serve(testCtx, testServerListener); err != grpc.ErrServerStopped && !strings.Contains(err.Error(), "use of closed network connection") {
			mustBeNil(err)
		}

	}()
	testClientConn, err := grpc.Dial(testServerListener.Addr().String(), grpc.WithInsecure())
	mustBeNil(err)
	testMessageClient = protos.NewAIDecisionMessageServiceClient(testClientConn)
	testStateClient = protos.NewAIDecisionStateServiceClient(testClientConn)
}

func suiteTeardown() {
	testCtxCancel()
	if testClientConn != nil {
		testClientConn.Close()
	}
}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		suiteSetup()
		defer suiteTeardown()
		return m.Run()
	}())
}
