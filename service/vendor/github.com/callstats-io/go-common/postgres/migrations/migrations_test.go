package migrations_test

import (
	"context"
	"fmt"
	"os"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/testutil/pgtestutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createTestMigration(logger log.Logger, opts *migrations.Options) migrations.Migration {
	return migrations.Migration{
		Version: 1,
		Up: func(db migrations.DB) error {
			logger.Info("creating table example...")
			_, err := db.Exec(fmt.Sprintf(`
                SET ROLE '%s';
                CREATE TABLE example(
                    id              SERIAL,
                    name            VARCHAR(100) NOT NULL,
                    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()::timestamp,
                    PRIMARY KEY(id)
                );
                CREATE UNIQUE INDEX example_name_idx ON example (name);
            `, opts.RootRole))

			return err
		},
		Down: func(db migrations.DB) error {
			logger.Warn("dropping table example...")
			_, err := db.Exec(fmt.Sprintf(`
                SET ROLE '%s';
                DROP TABLE IF EXISTS example;
                DROP INDEX IF EXISTS example_email_idx;
            `, opts.RootRole))

			return err
		},
	}
}

var _ = Describe("Migrator", func() {
	var testMigrator *migrations.Migrator
	BeforeEach(func() {
		testMigrator = migrations.NewMigrator()
	})
	Context("Success", func() {
		It("should create migrations table", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				testPostgresClient := newTestPostgresClient(ctx, nil)
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "init", testRoleOptions)).To(BeNil())
				db, err := testPostgresClient.DB(ctx)
				Expect(err).To(BeNil())
				_, err = db.Exec("SELECT id FROM gopg_migrations")
				Expect(err).To(BeNil())
				testPostgresClient.Close()
			})
		})
		It("should be idempotent", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				testPostgresClient := newTestPostgresClient(ctx, nil)
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "init", testRoleOptions)).To(BeNil())
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "init", testRoleOptions)).To(BeNil())
				testPostgresClient.Close()
			})
		})
		It("should run the migrations", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				testMigrator.Register(createTestMigration)
				testPostgresClient := newTestPostgresClient(ctx, nil)
				// create migration table
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "init", testRoleOptions)).To(BeNil())
				// execute migrations
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "up", testRoleOptions)).To(BeNil())

				db, err := testPostgresClient.DB(ctx)
				Expect(err).To(BeNil())
				_, err = db.Exec("SELECT id FROM example")
				Expect(err).To(BeNil())

				// validate down decreases version
				maxVersion := 0
				_, err = db.QueryOne(&maxVersion, "SELECT max(version) FROM gopg_migrations")
				Expect(err).To(BeNil())
				Expect(maxVersion).To(Equal(1))
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "down", testRoleOptions)).To(BeNil())
				_, err = db.QueryOne(&maxVersion, "SELECT max(version) FROM gopg_migrations")
				Expect(err).To(BeNil())
				Expect(maxVersion).To(Equal(1))

				testPostgresClient.Close()
			})
		})
	})
	Context("Error", func() {
		BeforeEach(func() {
			Expect(pgtestutil.CreateTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
		})

		It("should return an error if postgres client returns an error", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				testErr := fmt.Errorf("EXP TEST ERROR")
				testPostgresClient := newTestPostgresClient(ctx, testErr)
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "init", testRoleOptions)).To(MatchError(testErr))
				testPostgresClient.Close()
			})
		})
		It("should return an error on up if init has not been run", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				testMigrator.Register(createTestMigration)
				testPostgresClient := newTestPostgresClient(ctx, nil)
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "up", testRoleOptions)).ToNot(BeNil())
				testPostgresClient.Close()
			})
		})
		It("should return an error if two migrations have the same version", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				// create the same migration twice
				testMigrator.Register(createTestMigration)
				testMigrator.Register(createTestMigration)

				testPostgresClient := newTestPostgresClient(ctx, nil)
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "init", testRoleOptions)).To(BeNil())
				Expect(testMigrator.Migrate(ctx, testPostgresClient, "up", testRoleOptions)).To(MatchError(fmt.Errorf("Conflicting migrations, found two migrations with version %d", 1)))
				testPostgresClient.Close()
			})
		})
	})
})
