### Using with Vault

First [bootstrap dev Vault](../tools), then [setup a Vault auto-renew client](../vault).

Add to `docker-compose.yaml`:

```yaml
- VAULT_ENABLE_POSTGRES=true
- VAULT_POSTGRES_CREDS_PATH=dev/postgresql/<cluster name>/creds/<project name>
- POSTGRES_CONN_TMPL=postgres://%s:%s@mpostgres:5432/<project (db) name>
```

Match `<project-name>` with your project's name and `<cluster name>` with value of `VAULT_POSTGRES_NAME` in [Vault bootstrap](../tools).

An example to setup the postgres client:

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/postgres"
  "github.com/callstats-io/go-common/vault"
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

  // get an active postgres db with valid credentials
  db, err := pgc.DB(ctx)
  defer db.Close()
  if err != nil {
    logger.Fatal("Failed to get a valid postgres db", log.Error(err))
  }

  // get the db status
  err := db.Status()
  if err != nil {
    logger.Fatal("postgres db status error", log.Error(err))
  }
}
```

### Using without Vault

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/postgres"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx := context.Background()

  // with auth enabled in postgres
  postgresUri := "postgres://user:pass@localhost:5432/test"
  // without auth
  // postgresUri := "postgres://localhost:5432"

  // create the client
  postgresOpts := &postgres.Options{
    ConnectionTemplate: postgresUri,
  }
  pgc, err := postgres.NewStaticClient(ctx, postgresOpts)
  if err != nil {
    logger.Fatal("Failed to create postgres StaticClient", log.Error(err))
  }

  // get an active postgres db
  db, err := pgc.DB(ctx)
  defer db.Close()
  if err != nil {
    logger.Fatal("Failed to get an open postgres db", log.Error(err))
  }

  // get the db status
  err := db.Status()
  if err != nil {
    logger.Fatal("postgres db status error", log.Error(err))
  }
}
```
