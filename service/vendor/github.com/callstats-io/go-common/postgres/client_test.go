package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/callstats-io/go-common/postgres"
	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakePostgresClient struct {
	postgresSecretError error
	postgresSecret      *vault.UserPassSecret
}

func (fmc *fakePostgresClient) PostgresSecret(ctx context.Context) (*vault.UserPassSecret, error) {
	if fmc.postgresSecretError != nil {
		return nil, fmc.postgresSecretError
	}
	return fmc.postgresSecret, nil
}

var _ = Describe("ParseURL", func() {
	Context("Success", func() {
		It("should return correct user", func() {
			conf, err := postgres.ParseURL("postgres://user:pass@testhost:1111/testdb")
			Expect(err).To(BeNil())
			Expect(conf.User).To(Equal("user"))
		})

		It("should return correct password", func() {
			conf, err := postgres.ParseURL("postgres://user:pass@testhost:1111/testdb")
			Expect(err).To(BeNil())
			Expect(conf.Password).To(Equal("pass"))
		})

		It("should return correct host", func() {
			conf, err := postgres.ParseURL("postgres://user:pass@testhost:1111/testdb")
			Expect(err).To(BeNil())
			Expect(conf.Addr).To(Equal("testhost:1111"))
		})

		It("should return correct db", func() {
			conf, err := postgres.ParseURL("postgres://user:pass@testhost:1111/testdb")
			Expect(err).To(BeNil())
			Expect(conf.Database).To(Equal("testdb"))
		})
	})
	Context("Failure", func() {
		It("should return error on invalid url", func() {
			_, err := postgres.ParseURL("%gh&%ij")
			Expect(err).ToNot(BeNil())
		})

		It("should return error on missing host", func() {
			_, err := postgres.ParseURL("postgres://user:pass@/testdb")
			Expect(err).To(MatchError(postgres.ErrInvalidAddr))
		})

		It("should return error on missing db", func() {
			_, err := postgres.ParseURL("postgres://user:pass@testhost:1111")
			Expect(err).To(MatchError(postgres.ErrInvalidDB))
		})
	})
})

