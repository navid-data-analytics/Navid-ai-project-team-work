### Using with Vault

First [bootstrap dev Vault](../tools), then [setup a Vault auto-renew client](../vault).

Add to `docker-compose.yaml`:

```yaml
- VAULT_ENABLE_MONGO=true
- VAULT_MONGO_CREDS_PATH=dev/mongodb/<cluster name>/creds/<project name>
- MONGO_CONN_TMPL=mongodb://%s:%s@mongo:27017/admin
```

Match `<project-name>` with your project's name and `<cluster name>` with value of `VAULT_MONGO_CLUSTER_NAME` in [Vault bootstrap](../tools).

An example to setup the MongoDB client:

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/mongo"
  "github.com/callstats-io/go-common/vault"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx := context.Background()

  // set up vault auto-renew client, details in Vault README
  var avc *vault.AutoRenewClient
  avc = setupAvc()

  // read mongo options from environment variables
  mongoOpts := mongo.OptionsFromEnv()
  mc, err := mongo.NewStandardClient(avc, mongoOpts)
  if err != nil {
    logger.Fatal("Failed to create mongo StandardClient", log.Error(err))
  }

  // get an active mongo session with valid credentials
  session, err := mc.Session(ctx)
  defer session.Close()
  if err != nil {
    logger.Fatal("Failed to get a valid mongo session", log.Error(err))
  }

  // get the session status
  status, err := session.Status()
  if err != nil {
    logger.Fatal("Mongo session status error", log.Error(err))
  } else {
    connections := mongo.M(status)["connections"].(mongo.M)["current"].(int)
    logger.Debug("Mongo session ok!", log.Int("open connections", connections))
  }
}
```

### Using without Vault

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/mongo"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx := context.Background()

  // with auth enabled in mongo
  mongoUri := "mongodb://user:pass@localhost:27017/admin"
  // without auth
  // mongoUri := "mongodb://localhost:27017"

  // create the client
  mongoOpts := &mongo.Options{
    ConnectionTemplate: mongoUri,
  }
  mc, err := mongo.NewStaticClient(ctx, mongoOpts)
  if err != nil {
    logger.Fatal("Failed to create mongo StaticClient", log.Error(err))
  }

  // get an active mongo session
  session, err := mc.Session(ctx)
  defer session.Close()
  if err != nil {
    logger.Fatal("Failed to get an open mongo session", log.Error(err))
  }

  // get the session status
  status, err := session.Status()
  if err != nil {
    logger.Fatal("Mongo session status error", log.Error(err))
  } else {
    connections := mongo.M(status)["connections"].(mongo.M)["current"].(int)
    logger.Debug("Mongo session ok!", log.Int("open connections", connections))
  }
}
```
