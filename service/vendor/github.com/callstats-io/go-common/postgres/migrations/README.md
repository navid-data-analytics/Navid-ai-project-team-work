### Using with Vault

First [get a postgres client](../).

Match `<YOUR_POSTGRES_ROLE>` with your project's postgres role name (generally same as the service name) and `<YOUR_MIGRATION_COMMAND>` with one of the go-pg migrations supported commands. Most common are `init`, `up`, `down`, where `init` has to be run first to create the migration tracking table.


####To run migrations in your project you need to:

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/vault"
  "github.com/callstats-io/go-common/postgres"
  "github.com/callstats-io/go-common/postgres/migrations"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx := context.Background()

  // set up vault auto-renew client, details in Vault README
  var avc *vault.AutoRenewClient
  avc = setupAvc()

  // read postgres options from environment variables
  postgresOpts := postgres.OptionsFromEnv()
  pgc, err := postgres.NewStandardClient(avc, postgresOpts)
  if err != nil {
    logger.Fatal("Failed to create postgres StandardClient", log.Error(err))
  }

  if err := migrations.Migrate(ctx, pgc, "YOUR_POSTGRES_ROLE", "YOUR_MIGRATION_COMMAND"); err != nil {
    logger.Fatal("Failed to create postgres StandardClient", log.Error(err))
  }

```


#### To create migrations:

```golang

package main

import (
	"fmt"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres/migrations"
)

func init() {
	migrations.Register(func(logger log.Logger, role string) migrations.Migration {
		return migrations.Migration{
			Version: 1, // HAS TO BE INCREMENTALLY DIFFERENT FOR EACH MIGRATION
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
					`, role))

				return err
			},
			Down: func(db migrations.DB) error {
				logger.Warn("dropping table example...")
				_, err := db.Exec(fmt.Sprintf(`
					SET ROLE '%s';
					DROP TABLE IF EXISTS example;
					DROP INDEX IF EXISTS example_email_idx;
				`, role))

				return err
			},
		}
	})
}
```