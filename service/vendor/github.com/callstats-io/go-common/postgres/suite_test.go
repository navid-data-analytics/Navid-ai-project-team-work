package postgres_test

import (
	"os"
	"testing"

	"github.com/callstats-io/go-common/testutil/pgtestutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	Expect(pgtestutil.CreateTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
})

var _ = AfterSuite(func() {
	Expect(pgtestutil.DropTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
})

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Black Box Suite")
}
