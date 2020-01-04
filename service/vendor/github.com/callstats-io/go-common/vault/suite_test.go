package vault_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"path/filepath"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/testutil/pgtestutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/callstats-io/go-common/vaultbootstrap"
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testLogBuffer            *testutil.LogBuffer
	testLogger               log.Logger
	testVaultBootstrapConfig *vaultbootstrap.BootstrapConfig
	testVaultBootstrap       *vaultbootstrap.BootstrapClient
	testCertFilePath         string
	testCertKeyFilePath      string
)

func init() {
	var err error
	testCertFilePath, err = filepath.Abs("./testdata/test_cert.pem")
	if err != nil {
		panic(err)
	}
	testCertKeyFilePath, err = filepath.Abs("./testdata/test_cert_key.pem")
	if err != nil {
		panic(err)
	}
}

var _ = BeforeSuite(func() {
	Expect(pgtestutil.CreateTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
	testLogBuffer = testutil.NewLogBuffer()
	testLogger = testLogBuffer.Logger()
	log.SetRootLogger(testLogger)

	testVaultBootstrap = vaultbootstrap.NewBootstrapClient().
		WithVaultRootToken(os.Getenv("VAULT_TEST_BOOTSTRAP_TOKEN")).
		WithTestTLSCertData().
		WithTestTLSCertKeyData().
		WithTestGenericData().
		WithTestTransitData().
		UnmountAll().
		MountAWS().
		MountTLSCert().
		MountMongo().
		MountPostgres().
		MountAppRoleAuth().
		MountGeneric().
		MountTransit().
		WriteCredentialsEnv()
	testVaultBootstrapConfig = testVaultBootstrap.Config()
})

var _ = AfterSuite(func() {
	testVaultBootstrap.UnmountAll()
	Expect(pgtestutil.DropTestPgDb(os.Getenv("VAULT_POSTGRES_ROOT_URL"))).To(BeNil())
})

var _ = BeforeEach(func() {
	testLogBuffer.Reset()
})

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vault Black Box Suite")
}

// ========= TEST UTILS =========

type fakeTestClient struct {
	vaultHTTPClient    *http.Client
	options            *vault.Options
	authCalls          int
	authError          error
	authSecretOnce     *vault.StandardSecret
	authSecret         *vault.StandardSecret
	postgresCalls      int
	postgresError      error
	postgresSecretOnce *vault.UserPassSecret
	postgresSecret     *vault.UserPassSecret
	mongoCalls         int
	mongoError         error
	mongoSecretOnce    *vault.UserPassSecret
	mongoSecret        *vault.UserPassSecret
	tlsCalls           int
	tlsError           error
	tlsSecretOnce      *vault.TLSCertSecret
	tlsSecret          *vault.TLSCertSecret
	awsCalls           int
	awsError           error
	awsSecretOnce      *vault.AWSSecret
	awsSecret          *vault.AWSSecret
	readCalls          int
	readError          error
	readPath           string
	readSecretOnce     *vault.StandardSecret
	readSecret         *vault.StandardSecret
	writeCalls         int
	writeError         error
	writePath          string
	writeData          map[string]interface{}
	writeSecretOnce    *vault.StandardSecret
	writeSecret        *vault.StandardSecret
}

func (ftc *fakeTestClient) Options() *vault.Options {
	return ftc.options
}

func (ftc *fakeTestClient) VaultHTTPClient() *http.Client {
	return ftc.vaultHTTPClient
}

func (ftc *fakeTestClient) Authenticate(ctx context.Context) (*vault.StandardSecret, error) {
	ftc.authCalls++

	if ftc.authSecretOnce != nil {
		s := ftc.authSecretOnce
		ftc.authSecretOnce = nil
		return s, nil
	}

	if ftc.authError != nil {
		return nil, ftc.authError
	}
	return ftc.authSecret, nil
}

func (ftc *fakeTestClient) Read(ctx context.Context, path string) (*vault.StandardSecret, error) {
	ftc.readCalls++
	ftc.readPath = path
	if ftc.readSecretOnce != nil {
		s := ftc.readSecretOnce
		ftc.readSecretOnce = nil
		return s, nil
	}

	if ftc.readError != nil {
		return nil, ftc.readError
	}
	return ftc.readSecret, nil
}

func (ftc *fakeTestClient) Write(ctx context.Context, path string, data map[string]interface{}) (*vault.StandardSecret, error) {
	ftc.writeCalls++
	ftc.writePath = path
	ftc.writeData = data
	if ftc.writeSecretOnce != nil {
		s := ftc.writeSecretOnce
		ftc.writeSecretOnce = nil
		return s, nil
	}

	if ftc.writeError != nil {
		return nil, ftc.writeError
	}
	return ftc.writeSecret, nil
}

func (ftc *fakeTestClient) PostgresSecret(ctx context.Context) (*vault.UserPassSecret, error) {
	ftc.postgresCalls++

	if ftc.postgresSecretOnce != nil {
		s := ftc.postgresSecretOnce
		ftc.postgresSecretOnce = nil
		return s, nil
	}

	if ftc.postgresError != nil {
		return nil, ftc.postgresError
	}
	return ftc.postgresSecret, nil
}

func (ftc *fakeTestClient) MongoSecret(ctx context.Context) (*vault.UserPassSecret, error) {
	ftc.mongoCalls++
	if ftc.mongoSecretOnce != nil {
		s := ftc.mongoSecretOnce
		ftc.mongoSecretOnce = nil
		return s, nil
	}

	if ftc.mongoError != nil {
		return nil, ftc.mongoError
	}
	return ftc.mongoSecret, nil
}

func (ftc *fakeTestClient) TLSCertSecret(ctx context.Context) (*vault.TLSCertSecret, error) {
	ftc.tlsCalls++
	if ftc.tlsSecretOnce != nil {
		s := ftc.tlsSecretOnce
		ftc.tlsSecretOnce = nil
		return s, nil
	}

	if ftc.tlsError != nil {
		return nil, ftc.tlsError
	}
	return ftc.tlsSecret, nil
}

func (ftc *fakeTestClient) AWSSecret(ctx context.Context) (*vault.AWSSecret, error) {
	ftc.awsCalls++
	if ftc.awsSecretOnce != nil {
		s := ftc.awsSecretOnce
		ftc.awsSecretOnce = nil
		return s, nil
	}

	if ftc.awsError != nil {
		return nil, ftc.awsError
	}
	return ftc.awsSecret, nil
}

func isRenewable(secret vault.Secret) bool {
	switch secret.(type) {
	case *vault.UserPassSecret:
		return secret.(*vault.UserPassSecret).Renewable
	case *vault.StandardSecret:
		return secret.(*vault.StandardSecret).Renewable
	case *vault.TLSCertSecret:
		return secret.(*vault.TLSCertSecret).Renewable
	case *vault.AWSSecret:
		return secret.(*vault.AWSSecret).Renewable
	default:
		return false
	}
}

func makeStandardSecret(duration int, data map[string]interface{}) *vault.StandardSecret {
	raw := &api.Secret{
		LeaseDuration: duration,
		Data:          data,
	}
	return vault.NewStandardSecret(raw, nil)
}
