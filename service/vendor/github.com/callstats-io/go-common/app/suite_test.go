package app_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/callstats-io/go-common/response"
	"github.com/callstats-io/go-common/vault"
	"github.com/callstats-io/go-common/vaultbootstrap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testVaultClient    vault.Client
	testVaultBootstrap *vaultbootstrap.BootstrapClient
	testClient         *http.Client

	testCtx, testCtxCancel = context.WithCancel(context.Background())
)

func init() {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ExpectContinueTimeout: 0,
	}
	// Create new client with configured transport
	testClient = &http.Client{
		Transport: tr,
		Timeout:   0,
	}
}

var _ = BeforeSuite(func() {
	testVaultBootstrap = vaultbootstrap.NewBootstrapClient().
		WithVaultRootToken(os.Getenv("VAULT_TEST_BOOTSTRAP_TOKEN")).
		WithTestTLSCertData().
		WithTestTLSCertKeyData().
		UnmountAll().
		MountAppRoleAuth().
		MountTLSCert().
		WriteCredentialsEnv()
	vaultOpts, err := vault.OptionsFromEnv()
	Expect(err).To(BeNil())
	vaultOpts.EnableMongo = false    // not needed
	vaultOpts.EnablePostgres = false // not needed
	vc, err := vault.NewStandardClient(testCtx, vaultOpts)
	Expect(err).To(BeNil())
	testVaultClient, err = vault.NewAutoRenewClient(testCtx, vc)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	testCtxCancel()
	testVaultBootstrap.UnmountAll()
})

// ===== TEST SETUP =====
func TestAll(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Black Box Suite")
}

func validateResponse(resp *http.Response, err error) {
	Expect(err).To(BeNil())
	defer resp.Body.Close()
	Expect(resp.StatusCode).To(Equal(200))
	data, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	Expect(data).To(Equal(response.RespOK))
}

func validateShutdown(resp *http.Response, err error) {
	// ensure the app was shut down, sleep added as sometimes the app is pending shutdown for a while
	time.Sleep(10 * time.Millisecond)
	Expect(resp).To(BeNil())
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(ContainSubstring(": connection refused"))
}

type fakeVaultClient struct {
	vault.Client
	tlsCallCount int
	tlsError     error
}

func (f *fakeVaultClient) TLSCertSecret(ctx context.Context) (*vault.TLSCertSecret, error) {
	f.tlsCallCount++
	if f.tlsError != nil {
		return nil, f.tlsError
	}

	return f.Client.TLSCertSecret(ctx)
}
