package migrations

import (
	"context"
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

// Log message keys
const (
	LogKeyCommand = "cmd"
)

var (
	defaultMigrator = Migrator{}
)

// Options contains the migration options, role is required
type Options struct {
	RootRole string
	Meta     map[string]interface{}
}

// DB exports the migrations.DB as DB
type DB migrations.DB

// Migration wraps the go-pg migration
type Migration struct {
	Version int64
	Up      func(DB) error
	Down    func(DB) error
}

// MigrationGeneratorFunc is the type all migration register operations need to implement
type MigrationGeneratorFunc func(logger log.Logger, opts *Options) Migration

// Migrator contains logic for running the migrations
type Migrator struct {
	generators []func(logger log.Logger, opts *Options) migrations.Migration
}

// Register registers a migration on this migrator
func (m *Migrator) Register(genFunc MigrationGeneratorFunc) {
	m.generators = append(m.generators, func(logger log.Logger, opts *Options) migrations.Migration {
		mig := genFunc(logger, opts)

		// wrap the resulting migration as go-pg migrations migration
		return migrations.Migration{
			Version: mig.Version,
			Up: func(db migrations.DB) error {
				return mig.Up(db)
			},
			Down: func(db migrations.DB) error {
				return mig.Down(db)
			},
		}
	})
}

// Migrate runs all registered migrations on this migrator
func (m *Migrator) Migrate(ctx context.Context, client postgres.Client, cmd string, opts *Options) error {
	logger := log.FromContextWithPackageName(ctx, "go-common/postgres/migrations").With(log.String(LogKeyCommand, cmd))
	logger.Info("Run migrations")

	migrationCtx, migrationCtxCancel := context.WithCancel(ctx)
	defer migrationCtxCancel()

	db, err := client.DB(migrationCtx)
	if err != nil {
		return err
	}

	switch cmd {
	case "init":
		logger.Info("Create migrations tracking table")
		query := fmt.Sprintf(`
                SET ROLE '%s';
                CREATE TABLE IF NOT EXISTS gopg_migrations (id serial, version bigint, created_at timestamptz);
            `, opts.RootRole)
		if _, err := db.Exec(query); err != nil {
			return err
		}
		return migrateInTransaction(db, nil, "version")

	case "up", "down":
		// ensure migrations table has been created, no longer done by go-pg migrations itself and our table is different
		if _, err := db.Exec("select 1 as one from gopg_migrations"); err != nil {
			return err
		}

		migs := make([]migrations.Migration, len(m.generators))
		versionSeen := make(map[int64]bool, len(m.generators))
		for idx, generator := range m.generators {
			// get nested go-pg migration from the wrapped migration
			mig := generator(logger, opts)
			if versionSeen[mig.Version] {
				return fmt.Errorf("Conflicting migrations, found two migrations with version %d", mig.Version)
			}

			migs[idx] = mig
			versionSeen[mig.Version] = true
		}
		return migrateInTransaction(db, migs, cmd)

	default:
		return fmt.Errorf("Unmapped command %s", cmd)
	}
}

// NewMigrator returns a new empty Migrator
func NewMigrator() *Migrator {
	return &Migrator{}
}

// Register registers a new migration generator callback on the default migrator
func Register(genFunc MigrationGeneratorFunc) {
	defaultMigrator.Register(genFunc)
}

// Migrate runs the migrations registered on the default migrator
func Migrate(ctx context.Context, client postgres.Client, cmd string, opts *Options) error {
	return defaultMigrator.Migrate(ctx, client, cmd, opts)
}

func migrateInTransaction(db *postgres.DB, migs []migrations.Migration, cmd string) error {
	return db.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, _, err := migrations.RunMigrations(tx, migs, cmd)
		return err
	})
}
