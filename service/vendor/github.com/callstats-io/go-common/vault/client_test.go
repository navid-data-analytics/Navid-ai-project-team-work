package vault_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/h2non/gock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ctxKey int

const (
	ctxKeyMockRequest ctxKey = iota
)

// TODO refactor these tests, so much hackety whackety.
// Depends on deprecation of fixed read methods in favor of generic Read.
type secretTestCase struct {
	Name                string
	SuccessMessage      string
	FailureMessage      string
	RenewSuccessMessage string
	RenewFailureMessage string
	Exec                func(vault.Client, context.Context, string) (vault.Secret, error)
	ValidateSecret      func(vault.Secret)
	GetPath             func(*vault.Options, bool) string
	SetEnabledFlag      func(*vault.Options, bool)
	ErrDisabled         error
	ErrBadRequest       error
	Renewable           bool
}

var _ = Describe("NewStandardClient", func() {
	Context("Failure", func() {
		It("should fail if the opts are not valid", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				_, err := vault.NewStandardClient(ctx, &vault.Options{})
				Expect(err).To(MatchError(vault.ErrEmptyAppRoleAuthPath))
			})
		})
		It("should fail if the env doesn't have correct vault config", func() {
			prev := os.Getenv("VAULT_SKIP_VERIFY")
			defer os.Setenv("VAULT_SKIP_VERIFY", prev)

			os.Setenv("VAULT_SKIP_VERIFY", "abc")
			testutil.WithCancelContext(func(ctx context.Context) {
				_, err := vault.NewStandardClient(ctx, &vault.Options{
					AppRoleAuthPath: "abc",
				})
				Expect(err).To(MatchError(errors.New("Could not parse VAULT_SKIP_VERIFY")))
			})
		})
		It("should fail if vault addr is invalid", func() {
			prev := os.Getenv("VAULT_ADDR")
			defer os.Setenv("VAULT_ADDR", prev)

			os.Setenv("VAULT_ADDR", "http://%s.com") //invalid url
			testutil.WithCancelContext(func(ctx context.Context) {
				_, err := vault.NewStandardClient(ctx, &vault.Options{
					AppRoleAuthPath: "abc",
				})
				_, expErr := url.Parse("http://%s.com")
				Expect(err).To(MatchError(expErr))
			})
		})
	})
})

