package pgtestutil

import (
	"fmt"

	"github.com/callstats-io/go-common/postgres"

	"github.com/go-pg/pg"
)

func runWithoutDatabase(opts *pg.Options, cb func(*pg.DB) error) error {
	pgClientWithoutDb := pg.Connect(&pg.Options{
		Addr:     opts.Addr,
		User:     opts.User,
		Password: opts.Password,
	})
	defer pgClientWithoutDb.Close()
	return cb(pgClientWithoutDb)
}

// CreateTestPgDb recreates the postgres DB based on environment variables
func CreateTestPgDb(configURL string) error {
	opts, err := postgres.ParseURL(configURL)
	if err != nil {
		return err
	}

	return runWithoutDatabase(opts, func(pgClient *pg.DB) error {
		if err := closeConnectionsToTestDB(pgClient, opts); err != nil {
			return err
		}

		//NOTE SH: NEVER do this with user provided info (injection vulnerability)
		_, err = pgClient.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", opts.Database))

		if err != nil {
			return err
		}

		//NOTE SH: NEVER do this with user provided info (injection vulnerability)
		_, err = pgClient.Exec(fmt.Sprintf("CREATE DATABASE %v OWNER %v", opts.Database, opts.User))

		if err != nil {
			return err
		}
		return nil
	})
}

// DropTestPgDb drops the test database
func DropTestPgDb(configURL string) error {
	opts, err := postgres.ParseURL(configURL)
	if err != nil {
		return err
	}

	return runWithoutDatabase(opts, func(pgClient *pg.DB) error {
		if err := closeConnectionsToTestDB(pgClient, opts); err != nil {
			return err
		}
		// drop database
		_, err = pgClient.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", opts.Database))

		if err != nil {
			return err
		}
		return nil
	})
}

func closeConnectionsToTestDB(pgDB *pg.DB, opts *pg.Options) error {
	//NOTE SH: NEVER inline params with user provided info in production scenarios (injection vulnerability)
	// kill all connections
	_, err := pgDB.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%v'
		AND pid <> pg_backend_pid();
	`, opts.Database))

	return err
}
