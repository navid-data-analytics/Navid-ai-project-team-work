### Usage

First follow the [vault bootstrap instructions](../tools).

Next add these to `docker-compose.yaml` (replace `<project name>`):

```yaml
- VAULT_ADDR=http://vault:8200/
- VAULT_AUTHCREDENTIALS=file:/go/src/github.com/callstats-io/<project name>/creds/vault_dev_creds.json
```

An example to import and setup the vault auto-renew client:

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/vault"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx := context.Background()

  vaultOpts, err := vault.OptionsFromEnv()
  if err != nil {
    logger.Fatal("Failed to read vault options", log.Error(err))
  }
  vc, err := vault.NewStandardClient(ctx, vaultOpts)
  if err != nil {
    logger.Fatal("Failed to create vault StandardClient", log.Error(err))
  }
  avc, err := vault.NewAutoRenewClient(ctx, vc)
  if err != nil {
    logger.Fatal("Failed to create vault AutoRenewClient", log.Error(err))
  }

  // Prefer to use the Read-interface as the other methods will be removed eventually
  secret, err := avc.Read(ctx, "/your/secret/vault/path")
  if err != nil {
    logger.Fatal("Failed to read vault secret", log.Error(err))
  }
  // utilize avc the way you need
}
```

See also: [Bootstrapping Vault for unit tests](../vaultbootstrap)
