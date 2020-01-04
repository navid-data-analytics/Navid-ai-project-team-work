package mongo_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/callstats-io/go-common/mongo"
	"github.com/callstats-io/go-common/testutil"
	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeMongoClient struct {
	mongoSecretError error
	mongoSecret      *vault.UserPassSecret
}

func (fmc *fakeMongoClient) MongoSecret(ctx context.Context) (*vault.UserPassSecret, error) {
	if fmc.mongoSecretError != nil {
		return nil, fmc.mongoSecretError
	}
	return fmc.mongoSecret, nil
}

var _ = Describe("StandardClient", func() {
	var testVaultClient *fakeMongoClient
	BeforeEach(func() {
		testVaultClient = &fakeMongoClient{}
	})

	Context("NewStandardClient", func() {
		It("should validate the options", func() {
			opts, err := mongo.OptionsFromEnv()
			Expect(err).To(BeNil())
			_, err = mongo.NewStandardClient(testVaultClient, opts)
			Expect(err).To(BeNil())
			opts.ConnectionTemplate = ""
			_, err = mongo.NewStandardClient(testVaultClient, opts)
			Expect(err).To(MatchError(mongo.ErrEmptyConnectionTemplate))
		})
	})

	Context("Session", func() {
		var testMongoStandardClient *mongo.StandardClient
		BeforeEach(func() {
			opts, err := mongo.OptionsFromEnv()
			Expect(err).To(BeNil())
			testMongoStandardClient, err = mongo.NewStandardClient(testVaultClient, opts)
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			testMongoStandardClient.Close()
		})

		It("should return a new session", func() {
			secret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
				LeaseDuration: 60,
				Data: map[string]interface{}{
					"username": "vault",
					"password": "vault",
				},
			}, nil))
			Expect(err).To(BeNil())
			testVaultClient.mongoSecret = secret
			ses, err := testMongoStandardClient.Session(context.Background())
			Expect(err).To(BeNil())
			defer ses.Close()
			_, err = ses.Status()
			Expect(err).To(BeNil())
		})
		It("should error if vault mongo secret fetch returns an error", func() {
			testErr := errors.New("test error")
			testVaultClient.mongoSecretError = testErr
			_, err := testMongoStandardClient.Session(context.Background())
			Expect(err).To(MatchError(testErr))
		})
		It("should return an error if mongo dial fails", func() {
			secret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
				LeaseDuration: 60,
				Data: map[string]interface{}{
					"username": "invalid",
					"password": "invalid",
				},
			}, nil))
			Expect(err).To(BeNil())
			testVaultClient.mongoSecret = secret
			ses, err := testMongoStandardClient.Session(context.Background())
			if ses != nil {
				ses.Close()
			}
			Expect(err).To(MatchError(errors.New("server returned error on SASL authentication step: Authentication failed.")))
			_, err = testMongoStandardClient.Status(context.Background())
			Expect(err).To(MatchError(errors.New("server returned error on SASL authentication step: Authentication failed.")))
		})
		It("should use the cached session if the mongo secret id has not changed", func() {
			validMongoSecret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
				LeaseID:       "id1",
				LeaseDuration: 60,
				Data: map[string]interface{}{
					"username": "vault",
					"password": "vault",
				},
			}, nil))
			// test valid session
			testVaultClient.mongoSecret = validMongoSecret
			ses, err := testMongoStandardClient.Session(context.Background())
			Expect(err).To(BeNil())
			defer ses.Close()
			_, err = ses.Status()
			Expect(err).To(BeNil())

			// set the fake StandardClient to have a secret with new credentials but same ID, which should not cause a new connection in the mongo StandardClient
			invalidMongoSecret, err := vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
				LeaseID:       "id1",
				LeaseDuration: 60,
				Data: map[string]interface{}{
					"username": "invalid",
					"password": "invalid",
				},
			}, nil))

			// update the ID to cause a "new" secret to appear
			testVaultClient.mongoSecret = invalidMongoSecret
			ses2, err := testMongoStandardClient.Session(context.Background())
			Expect(err).To(MatchError(errors.New("server returned error on SASL authentication step: Authentication failed.")))
			if ses2 != nil {
				ses2.Close()
			}
		})
	})

	Context("Status", func() {
		var (
			testMongoStandardClient *mongo.StandardClient
		)

		BeforeEach(func() {
			var err error
			testVaultClient.mongoSecret, err = vault.NewUserPassSecret(vault.NewStandardSecret(&api.Secret{
				LeaseID:       "id1",
				LeaseDuration: 60,
				Data: map[string]interface{}{
					"username": "vault",
					"password": "vault",
				},
			}, nil))
			Expect(err).To(BeNil())
			opts, err := mongo.OptionsFromEnv()
			Expect(err).To(BeNil())
			testMongoStandardClient, err = mongo.NewStandardClient(testVaultClient, opts)
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			testMongoStandardClient.Close()
		})

		Context("with a healthy connection", func() {
			It("should return the server status", func() {
				status, err := testMongoStandardClient.Status(context.Background())
				Expect(err).To(BeNil())
				Expect(status["ok"]).To(Equal(float64(1)))
			})
		})

		Context("with a closed connection", func() {
			BeforeEach(func() {
				testMongoStandardClient.Close()
			})

			It("should return an error", func() {
				status, err := testMongoStandardClient.Status(context.Background())
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError(errors.New("StandardClient has been closed")))
				Expect(status).To(BeNil())
			})
		})
	})
})

var _ = Describe("StaticClient", func() {
	var testOpts *mongo.Options
	BeforeEach(func() {
		opts, err := mongo.OptionsFromEnv()
		Expect(err).To(BeNil())
		opts.ConnectionTemplate = fmt.Sprintf(opts.ConnectionTemplate, "vault", "vault")
		testOpts = opts
	})
	Context("NewStaticClient", func() {
		Context("Success", func() {
			It("should return a new client with valid connection", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					client, err := mongo.NewStaticClient(ctx, testOpts)
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					ses, err := client.Session(ctx)
					Expect(err).To(BeNil())
					defer ses.Close()
					_, err = ses.Status()
					Expect(err).To(BeNil())
				})
			})
		})
		Context("Failure", func() {
			It("should return an error", func() {
				testOpts.ConnectionTemplate = "abc%2F"
				testutil.WithCancelContext(func(ctx context.Context) {
					_, err := mongo.NewStaticClient(ctx, testOpts)
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})

	Context("Session", func() {
		Context("Success", func() {
			It("should return a valid connection", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					client, err := mongo.NewStaticClient(ctx, testOpts)
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					ses, err := client.Session(ctx)
					Expect(err).To(BeNil())
					defer ses.Close()
					_, err = ses.Status()
					Expect(err).To(BeNil())
				})
			})
		})

		Context("Failure", func() {
			It("should return an error if client is closed", func() {
				testutil.WithCancelContext(func(ctx context.Context) {
					client, err := mongo.NewStaticClient(ctx, testOpts)
					Expect(err).To(BeNil())
					Expect(client).ToNot(BeNil())
					client.Close()
					_, err = client.Session(ctx)
					Expect(err).To(MatchError(mongo.ErrClosed))
				})
			})
		})
	})
})
