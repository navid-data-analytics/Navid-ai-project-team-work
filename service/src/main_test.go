package main_test

/*
TODO(SH): This test breaks with the other postgres tests because of vault env or something like that.
An obscure error of a database not existing is raised on unmount. As time does not permit further debugging,
this is left for future to be fixed. Likely requires a bit of rethinking in the bootstrapping so that each test
package has it's own isolated environment (database names etc.) which is currently not feasible.

import (
	"context"
	"os"
	"testing"
	"time"

	main "github.com/callstats-io/ai-decision/service/src"
	"github.com/callstats-io/go-common/testutil/pgtestutil"
	"github.com/callstats-io/go-common/vaultbootstrap"
	"github.com/stretchr/testify/require"
)


func TestStart(t *testing.T) {
	assert := require.New(t)
	pgConnectionURL := os.Getenv("VAULT_POSTGRES_ROOT_URL")

	assert.Nil(pgtestutil.CreateTestPgDb(pgConnectionURL))
	defer pgtestutil.DropTestPgDb(pgConnectionURL)

	bc := vaultbootstrap.NewBootstrapClient().
		WithVaultRootToken(os.Getenv("VAULT_TEST_BOOTSTRAP_TOKEN")).
		UnmountAll().
		MountAppRoleAuth().
		MountPostgres().
		WriteCredentialsEnv()

	defer bc.UnmountAll()

	testCtx, testCtxCancel := context.WithTimeout(context.Background(), time.Second)
	defer testCtxCancel()

	main.Start(testCtx)
}
*/
