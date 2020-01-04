package vault_test

import (
	"context"
	"errors"
	"time"

	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Each autoRenewTestCase contains the methods for testing a simple auto renew handler based proxy.
type autoRenewTestCase struct {
	Name string

	// Vault options to use for creating the test client
	VaultOptions *vault.Options

	// Errors
	ExpDisabledError error

	// SetupAuthError sets up the fakeTestClient to return an error on authenticate call.
	// This method is expected to return the error to match against.
	SetupAuthError func(*fakeTestClient) error

	// SetupAuthValidOnce sets up the fakeTestClient to return a valid auth secret once,
	// after which the result of either SetupError/SetupValid is returned.
	// It is expected to return the secret to match against on the tested method.
	SetupAuthValidOnce func(*fakeTestClient) *vault.StandardSecret

	// SetupAuthValid sets up the fakeTestClient to return a valid auth secret
	// It is expected to return the secret to match against on authenticate-call
	SetupAuthValid func(*fakeTestClient) *vault.StandardSecret

	// SetupError sets up the fakeTestClient to return an error on the tested method.
	// This method is expected to return the error to match against.
	SetupError func(*fakeTestClient) error

	// SetupValidOnce sets up the fakeTestClient to return a valid secret once,
	// after which the result of either SetupError/SetupValid is returned.
	// It is expected to return the secret to match against on the tested method.
	SetupValidOnce func(*fakeTestClient, *vault.StandardSecret, time.Duration) vault.Secret

	// SetupValid sets up the fakeTestClient to return a valid secret.
	// It is expected to return the secret to match against on the tested method.
	SetupValid func(*fakeTestClient, *vault.StandardSecret) vault.Secret

	// Exec is called when the actual method to test should be called on the AutoRenewClient.
	// It is expected to return the result of that call (e.g. secret + error).
	Exec func(*vault.AutoRenewClient, context.Context) (vault.Secret, error)

	// Calls should return the number of auth calls made and the number of calls made to the tested method
	Calls func(*fakeTestClient) (int, int)
}

var _ = Describe("AutoRenewClient", func() {
	Describe("Authenticate", func() {
		var testClient *fakeTestClient

		BeforeEach(func() {
			testClient = &fakeTestClient{
				options: &vault.Options{},
			}
		})
		Context("Success", func() {
			BeforeEach(func() {
				testClient.authSecret = vault.NewStandardSecret(&api.Secret{
					LeaseDuration: 300, // 5 minutes
				}, nil)
			})
			It("should authenticate on new client create", func() {
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					_, err := vault.NewAutoRenewClient(ctx, testClient)
					Expect(err).To(BeNil())
					Expect(testClient.authCalls).To(Equal(1))
				})
				Expect(err).To(BeNil())
			})
			It("should not authenticate again if the auth secret is valid", func() {
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
					Expect(err).To(BeNil())
					_, err = testAutoRenewClient.Authenticate(ctx)
					Expect(err).To(BeNil())
					Expect(testClient.authCalls).To(Equal(1))
				})
				Expect(err).To(BeNil())
			})
		})

		Context("Failure", func() {
			It("should return an error if the underlying client returns an error on the first time", func() {
				expAuthErr := errors.New("AUTHERR")
				testClient.authError = expAuthErr

				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					_, err := vault.NewAutoRenewClient(ctx, testClient)
					Expect(err).To(MatchError(expAuthErr))
				})
				Expect(err).To(BeNil())
			})
			It("should return an error if the underlying client returns error on refresh", func() {
				authSecret := vault.NewStandardSecret(&api.Secret{
					LeaseDuration: 1, // 1 second
				}, nil)
				testClient.authSecretOnce = authSecret
				expAuthErr := errors.New("AUTHERR")
				testClient.authError = expAuthErr

				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
					Expect(err).To(BeNil())

					select {
					case <-authSecret.RenewContext().Done():
						// wait to give time for context change to propagate
						time.Sleep(10 * time.Millisecond)
					case <-ctx.Done():
						Fail("Timed out")
					}

					_, err = testAutoRenewClient.Authenticate(ctx)
					Expect(err).To(MatchError(expAuthErr))
				})
				Expect(err).To(BeNil())
			})
		})
	})

	setupValidAuthSecret := func(testClient *fakeTestClient) *vault.StandardSecret {
		testClient.authSecret = vault.NewStandardSecret(&api.Secret{
			LeaseDuration: 300, // 5 minutes
		}, nil)
		return testClient.authSecret
	}
	setupAuthValidSecretOnce := func(testClient *fakeTestClient) *vault.StandardSecret {
		testClient.authSecretOnce = vault.NewStandardSecret(&api.Secret{
			LeaseDuration: 5, // 5 seconds
		}, nil)
		return testClient.authSecretOnce
	}
	setupAuthError := func(testClient *fakeTestClient) error {
		testClient.authError = errors.New("EXPAUTHERR")
		return testClient.authError
	}
	newUserPassSecret := func(authSecret *vault.StandardSecret, leaseDuration time.Duration) (*vault.UserPassSecret, error) {
		return vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
			LeaseDuration: int(leaseDuration / time.Second),
			Data: map[string]interface{}{
				vault.SecretDataKeyUsername: "vault",
				vault.SecretDataKeyPassword: "vault",
			},
		}, authSecret))
	}
	newTLSCertSecret := func(authSecret *vault.StandardSecret, leaseDuration time.Duration) (*vault.TLSCertSecret, error) {
		certSecret := vault.NewStandardSecret(&api.Secret{
			LeaseDuration: int(leaseDuration / time.Second),
			Data: map[string]interface{}{
				vault.SecretDataKeyData: testVaultBootstrapConfig.TLSCertData,
			},
		}, authSecret)
		keySecret := vault.NewStandardSecret(&api.Secret{
			LeaseDuration: int(leaseDuration / time.Second),
			Data: map[string]interface{}{
				vault.SecretDataKeyData: testVaultBootstrapConfig.TLSCertKeyData,
			},
		}, authSecret)
		return vault.NewTLSCertSecret(certSecret, keySecret)
	}
	newAWSSecret := func(authSecret *vault.StandardSecret, leaseDuration time.Duration) (*vault.AWSSecret, error) {
		return vault.NewAWSSecret(vault.NewStandardSecret(&api.Secret{
			LeaseDuration: int(leaseDuration / time.Second),
			Data: map[string]interface{}{
				vault.SecretDataKeyAccessKey:     "vault_acc_key",
				vault.SecretDataKeySecretKey:     "vault_sec_key",
				vault.SecretDataKeySecurityToken: "vault_sec_token",
			},
		}, authSecret))
	}
	newStandardSecret := func(authSecret *vault.StandardSecret, leaseDuration time.Duration) (*vault.StandardSecret, error) {
		return vault.NewStandardSecret(&api.Secret{
			LeaseDuration: int(leaseDuration / time.Second),
			Data: map[string]interface{}{
				vault.SecretDataKeyUsername: "vault",
				vault.SecretDataKeyPassword: "vault",
			},
		}, authSecret), nil
	}
	// Test all auto renew methods
	testCases := []autoRenewTestCase{
		autoRenewTestCase{
			Name:               "MongoSecret",
			VaultOptions:       &vault.Options{EnableMongo: true},
			ExpDisabledError:   vault.ErrMongoDisabled,
			SetupAuthError:     setupAuthError,
			SetupAuthValidOnce: setupAuthValidSecretOnce,
			SetupAuthValid:     setupValidAuthSecret,
			SetupError: func(testClient *fakeTestClient) error {
				testClient.mongoError = errors.New("EXPMONGOERR")
				return testClient.mongoError
			},
			SetupValidOnce: func(testClient *fakeTestClient, authSecret *vault.StandardSecret, lease time.Duration) vault.Secret {
				testClient.mongoSecretOnce, _ = newUserPassSecret(authSecret, lease)
				return testClient.mongoSecretOnce
			},
			SetupValid: func(testClient *fakeTestClient, authSecret *vault.StandardSecret) vault.Secret {
				testClient.mongoSecret, _ = newUserPassSecret(authSecret, 10*time.Minute)
				return testClient.mongoSecret
			},
			Exec: func(arc *vault.AutoRenewClient, ctx context.Context) (vault.Secret, error) {
				return arc.MongoSecret(ctx)
			},
			Calls: func(testClient *fakeTestClient) (int, int) {
				return testClient.authCalls, testClient.mongoCalls
			},
		},
		autoRenewTestCase{
			Name:               "PostgresSecret",
			VaultOptions:       &vault.Options{EnablePostgres: true},
			ExpDisabledError:   vault.ErrPostgresDisabled,
			SetupAuthError:     setupAuthError,
			SetupAuthValidOnce: setupAuthValidSecretOnce,
			SetupAuthValid:     setupValidAuthSecret,
			SetupError: func(testClient *fakeTestClient) error {
				testClient.postgresError = errors.New("EXPPOSTGRESERR")
				return testClient.postgresError
			},
			SetupValidOnce: func(testClient *fakeTestClient, authSecret *vault.StandardSecret, lease time.Duration) vault.Secret {
				testClient.postgresSecretOnce, _ = newUserPassSecret(authSecret, lease)
				return testClient.postgresSecretOnce
			},
			SetupValid: func(testClient *fakeTestClient, authSecret *vault.StandardSecret) vault.Secret {
				testClient.postgresSecret, _ = newUserPassSecret(authSecret, 10*time.Minute)
				return testClient.postgresSecret
			},
			Exec: func(arc *vault.AutoRenewClient, ctx context.Context) (vault.Secret, error) {
				return arc.PostgresSecret(ctx)
			},
			Calls: func(testClient *fakeTestClient) (int, int) {
				return testClient.authCalls, testClient.postgresCalls
			},
		},
		autoRenewTestCase{
			Name:               "TLSCertSecret",
			VaultOptions:       &vault.Options{EnableTLSCert: true},
			ExpDisabledError:   vault.ErrTLSCertDisabled,
			SetupAuthError:     setupAuthError,
			SetupAuthValidOnce: setupAuthValidSecretOnce,
			SetupAuthValid:     setupValidAuthSecret,
			SetupError: func(testClient *fakeTestClient) error {
				testClient.tlsError = errors.New("EXPTLSERR")
				return testClient.tlsError
			},
			SetupValidOnce: func(testClient *fakeTestClient, authSecret *vault.StandardSecret, lease time.Duration) vault.Secret {
				s, err := newTLSCertSecret(authSecret, lease)
				Expect(err).To(BeNil())
				testClient.tlsSecretOnce = s
				return testClient.tlsSecretOnce
			},
			SetupValid: func(testClient *fakeTestClient, authSecret *vault.StandardSecret) vault.Secret {
				s, err := newTLSCertSecret(authSecret, 10*time.Minute)
				Expect(err).To(BeNil())
				testClient.tlsSecret = s
				return testClient.tlsSecret
			},
			Exec: func(arc *vault.AutoRenewClient, ctx context.Context) (vault.Secret, error) {
				return arc.TLSCertSecret(ctx)
			},
			Calls: func(testClient *fakeTestClient) (int, int) {
				return testClient.authCalls, testClient.tlsCalls
			},
		},
		autoRenewTestCase{
			Name:               "AWSSecret",
			VaultOptions:       &vault.Options{EnableAWS: true},
			ExpDisabledError:   vault.ErrAWSDisabled,
			SetupAuthError:     setupAuthError,
			SetupAuthValidOnce: setupAuthValidSecretOnce,
			SetupAuthValid:     setupValidAuthSecret,
			SetupError: func(testClient *fakeTestClient) error {
				testClient.awsError = errors.New("EXPAWSERR")
				return testClient.awsError
			},
			SetupValidOnce: func(testClient *fakeTestClient, authSecret *vault.StandardSecret, lease time.Duration) vault.Secret {
				s, err := newAWSSecret(authSecret, lease)
				Expect(err).To(BeNil())
				testClient.awsSecretOnce = s
				return testClient.awsSecretOnce
			},
			SetupValid: func(testClient *fakeTestClient, authSecret *vault.StandardSecret) vault.Secret {
				s, err := newAWSSecret(authSecret, 10*time.Minute)
				Expect(err).To(BeNil())
				testClient.awsSecret = s
				return testClient.awsSecret
			},
			Exec: func(arc *vault.AutoRenewClient, ctx context.Context) (vault.Secret, error) {
				return arc.AWSSecret(ctx)
			},
			Calls: func(testClient *fakeTestClient) (int, int) {
				return testClient.authCalls, testClient.awsCalls
			},
		},
		// Duplicate AWS secret test case as "Read" test case
		autoRenewTestCase{
			Name:               "Read",
			VaultOptions:       &vault.Options{},
			SetupAuthError:     setupAuthError,
			SetupAuthValidOnce: setupAuthValidSecretOnce,
			SetupAuthValid:     setupValidAuthSecret,
			SetupError: func(testClient *fakeTestClient) error {
				testClient.readError = errors.New("EXPREADERR")
				return testClient.readError
			},
			SetupValidOnce: func(testClient *fakeTestClient, authSecret *vault.StandardSecret, lease time.Duration) vault.Secret {
				s, err := newStandardSecret(authSecret, lease)
				Expect(err).To(BeNil())
				testClient.readSecretOnce = s
				return testClient.readSecretOnce
			},
			SetupValid: func(testClient *fakeTestClient, authSecret *vault.StandardSecret) vault.Secret {
				s, err := newStandardSecret(authSecret, 10*time.Minute)
				Expect(err).To(BeNil())
				testClient.readSecret = s
				return testClient.readSecret
			},
			Exec: func(arc *vault.AutoRenewClient, ctx context.Context) (vault.Secret, error) {
				return arc.Read(ctx, arc.Options().AWSCredsReadPath)
			},
			Calls: func(testClient *fakeTestClient) (int, int) {
				return testClient.authCalls, testClient.readCalls
			},
		},
	}

	for idx := range testCases {
		testCase := &testCases[idx]

		Describe(testCase.Name, func() {
			var testClient *fakeTestClient

			BeforeEach(func() {
				testClient = &fakeTestClient{
					options: testCase.VaultOptions,
				}
			})

			Context("Success", func() {
				It("should return a secret if fetch is successful", func() {
					testCase.SetupAuthValid(testClient)
					expSecret := testCase.SetupValid(testClient, testClient.authSecret)

					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
						Expect(err).To(BeNil())
						secret, err := testCase.Exec(testAutoRenewClient, ctx)
						Expect(err).To(BeNil())
						Expect(secret).To(Equal(expSecret))
					})
					Expect(err).To(BeNil())

					// sanity check for no new calls
					authCalls, methodCalls := testCase.Calls(testClient)
					Expect(authCalls).To(Equal(1))
					Expect(methodCalls).To(Equal(1))
				})

				It("should renew the secret automatically on lease expire", func() {
					testCase.SetupAuthValid(testClient)
					expOnceSecret := testCase.SetupValidOnce(testClient, testClient.authSecret, time.Second)
					expSecret := testCase.SetupValid(testClient, testClient.authSecret)

					err := testutil.WithDeadlineContext(5*time.Second, func(ctx context.Context) {
						testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
						Expect(err).To(BeNil())
						_, err = testCase.Exec(testAutoRenewClient, ctx)
						Expect(err).To(BeNil())

						select {
						case <-expOnceSecret.RenewContext().Done():
							// wait to give time for context change to propagate
							time.Sleep(10 * time.Millisecond)
						case <-ctx.Done():
							Fail("Timed out")
						}

						// check that the renew happened and verify call counts
						authCalls, methodCalls := testCase.Calls(testClient)
						Expect(authCalls).To(Equal(1), "expected 1 auth call")
						Expect(methodCalls).To(Equal(2), "expected 2 method calls")

						// expect the secret to now be different
						secret2, err := testCase.Exec(testAutoRenewClient, ctx)
						Expect(err).To(BeNil())
						Expect(secret2).To(Equal(expSecret))
						Expect(secret2).ToNot(Equal(expOnceSecret))
					})
					Expect(err).To(BeNil())

					// sanity check for no new calls
					authCalls, methodCalls := testCase.Calls(testClient)
					Expect(authCalls).To(Equal(1), "expected 1 auth call")
					Expect(methodCalls).To(Equal(2), "expected 2 method calls")
				})

				It("should renew when authentication expires", func() {
					authRenewSecret := testCase.SetupAuthValidOnce(testClient)
					testCase.SetupAuthValid(testClient)
					testCase.SetupValidOnce(testClient, authRenewSecret, 10*time.Minute)
					testCase.SetupValid(testClient, testClient.authSecret)

					err := testutil.WithDeadlineContext(5*time.Second, func(ctx context.Context) {
						testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
						Expect(err).To(BeNil())
						_, err = testCase.Exec(testAutoRenewClient, ctx)
						Expect(err).To(BeNil())

						select {
						case <-authRenewSecret.ExpireContext().Done():
							// wait to give time for context change to propagate
							time.Sleep(10 * time.Millisecond)
						case <-ctx.Done():
							Fail("Timed out")
						}

						authCalls, methodCalls := testCase.Calls(testClient)
						Expect(authCalls).To(Equal(2), "expected 2 auth calls")
						Expect(methodCalls).To(Equal(2), "expected 2 method calls")
					})
					Expect(err).To(BeNil())
				})

				It("should retry renew with backoff if it fails", func() {
					testCase.SetupAuthValid(testClient)
					expOnceSecret := testCase.SetupValidOnce(testClient, testClient.authSecret, time.Second)
					_ = testCase.SetupError(testClient)

					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
						Expect(err).To(BeNil())
						secret, err := testCase.Exec(testAutoRenewClient, ctx)
						Expect(err).To(BeNil())
						Expect(secret).To(Equal(expOnceSecret))

						Eventually(func() int {
							authCalls, methodCalls := testCase.Calls(testClient)
							Expect(authCalls).To(Equal(1))
							return methodCalls
						}).Should(BeNumerically(">", 2))
					})
					Expect(err).To(BeNil())
				})
			})

			Context("Failure", func() {
				// Read secrets are not verified during create so we need to skip these tests
				if testCase.Name != "Read" {
					It("should return an error if the method is not enabled", func() {
						disabledClient := &fakeTestClient{
							options: &vault.Options{},
						}
						testCase.SetupAuthValid(disabledClient)
						err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
							testAutoRenewClient, _ := vault.NewAutoRenewClient(ctx, disabledClient)
							_, err := testCase.Exec(testAutoRenewClient, ctx)
							Expect(err).To(MatchError(testCase.ExpDisabledError))
						})
						Expect(err).To(BeNil())
					})

					It("should return an error at create if the underlying client returns an error on authenticate", func() {
						expAuthErr := testCase.SetupAuthError(testClient)

						err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
							_, err := vault.NewAutoRenewClient(ctx, testClient)
							Expect(err).To(MatchError(expAuthErr))
						})
						Expect(err).To(BeNil())
					})

					It("should return an error at create if the underlying client returns an error on the first call", func() {
						testCase.SetupAuthValid(testClient)
						expErr := testCase.SetupError(testClient)

						err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
							_, err := vault.NewAutoRenewClient(ctx, testClient)
							Expect(err).To(MatchError(expErr))
						})
						Expect(err).To(BeNil())
					})

					It("should return an error at create if the underlying client returns an error after expired secret", func() {
						testCase.SetupAuthValid(testClient)
						expSecret := testCase.SetupValidOnce(testClient, testClient.authSecret, time.Second)
						expErr := testCase.SetupError(testClient)

						err := testutil.WithDeadlineContext(5*time.Second, func(ctx context.Context) {
							testAutoRenewClient, err := vault.NewAutoRenewClient(ctx, testClient)
							Expect(err).To(BeNil())
							select {
							case <-expSecret.ExpireContext().Done():
								// wait to give time for context change to propagate
								time.Sleep(10 * time.Millisecond)
							case <-ctx.Done():
								Fail("Timed out")
							}
							_, err = testCase.Exec(testAutoRenewClient, ctx)
							Expect(err).To(MatchError(expErr))
						})
						Expect(err).To(BeNil())
					})
				}
			})

		})
	}
})
