package main

import (
	"flag"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/callstats-io/ai-decision/service/src/config"
	"github.com/callstats-io/ai-decision/service/src/flowdock"
	"github.com/callstats-io/ai-decision/service/src/grpc"
	"github.com/callstats-io/ai-decision/service/src/http"
	"github.com/callstats-io/ai-decision/service/src/service"
	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/callstats-io/go-common/app"
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/metrics"
	"github.com/callstats-io/go-common/postgres"
	"github.com/callstats-io/go-common/postgres/migrations"
	raven "github.com/getsentry/raven-go"

	"context"

	// import unnamed to register all migrations
	defined_migrations "github.com/callstats-io/ai-decision/service/migrations"
)

var (
	cmdRunServer = flag.Bool("server", true, "Start server for services")
	cmdMigrate   = flag.String("migrate", "", "Run migrations, value should be a supported command for go-pg/migrations (e.g. init, up, down).")
	cmdDryRun    = flag.Bool("dry-run", false, "Read-only mode")
	deleteList   = flag.String("delete", "", "Comma separated list of notifications to delete")
)

func main() {
	os.Exit(Serve())
}

// Serve starts the server
func Serve() int {
	flag.Parse()
	raven.SetRelease(config.ServiceVersion)
	p, _ := raven.CapturePanicAndWait(func() {
		serviceCtx, serviceCtxCancel := context.WithCancel(context.Background())
		defer serviceCtxCancel()
		Start(serviceCtx)

	}, nil)

	if p != nil {
		switch p.(type) {
		case error:
			log.RootLogger().Error("Panicked", log.Error(p.(error)))
		case string:
			log.RootLogger().Error("Panicked", log.String("error", p.(string)))
		default:
			log.RootLogger().Error("Panicked with unknown reason")
		}
		return 1
	}
	return 0
}

// Start starts a new service instance given configuration from environment
func Start(ctx context.Context) {
	logger := log.FromContext(ctx)

	settings, err := config.FromEnv()
	if err != nil {
		logger.Panic("Error reading the settings from environment", log.Error(err))
	}

	app := app.NewApp(ctx)
	postgresClient, err := postgres.NewStandardClient(app.Context(), app.VaultClient(), &postgres.Options{
		ConnectionTemplate: settings.PostgresConnectionTemplate,
	})
	if err != nil {
		logger.Panic("Failed to create postgres client", log.Error(err))
	}

	if *cmdMigrate != "" {
		logger.Info("Run migrations")
		migrationOptions := &migrations.Options{
			RootRole: settings.PostgresRootRole,
			Meta: map[string]interface{}{
				defined_migrations.MetaKeyReadRole: settings.PostgresReadOnlyRole,
			},
		}
		if err := migrations.Migrate(app.Context(), postgresClient, *cmdMigrate, migrationOptions); err != nil {
			logger.Panic("Failed to run migrations", log.Error(err))
		}
	}

	if *deleteList != "" {
		deletionCtx := app.Context()
		idList := parseStringFlag(string(*deleteList))
		if *cmdDryRun {
			logger.Info("Read-only mode of deletion running, no message will be deleted.")
			_ = storage.CheckExistingMessageIDs(deletionCtx, postgresClient, idList)
			logger.Info("Messages to delete: " + string(idList))
		} else {
			logger.Info("Delete messages with ids: " + string(*deleteList))
			storage.DeleteMessageByID(deletionCtx, postgresClient, idList)
		}
	}

	if *cmdRunServer {
		logger.Info("Run server")

		grpcLn, err := net.Listen("tcp", ":"+strconv.Itoa(settings.GRPCPort))
		if err != nil {
			logger.Panic("Failed to start gRPC listener", log.Int("grpcPort", settings.GRPCPort), log.Error(err))
		}

		storage := storage.NewPostgres(postgresClient)
		flowdockClient := flowdock.NewClient(settings.FlowdockToken)
		messageService, err := service.NewAIDecisionMessageService(storage, flowdockClient)
		if err != nil {
			logger.Panic("Error creating a new ai-decision message service", log.Error(err))
		}

		stateService, err := service.NewAIDecisionStateService(storage)
		if err != nil {
			logger.Panic("Error creating a new ai-decision state service", log.Error(err))
		}

		app.WithHTTPPort(settings.HTTPStatusPort).
			ServeHTTP(http.NewInternalRequestRouter(metrics.PrometheusEndpointWithoutCompression(), postgresStatusCheck(postgresClient)))

		grpcServer, err := grpc.NewServer(ctx, messageService, stateService)
		if err != nil {
			logger.Panic("Error creating a new gRPC server", log.Error(err))
		}

		go func() {
			logger.Info("Starting GRPC server", log.Int("grpcPort", settings.GRPCPort))
			grpcServer.Serve(app.Context(), grpcLn)
		}()

		<-app.Context().Done()
	}
}

func postgresStatusCheck(postgresClient postgres.Client) func(context.Context) error {
	return func(ctx context.Context) error {
		// check connection to postgres works
		db, err := postgresClient.DB(ctx)
		if err != nil {
			return err
		}
		return db.Status()
	}
}

// parseStringFlag converts string flag to int slice
func parseStringFlag(str string) []int32 {
	splitStr := strings.Split(str, ",")
	result := make([]int32, len(splitStr), len(splitStr))
	for i, strNum := range splitStr {
		intNum, err := strconv.Atoi(strNum)
		if err != nil {
			panic(err)
		}
		result[i] = int32(intNum)
	}
	return result
}
