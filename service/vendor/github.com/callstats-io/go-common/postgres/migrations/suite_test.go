package migrations_test

import (
	"context"
	"os"
	"testing"

	"github.com/callstats-io/go-common/postgres"
	"github.com/callstats-io/go-common/postgres/migrations"
	"github.com/callstats-io/go-common/testutil/pgtestutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testRoleOptions = &migrations.Options{RootRole: os.Getenv("SERVICE_NAME")}
)

var _ = BeforeSuite(func() {
	Expect(pgtestutil.CreateTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
})

var _ = AfterSuite(func() {
	Expect(pgtestutil.DropTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
})

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Migrations Black Box Suite")
}

type fakeVaultPostgresClient struct {
	postgresSecretError error
	postgresSecret      *vault.UserPassSecret
}

func (fmc *fakeVaultPostgresClient) PostgresSecret(ctx context.Context) (*vault.UserPassSecret, error) {
	if fmc.postgresSecretError != nil {
		return nil, fmc.postgresSecretError
	}
	return fmc.postgresSecret, nil
}

func newTestPostgresClient(ctx context.Context, vaultErr error) postgres.Client {
	testVaultClient := &fakeVaultPostgresClient{}
	if vaultErr == nil {
		secret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
			LeaseDuration: 60,
			Data: map[string]interface{}{
				"username": "go_common",
				"password": "test",
			},
		}, nil))
		Expect(err).To(BeNil())
		testVaultClient.postgresSecret = secret
	} else {
		testVaultClient.postgresSecretError = vaultErr
	}

	opts, err := postgres.OptionsFromEnv()
	Expect(err).To(BeNil())

	client, err := postgres.NewStandardClient(ctx, testVaultClient, opts)
	Expect(err).To(BeNil())
	return client
}
