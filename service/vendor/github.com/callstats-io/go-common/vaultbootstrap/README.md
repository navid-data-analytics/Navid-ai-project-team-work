### Using in unit tests

`docker-compose.yaml`:

```yaml
- SERVICE_NAME=<project name>
- ENV=test
- VAULT_ADDR=http://vault:8200/
- VAULT_TEST_ROOT_TOKEN_ID=<uuid token>
```

`main_test.go`:

```golang
package main_test

import (
  "testing"

  "github.com/callstats-io/go-common/vaultbootstrap"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var vaultbootstrapClient *vaultbootstrap.BootstrapClient

var _ = BeforeSuite(func() {
		vaultbootstrapClient = vaultbootstrap.NewBootstrapClient().
			WithVaultRootToken(os.Getenv("VAULT_TEST_BOOTSTRAP_TOKEN")).
			WithTestTLSCertData().
			WithTestTLSCertKeyData().
			UnmountAll().
			MountAWS().
			MountTLSCert().
			MountMongo().
			MountPostgres().
			MountAppRoleAuth().
			WriteCredentialsEnv()
})

var _ = AfterSuite(func() {
  vaultbootstrapClient.ClearVault()
})

func TestMain(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Main Suite")
}
```
