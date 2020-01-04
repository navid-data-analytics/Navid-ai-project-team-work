package vault_test

import (
	"errors"
	"os"

	"github.com/callstats-io/go-common/vault"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options", func() {
	makeValidFullOptions := func() *vault.Options {
		return &vault.Options{
			AppRoleAuthPath:       "/fake/approle/auth",
			EnableMongo:           true,
			MongoCredsReadPath:    "/fake/mongo/creds/path",
			EnablePostgres:        true,
			PostgresCredsReadPath: "/fake/postgres/creds/path",
			EnableTLSCert:         true,
			TLSCertReadPath:       "/fake/tls/certs/path",
			TLSCertKeyReadPath:    "/fake/tls/certkeys/path",
			EnableAWS:             true,
			AWSCredsReadPath:      "/fake/aws/sts/path",
		}
	}
	Describe("Validate", func() {
		Context("Success", func() {
			It("should be valid", func() {
				// minimal valid options
				opts := makeValidFullOptions()
				Expect(opts.Validate()).To(BeNil())
			})
		})
		Context("Failure", func() {
			It("should return an error if AppRoleAuthPath is not set", func() {
				opts := makeValidFullOptions()
				opts.AppRoleAuthPath = ""
				Expect(opts.Validate()).To(Equal(vault.ErrEmptyAppRoleAuthPath))
			})
			Context("With EnableMongo = true", func() {
				It("should return an error if MongoCredsReadPath is not set", func() {
					opts := makeValidFullOptions()
					opts.MongoCredsReadPath = ""
					Expect(opts.Validate()).To(Equal(vault.ErrEmptyMongoCredsReadPath))
				})
			})
			Context("With EnablePostgres = true", func() {
				It("should return an error if PostgresCredsReadPath is not set", func() {
					opts := makeValidFullOptions()
					opts.PostgresCredsReadPath = ""
					Expect(opts.Validate()).To(Equal(vault.ErrEmptyPostgresCredsReadPath))
				})
			})
			Context("With EnableTLSCert = true", func() {
				It("should return an error if TLSCertReadPath is not set", func() {
					opts := makeValidFullOptions()
					opts.TLSCertReadPath = ""
					Expect(opts.Validate()).To(Equal(vault.ErrEmptyTLSCertReadPath))
				})
				It("should return an error if TLSCertKeyReadPath is not set", func() {
					opts := makeValidFullOptions()
					opts.TLSCertKeyReadPath = ""
					Expect(opts.Validate()).To(Equal(vault.ErrEmptyTLSCertKeyReadPath))
				})
			})
			Context("With EnableAWS = true", func() {
				It("should return an error if AWSCredsReadPath is not set", func() {
					opts := makeValidFullOptions()
					opts.AWSCredsReadPath = ""
					Expect(opts.Validate()).To(Equal(vault.ErrEmptyAWSCredsReadPath))
				})
			})
		})
	})
})

type optionsEnvTestCase struct {
	// Test case title
	Title string

	// Enable* variable name + value (e.g. EnvVaultEnableMongo and "true")
	EnvVariableEnable string
	ValueEnable       string

	// Variable name + value settings (e.g. EnvVaultMongoCredsReadPath and /abc/def/g)
	EnvVariable string
	Value       string

	// Expected error
	ExpError error
}

var _ = Describe("OptionsFromEnv", func() {
	Describe("Validate", func() {
		Context("Success", func() {
			It("should be valid", func() {
				opts, err := vault.OptionsFromEnv()
				Expect(err).To(BeNil())
				Expect(opts).ToNot(BeNil())
			})
			It("should default AppRoleAuthPath to auth/test/approle/login", func() {
				prev := os.Getenv(vault.EnvVaultAppRoleCreds)
				os.Unsetenv(vault.EnvVaultAppRoleCreds)
				defer os.Setenv(vault.EnvVaultAppRoleCreds, prev)

				opts, err := vault.OptionsFromEnv()
				Expect(err).To(BeNil())
				Expect(opts.AppRoleAuthPath).To(Equal("auth/test/approle/login"))
			})
		})
		Context("Failure", func() {
			testCases := []optionsEnvTestCase{
				optionsEnvTestCase{
					Title:             "EnableMongo is not a boolean",
					EnvVariableEnable: vault.EnvVaultEnableMongo,
					ValueEnable:       "abc",
					ExpError:          errors.New("Failed to parse VAULT_ENABLE_MONGO, error: strconv.ParseBool: parsing \"abc\": invalid syntax"),
				},
				optionsEnvTestCase{
					Title:             "EnablePostgres is not a boolean",
					EnvVariableEnable: vault.EnvVaultEnablePostgres,
					ValueEnable:       "abc",
					ExpError:          errors.New("Failed to parse VAULT_ENABLE_POSTGRES, error: strconv.ParseBool: parsing \"abc\": invalid syntax"),
				},
				optionsEnvTestCase{
					Title:             "EnableTLSCert is not a boolean",
					EnvVariableEnable: vault.EnvVaultEnableTLSCert,
					ValueEnable:       "abc",
					ExpError:          errors.New("Failed to parse VAULT_ENABLE_TLS_CERT, error: strconv.ParseBool: parsing \"abc\": invalid syntax"),
				},
				optionsEnvTestCase{
					Title:             "EnableMongo = true and mongo creds read path is empty",
					EnvVariableEnable: vault.EnvVaultEnableMongo,
					ValueEnable:       "true",
					EnvVariable:       vault.EnvVaultMongoCredsPath,
					Value:             "",
					ExpError:          vault.ErrEmptyMongoCredsReadPath,
				},
				optionsEnvTestCase{
					Title:             "EnablePostgres = true and postgres creds read path is empty",
					EnvVariableEnable: vault.EnvVaultEnablePostgres,
					ValueEnable:       "true",
					EnvVariable:       vault.EnvVaultPostgresCredsPath,
					Value:             "",
					ExpError:          vault.ErrEmptyPostgresCredsReadPath,
				},
				optionsEnvTestCase{
					Title:             "EnableTLSCert = true and tls cert read path is empty",
					EnvVariableEnable: vault.EnvVaultEnableTLSCert,
					ValueEnable:       "true",
					EnvVariable:       vault.EnvVaultTLSCertPath,
					Value:             "",
					ExpError:          vault.ErrEmptyTLSCertReadPath,
				},
				optionsEnvTestCase{
					Title:             "EnableTLSCert = true and tls cert key read path is empty",
					EnvVariableEnable: vault.EnvVaultEnableTLSCert,
					ValueEnable:       "true",
					EnvVariable:       vault.EnvVaultTLSCertKeyPath,
					Value:             "",
					ExpError:          vault.ErrEmptyTLSCertKeyReadPath,
				},
			}

			for idx := range testCases {
				testCase := testCases[idx]
				It("should fail when "+testCase.Title, func() {
					prevEnable := os.Getenv(testCase.EnvVariableEnable)
					os.Setenv(testCase.EnvVariableEnable, testCase.ValueEnable)
					defer os.Setenv(testCase.EnvVariableEnable, prevEnable)

					if testCase.EnvVariable != "" {
						prev := os.Getenv(testCase.EnvVariable)
						os.Setenv(testCase.EnvVariable, testCase.Value)
						defer os.Setenv(testCase.EnvVariable, prev)
					}

					_, err := vault.OptionsFromEnv()
					Expect(err).To(Equal(testCase.ExpError))
				})
			}
		})
	})
})
