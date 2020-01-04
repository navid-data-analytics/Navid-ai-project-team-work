package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/callstats-io/go-common/postgres"
	"github.com/callstats-io/go-common/postgres/migrations"
	"github.com/callstats-io/go-common/testutil/pgtestutil"

	// Import migrations to register them to the global postgres migrater
	service_migrations "github.com/callstats-io/ai-decision/service/migrations"
)

var (
	pgConnectionURL              = os.Getenv("VAULT_POSTGRES_ROOT_URL")
	pgRootRole                   = os.Getenv("POSTGRES_ROOT_ROLE")
	testCtx, testCtxCancel       = context.WithCancel(context.Background())
	testPostgresClient           postgres.Client
	testPostgresClosedConnClient postgres.Client
	testPostgresDB               *postgres.DB
)

func mustBeNil(err error) {
	if err != nil {
		panic(err)
	}
}

func suiteSetup() {
	mustBeNil(pgtestutil.CreateTestPgDb(pgConnectionURL))

	pgc, err := postgres.NewStaticClient(testCtx, &postgres.Options{
		ConnectionTemplate: pgConnectionURL,
	})
	mustBeNil(err)
	testPostgresClient = pgc
	testPostgresDB, err = pgc.DB(testCtx)
	mustBeNil(err)

	// create a closed client to check failures
	pgcc, err := postgres.NewStaticClient(testCtx, &postgres.Options{
		ConnectionTemplate: pgConnectionURL,
	})
	mustBeNil(err)
	db, err := pgcc.DB(testCtx)
	mustBeNil(err)
	db.Close()
	testPostgresClosedConnClient = pgcc

	mustBeNil(migrateTestDB(testCtx, testPostgresClient, pgRootRole, "up"))
}

func suiteTeardown() {
	mustBeNil(migrateTestDB(testCtx, testPostgresClient, pgRootRole, "down"))
	if testPostgresClient != nil {
		testPostgresClient.Close()
	}
	pgtestutil.DropTestPgDb(pgConnectionURL)
	testCtxCancel()
}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		suiteSetup()
		defer suiteTeardown()
		return m.Run()
	}())
}

// migrateTestDB runs migrate init and migrate up on the db pointed to by client
func migrateTestDB(ctx context.Context, client postgres.Client, role string, direction string) error {
	opts := &migrations.Options{
		RootRole: role,
		Meta: map[string]interface{}{
			service_migrations.MetaKeyReadRole: role,
		},
	}
	if err := migrations.Migrate(ctx, client, "init", opts); err != nil {
		return err
	}

	return migrations.Migrate(ctx, client, direction, opts)
}