var _ = Describe("StandardClient", func() {
	var testVaultClient *fakePostgresClient
	BeforeEach(func() {
		testVaultClient = &fakePostgresClient{}
	})

	Context("NewStandardClient", func() {
		It("should validate the options", func() {
			testutil.WithCancelContext(func(ctx context.Context) {
				opts, err := postgres.OptionsFromEnv()
				Expect(err).To(BeNil())
				_, err = postgres.NewStandardClient(ctx, testVaultClient, opts)
				Expect(err).To(BeNil())
				opts.ConnectionTemplate = ""
				_, err = postgres.NewStandardClient(ctx, testVaultClient, opts)
				Expect(err).To(Equal(postgres.ErrEmptyConnectionTemplate))
			})
		})
	})

	Context("DB", func() {
		newTestPostgresClient := func(ctx context.Context, opts *postgres.Options) *postgres.StandardClient {
			if opts == nil {
				o, err := postgres.OptionsFromEnv()
				Expect(err).To(BeNil())
				opts = o
			}
			testPostgresClient, err := postgres.NewStandardClient(ctx, testVaultClient, opts)
			Expect(err).To(BeNil())
			return testPostgresClient
		}

		Context("Success", func() {
			validateConnection := func(db *postgres.DB) {
				_, err := db.Exec("Select 1 AS one")
				Expect(err).To(BeNil())
			}
			It("should return an active db pool", func() {
				secret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
					LeaseDuration: 60,
					Data: map[string]interface{}{
						"username": "go_common",
						"password": "test",
					},
				}, nil))
				Expect(err).To(BeNil())
				testVaultClient.postgresSecret = secret

				testutil.WithCancelContext(func(ctx context.Context) {
					testPostgresClient := newTestPostgresClient(ctx, nil)
					defer testPostgresClient.Close()

					db, err := testPostgresClient.DB(context.Background())
					Expect(err).To(BeNil())
					defer db.Close()
					validateConnection(db)
				})
			})

			It("should use the cached db pool if the postgres secret id has not changed", func() {
				validPostgresSecret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
					LeaseID:       "id1",
					LeaseDuration: 60,
					Data: map[string]interface{}{
						"username": "go_common",
						"password": "test",
					},
				}, nil))
				Expect(err).To(BeNil())

				testutil.WithCancelContext(func(ctx context.Context) {
					testPostgresClient := newTestPostgresClient(ctx, nil)
					defer testPostgresClient.Close()
					// test valid db
					testVaultClient.postgresSecret = validPostgresSecret
					db, err := testPostgresClient.DB(context.Background())
					Expect(err).To(BeNil())
					defer db.Close()
					validateConnection(db)

					// set the fake StandardClient to have a secret with new credentials but same ID, which should not cause a new connection in the postgres StandardClient
					invalidPostgresSecret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
						LeaseID:       "id1",
						LeaseDuration: 60,
						Data: map[string]interface{}{
							"username": "invalid",
							"password": "invalid",
						},
					}, nil))
					Expect(err).To(BeNil())

					// update the ID to cause a "new" secret to appear
					testVaultClient.postgresSecret = invalidPostgresSecret
					db2, err := testPostgresClient.DB(context.Background())
					Expect(err).To(BeNil())
					_, err = db2.Exec("Select 1 AS one")
					// Need to check for substring as postgres error includes the server addr (which is not the same e.g. locally and in jenkins)
					Expect(err.Error()).To(ContainSubstring("FATAL #28P01 password authentication failed for user \"invalid\""))
					if db2 != nil {
						db2.Close()
					}
				})
			})
			It("should handle multiple close calls gracefully", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					testPostgresClient := newTestPostgresClient(ctx, nil)
					testPostgresClient.Close()
					Expect(func() { testPostgresClient.Close() }).NotTo(Panic())
				})
			})
		})
		Context("Failure", func() {
			It("should error if vault postgres secret fetch returns an error", func() {
				testErr := errors.New("test error")
				testVaultClient.postgresSecretError = testErr

				testutil.WithCancelContext(func(ctx context.Context) {
					testPostgresClient := newTestPostgresClient(ctx, nil)
					defer testPostgresClient.Close()

					_, err := testPostgresClient.DB(context.Background())
					Expect(err).To(MatchError(testErr))
				})
			})
			It("should error if client is closed", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					testPostgresClient := newTestPostgresClient(ctx, nil)
					testPostgresClient.Close()

					_, err := testPostgresClient.DB(context.Background())
					Expect(err).To(MatchError(postgres.ErrClosed))
				})
			})
			It("should error if connection template is invalid", func() {
				secret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
					LeaseDuration: 60,
					Data: map[string]interface{}{
						"username": "go_common",
						"password": "test",
					},
				}, nil))
				Expect(err).To(BeNil())
				testVaultClient.postgresSecret = secret

				opts := &postgres.Options{
					ConnectionTemplate: "postgres://%s:%s@testhost:5432",
				}
				testutil.WithCancelContext(func(ctx context.Context) {
					testPostgresClient := newTestPostgresClient(ctx, opts)
					defer testPostgresClient.Close()

					testutil.WithCancelContext(func(callCtx context.Context) {
						_, err := testPostgresClient.DB(callCtx)
						Expect(err).To(MatchError(postgres.ErrInvalidDB))
					})
				})
			})
		})
	})
})

var _ = Describe("StaticClient", func() {
	var testOpts *postgres.Options
	BeforeEach(func() {
		opts, err := postgres.OptionsFromEnv()
		Expect(err).To(BeNil())
		// TODO configure this more reliably, for now has to be kept in sync with the docker-compose/Jenkinsfile
		opts.ConnectionTemplate = fmt.Sprintf(opts.ConnectionTemplate, os.Getenv("SERVICE_NAME"), "test")
		testOpts = opts
	})
	Context("NewStaticClient", func() {
		Context("Success", func() {
			It("should return a new client with valid connection", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					client, err := postgres.NewStaticClient(ctx, testOpts)
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					db, err := client.DB(ctx)
					Expect(err).To(BeNil())
					defer db.Close()
					Expect(db.Status()).To(BeNil())
				})
			})
		})
		Context("Failure", func() {
			It("should return an error", func() {
				testOpts.ConnectionTemplate = "abc%2F"
				testutil.WithCancelContext(func(ctx context.Context) {
					_, err := postgres.NewStaticClient(ctx, testOpts)
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})

	Context("DB", func() {
		Context("Success", func() {
			It("should return a valid connection", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					client, err := postgres.NewStaticClient(ctx, testOpts)
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					db, err := client.DB(ctx)
					Expect(err).To(BeNil())
					defer db.Close()
					Expect(db.Status()).To(BeNil())
				})
			})
		})

		Context("Failure", func() {
			It("should return an error if client is closed", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					client, err := postgres.NewStaticClient(ctx, testOpts)
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					client.Close()
					_, err = client.DB(ctx)
					Expect(err).To(MatchError(postgres.ErrClosed))
				})
			})
		})
	})
})
