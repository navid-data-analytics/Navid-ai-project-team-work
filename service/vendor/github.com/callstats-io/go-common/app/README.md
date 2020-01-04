### Using with Vault

First [bootstrap dev Vault](../tools), then [setup a Vault auto-renew client](../vault).

Add to `docker-compose.yaml`:

```yaml
- VAULT_ENABLE_TLS_CERT=true
- VAULT_TLS_CERT_PATH=dev/generic/<project name>/certificates/x509/cert
- VAULT_TLS_CERT_KEY_PATH=dev/generic/<project name>/certificates/x509/key
```

Match `<project-name>` with your project's name in [Vault bootstrap](../tools).


### Serving over HTTPS (with HTTP/2 support)
```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/app"
  "github.com/callstats-io/go-common/vault"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx, ctxCancel := context.WithCancel(context.Background())
  defer ctxCancel()

  mux := httptreemux.New()
  // ... Setup routes ...
  
  // SIMPLE: initalizes vault.AutoRenewClient from env, fetches TLS cert secret and reads HTTPS_PORT from env
  a := app.New(ctx).ServeHTTPS(mux)

  // COMPLEX: initalize everything yourself
  a := app.New(ctx).
    WithVaultClient(yourVaultClient).
    WithHTTPSPort(yourHttpsPort).
    ServeHTTPS(mux)

  // blocks until the server shuts down
  a.Done()
}
```

### Serving over HTTP (no HTTP/2 support)

```golang
package main

import (
  "context"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/app"
  "github.com/callstats-io/go-common/vault"
)

func main() {
  logger := log.NewLogger("DEBUG")
  ctx, ctxCancel := context.WithCancel(context.Background())
  defer ctxCancel()

  mux := httptreemux.New()
  // ... Setup routes ...
  
  // SIMPLE: reads HTTP_PORT from env
  a := app.New(ctx).ServeHTTP(mux)

  // COMPLEX: initalize port yourself
  a := app.New(ctx).
    WithHTTPPort(yourHttpPort).
    ServeHTTP(mux)
}
```