var _ = Describe("StandardClient", func() {
	newClientAndOptions := func(ctx context.Context) (*vault.StandardClient, *vault.Options) {
		opts, err := vault.OptionsFromEnv()
		Expect(err).To(BeNil())
		client, err := vault.NewStandardClient(ctx, opts)
		Expect(err).To(BeNil())
		return client, opts
	}

	Describe("Options", func() {
		It("should return the clients options", func() {
			err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				testClient, testOptions := newClientAndOptions(ctx)
				Expect(testClient.Options()).To(Equal(testOptions))
			})
			Expect(err).To(BeNil())
		})
	})

	Describe("Authenticate", func() {
		Context("Success", func() {
			It("should authenticate without error", func() {
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					testClient, _ := newClientAndOptions(ctx)
					secret, err := testClient.Authenticate(ctx)
					Expect(err).To(BeNil())
					Expect(secret).ToNot(BeNil())
				})
				Expect(err).To(BeNil())
			})

			It("should renew the token if already authenticated", func() {
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					testClient, _ := newClientAndOptions(ctx)
					secret, err := testClient.Authenticate(ctx)
					Expect(err).To(BeNil())
					Expect(secret).ToNot(BeNil())
					secret2, err := testClient.Authenticate(ctx)
					Expect(err).To(BeNil())
					Expect(secret.ID()).ToNot(Equal(secret2.ID()))
					Expect(secret.Auth.LeaseDuration).To(BeNumerically(">", 0))
					Expect(secret.Auth.Accessor).To(Equal(secret2.Auth.Accessor))
					Expect(secret.RenewTime().Before(secret2.RenewTime())).To(BeTrue())
					Expect(secret.ExpireTime().Before(secret2.ExpireTime())).To(BeTrue())
				})
				Expect(err).To(BeNil())
			})
		})
		Context("Failure", func() {
			It("should fail if credentials are not in environment", func() {
				prevCreds := os.Getenv(vault.EnvVaultAppRoleCreds)
				defer os.Setenv(vault.EnvVaultAppRoleCreds, prevCreds)

				os.Setenv(vault.EnvVaultAppRoleCreds, "")
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					testClient, _ := newClientAndOptions(ctx)
					_, err := testClient.Authenticate(ctx)
					Expect(err).To(MatchError(vault.ErrEmptyEnvAppRoleCreds))
				})
				Expect(err).To(BeNil())
			})
			It("should fail if vault auth request fails", func() {
				prevCreds := os.Getenv(vault.EnvVaultAppRoleCreds)
				defer os.Setenv(vault.EnvVaultAppRoleCreds, prevCreds)

				credsJSON, _ := json.Marshal(vault.AppRoleCredentials{
					SecretID: "invalid",
					RoleID:   "invalid",
				})
				os.Setenv(vault.EnvVaultAppRoleCreds, string(credsJSON))
				err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
					testClient, _ := newClientAndOptions(ctx)
					_, err := testClient.Authenticate(ctx)
					Expect(err).To(MatchError(errors.New("Error making API request.\n\nURL: PUT http://vault:8200/v1/auth/test/approle/login\nCode: 400. Errors:\n\n* failed to validate SecretID: failed to find secondary index for role_id \"invalid\"\n")))
				})
				Expect(err).To(BeNil())
			})
		})
	})

	readError := errors.New("Error making API request.\n\nURL: GET http://vault:8200/v1/abc\nCode: 403. Errors:\n\n* permission denied")
	writeError := errors.New("Error making API request.\n\nURL: PUT http://vault:8200/v1/abc\nCode: 403. Errors:\n\n* permission denied")

	testCases := []*secretTestCase{
		&secretTestCase{
			Name:                "MongoSecret",
			SuccessMessage:      vault.LogMongoCredsReadSuccess,
			FailureMessage:      vault.LogMongoCredsReadFailure,
			RenewSuccessMessage: vault.LogMongoCredsRenewSuccess,
			RenewFailureMessage: vault.LogMongoCredsRenewFailure,
			Renewable:           true,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				return testClient.MongoSecret(ctx)
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.UserPassSecret)
				Expect(secret.Credentials.User).ToNot(BeEmpty())
				Expect(secret.Credentials.Password).ToNot(BeEmpty())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					opts.MongoCredsReadPath = "abc"
				}
				return opts.MongoCredsReadPath
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				opts.EnableMongo = enabled
			},
			ErrDisabled:   vault.ErrMongoDisabled,
			ErrBadRequest: readError,
		},
		&secretTestCase{
			Name:                "PostgresSecret",
			SuccessMessage:      vault.LogPostgresCredsReadSuccess,
			FailureMessage:      vault.LogPostgresCredsReadFailure,
			RenewSuccessMessage: vault.LogPostgresCredsRenewSuccess,
			RenewFailureMessage: vault.LogPostgresCredsRenewFailure,
			Renewable:           true,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				return testClient.PostgresSecret(ctx)
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.UserPassSecret)
				Expect(secret.Credentials.User).ToNot(BeEmpty())
				Expect(secret.Credentials.Password).ToNot(BeEmpty())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					opts.PostgresCredsReadPath = "abc"
				}
				return opts.PostgresCredsReadPath
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				opts.EnablePostgres = enabled
			},
			ErrDisabled:   vault.ErrPostgresDisabled,
			ErrBadRequest: readError,
		},
		&secretTestCase{
			Name:                "TLSCertSecret (cert)",
			SuccessMessage:      vault.LogTLSCertReadSuccess,
			FailureMessage:      vault.LogTLSCertReadFailure,
			RenewSuccessMessage: vault.LogTLSCertRenewSuccess,
			RenewFailureMessage: vault.LogTLSCertRenewFailure,
			Renewable:           true,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				return testClient.TLSCertSecret(ctx)
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.TLSCertSecret)
				Expect(secret.Certificate).ToNot(BeNil())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					opts.TLSCertReadPath = "abc"
				}
				return opts.TLSCertReadPath
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				opts.EnableTLSCert = enabled
			},
			ErrDisabled:   vault.ErrTLSCertDisabled,
			ErrBadRequest: readError,
		},
		&secretTestCase{
			Name:                "TLSCertSecret (cert key)",
			SuccessMessage:      vault.LogTLSCertKeyReadSuccess,
			FailureMessage:      vault.LogTLSCertKeyReadFailure,
			RenewSuccessMessage: vault.LogTLSCertKeyRenewSuccess,
			RenewFailureMessage: vault.LogTLSCertKeyRenewFailure,
			Renewable:           true,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				return testClient.TLSCertSecret(ctx)
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.TLSCertSecret)
				Expect(secret.Certificate).ToNot(BeNil())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					opts.TLSCertKeyReadPath = "abc"
				}
				return opts.TLSCertKeyReadPath
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				opts.EnableTLSCert = enabled
			},
			ErrDisabled:   vault.ErrTLSCertDisabled,
			ErrBadRequest: readError,
		},
		&secretTestCase{
			Name:                "AWSSecret",
			SuccessMessage:      vault.LogAWSCredsReadSuccess,
			FailureMessage:      vault.LogAWSCredsReadFailure,
			RenewSuccessMessage: vault.LogAWSCredsRenewSuccess,
			RenewFailureMessage: vault.LogAWSCredsRenewFailure,
			Renewable:           true,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				if mock, ok := ctx.Value(ctxKeyMockRequest).(bool); ok && mock {
					defer gock.Off()
					gock.InterceptClient(testClient.VaultHTTPClient())
					gock.New(os.Getenv("VAULT_ADDR")).
						Get("/v1/" + os.Getenv(vault.EnvVaultAWSCredsPath)).
						Reply(200).
						JSON(map[string]interface{}{
							"lease_id":        "aws/abc/sts/go_common/31d771a6-fb39-f46b-fdc5-945109106422",
							"lease_duration":  3600,
							"lease_renewable": true,
							"data": map[string]interface{}{
								"access_key":     "ASIAJYYYY2AA5K4WIXXX",
								"secret_key":     "HSs0DYYYYYY9W81DXtI0K7X84H+OVZXK5BXXXX",
								"security_token": "AQoDYXdzEEwasAKwQyZUtZaCjVNDiXXXXXXXXgUgBBVUUbSyujLjsw6jYzboOQ89vUVIehUw/9MreAifXFmfdbjTr3g6zc0me9M+dB95DyhetFItX5QThw0lEsVQWSiIeIotGmg7mjT1//e7CJc4LpxbW707loFX1TYD1ilNnblEsIBKGlRNXZ+QJdguY4VkzXxv2urxIH0Sl14xtqsRPboV7eYruSEZlAuP3FLmqFbmA0AFPCT37cLf/vUHinSbvw49C4c9WQLH7CeFPhDub7/rub/QU/lCjjJ43IqIRo9jYgcEvvdRkQSt70zO8moGCc7pFvmL7XGhISegQpEzudErTE/PdhjlGpAKGR3d5qKrHpPYK/k480wk1Ai/t1dTa/8/3jUYTUeIkaJpNBnupQt7qoaXXXXXXXXXX",
							},
						})
				}

				return testClient.AWSSecret(ctx)
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.AWSSecret)
				Expect(secret.Credentials.AccessKey).ToNot(BeEmpty())
				Expect(secret.Credentials.SecretKey).ToNot(BeEmpty())
				Expect(secret.Credentials.SecurityToken).ToNot(BeEmpty())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					opts.AWSCredsReadPath = "abc"
				}
				return opts.AWSCredsReadPath
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				opts.EnableAWS = enabled
			},
			ErrDisabled:   vault.ErrAWSDisabled,
			ErrBadRequest: readError,
		},
		&secretTestCase{
			Name:                "Read",
			SuccessMessage:      "Successfully read secret at path " + os.Getenv(vault.EnvVaultPostgresCredsPath),
			FailureMessage:      "Failed to read secret at path abc",
			RenewSuccessMessage: "Successfully renewed secret at path " + os.Getenv(vault.EnvVaultPostgresCredsPath),
			RenewFailureMessage: "Failed to renew secret at path " + os.Getenv(vault.EnvVaultPostgresCredsPath),
			Renewable:           true,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				return testClient.Read(ctx, path)
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.StandardSecret)
				Expect(secret.Data["username"]).ToNot(BeEmpty())
				Expect(secret.Data["password"]).ToNot(BeEmpty())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					return "abc"
				}
				return opts.PostgresCredsReadPath
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				// do nothing, this is not applicable for Read
			},
			ErrBadRequest: readError,
		},
		&secretTestCase{
			Name:           "Write",
			SuccessMessage: "Successfully wrote secret at path test/transit/encrypt/gocommon",
			FailureMessage: "Failed to write secret at path abc",
			Renewable:      false,
			Exec: func(testClient vault.Client, ctx context.Context, path string) (vault.Secret, error) {
				return testClient.Write(ctx, path, map[string]interface{}{"plaintext": base64.StdEncoding.EncodeToString([]byte("abcdef"))})
			},
			ValidateSecret: func(s vault.Secret) {
				secret := s.(*vault.StandardSecret)
				Expect(secret.Data["ciphertext"]).ToNot(BeEmpty())
			},
			GetPath: func(opts *vault.Options, invalid bool) string {
				if invalid {
					return "abc"
				}
				return "test/transit/encrypt/gocommon"
			},
			SetEnabledFlag: func(opts *vault.Options, enabled bool) {
				// do nothing, this is not applicable for Write
			},
			ErrBadRequest: writeError,
		},
	}

	for idx := range testCases {
		testCase := testCases[idx]

		Describe(testCase.Name, func() {
			Context("Success", func() {
				It("should successfully fetch secret", func() {
					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						testClient, testOptions := newClientAndOptions(ctx)
						testCase.SetEnabledFlag(testOptions, true)
						_, err := testClient.Authenticate(ctx)
						Expect(err).To(BeNil())

						secret, err := testCase.Exec(testClient, context.WithValue(ctx, ctxKeyMockRequest, true), testCase.GetPath(testOptions, false))
						Expect(err).To(BeNil())
						testCase.ValidateSecret(secret)
						Expect(testLogBuffer.String()).To(ContainSubstring(testCase.SuccessMessage))
						Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.FailureMessage))
					})
					Expect(err).To(BeNil())
				})
				if testCase.Renewable {
					It("should try to renew an existing secret first", func() {
						err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
							testClient, testOptions := newClientAndOptions(ctx)
							testCase.SetEnabledFlag(testOptions, true)
							_, err := testClient.Authenticate(ctx)
							Expect(err).To(BeNil())

							secret1, err := testCase.Exec(testClient, context.WithValue(ctx, ctxKeyMockRequest, true), testCase.GetPath(testOptions, false))
							secret2, err := testCase.Exec(testClient, context.WithValue(ctx, ctxKeyMockRequest, true), testCase.GetPath(testOptions, false))
							Expect(err).To(BeNil())
							Expect(secret1.ID()).ToNot(Equal(secret2.ID()))
							Expect(secret1.VaultID()).To(Equal(secret2.VaultID()))
							Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.FailureMessage))

							if isRenewable(secret1) {
								Expect(testLogBuffer.String()).To(ContainSubstring(testCase.RenewSuccessMessage))
								Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.RenewFailureMessage))
							} else {
								Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.RenewSuccessMessage))
								Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.RenewFailureMessage))
							}
						})
						Expect(err).To(BeNil())
					})
				}
			})
			Context("Failure", func() {
				It("should fail if the client is not authenticated", func() {
					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						testClient, testOptions := newClientAndOptions(ctx)
						testCase.SetEnabledFlag(testOptions, true)
						secret, err := testCase.Exec(testClient, ctx, testCase.GetPath(testOptions, false))
						Expect(err).To(MatchError(vault.ErrUnauthenticated))
						Expect(secret).To(BeNil())
						Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.SuccessMessage))
					})
					Expect(err).To(BeNil())
				})
				It("should fail if vault request fails", func() {
					err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
						testClient, testOptions := newClientAndOptions(ctx)
						testCase.SetEnabledFlag(testOptions, true)
						_, err := testClient.Authenticate(ctx)
						Expect(err).To(BeNil())
						secret, err := testCase.Exec(testClient, ctx, testCase.GetPath(testOptions, true))
						Expect(err).To(MatchError(testCase.ErrBadRequest))
						Expect(secret).To(BeNil())
						Expect(testLogBuffer.String()).To(ContainSubstring(testCase.FailureMessage))
						Expect(testLogBuffer.String()).ToNot(ContainSubstring(testCase.SuccessMessage))
					})
					Expect(err).To(BeNil())
				})
			})
			// Hacky way to ignore this for the Read case
			// TODO remove this test once the existing code can be removed
			if testCase.Name != "Read" && testCase.Name != "Write" {
				Context("Invalid setup", func() {
					It("should fail if the method is not enabled", func() {
						err := testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
							testClient, testOptions := newClientAndOptions(ctx)
							testCase.SetEnabledFlag(testOptions, false)
							secret, err := testCase.Exec(testClient, ctx, testCase.GetPath(testOptions, false))
							Expect(err).To(MatchError(testCase.ErrDisabled))
							Expect(secret).To(BeNil())
						})
						Expect(err).To(BeNil())
					})
				})
			}
		})
	}
})
