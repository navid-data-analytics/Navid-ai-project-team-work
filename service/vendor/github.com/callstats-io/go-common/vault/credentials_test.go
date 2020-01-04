package vault_test

import (
	"encoding/json"
	"os"

	"github.com/callstats-io/go-common/vault"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppRoleCredentials", func() {
	Describe("Validate", func() {
		It("should return an error if role is invalid", func() {
			creds := vault.AppRoleCredentials{
				SecretID: "test",
			}
			Expect(creds.Validate()).To(Equal(vault.ErrInvalidRoleID))
		})
		It("should return an error if secret is invalid", func() {
			creds := vault.AppRoleCredentials{
				RoleID: "test",
			}
			Expect(creds.Validate()).To(Equal(vault.ErrInvalidSecretID))
		})
	})
	Describe("ReadEnvironment", func() {
		var prevCreds string
		BeforeEach(func() {
			prevCreds = os.Getenv(vault.EnvVaultAppRoleCreds)
		})
		AfterEach(func() {
			os.Setenv(vault.EnvVaultAppRoleCreds, prevCreds)
		})
		Context("Success", func() {
			It("should read json creds from env", func() {
				expCreds := &vault.AppRoleCredentials{
					RoleID:   "envsonroleid",
					SecretID: "envjsonsecretid",
				}
				credsJSON, _ := json.Marshal(expCreds)
				os.Setenv(vault.EnvVaultAppRoleCreds, string(credsJSON))
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadEnvironment()).To(BeNil())
				Expect(creds.SecretID).To(Equal(expCreds.SecretID))
				Expect(creds.RoleID).To(Equal(expCreds.RoleID))
			})
			It("should read file creds from env", func() {
				creds := &vault.AppRoleCredentials{}
				os.Setenv(vault.EnvVaultAppRoleCreds, "file:./testdata/approle_valid.json")
				Expect(creds.ReadEnvironment()).To(BeNil())
				Expect(creds.Validate()).To(BeNil())
				Expect(creds.SecretID).To(Equal("supersecretsecretid"))
				Expect(creds.RoleID).To(Equal("supersecretroleid"))
			})
			It("should read kubernetes creds from env", func() {
				creds := &vault.AppRoleCredentials{}
				os.Setenv(vault.EnvVaultAppRoleCreds, "kubernetes:./testdata/approle_kubernetes")
				Expect(creds.ReadEnvironment()).To(BeNil())
				Expect(creds.Validate()).To(BeNil())
				Expect(creds.SecretID).To(Equal("supersecretkubesecretid"))
				Expect(creds.RoleID).To(Equal("supersecretkuberoleid"))
			})
		})

		Context("Failure", func() {
			It("should fail if app role credentials cannot be found from environment", func() {
				os.Setenv(vault.EnvVaultAppRoleCreds, "")
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadEnvironment()).To(MatchError(vault.ErrEmptyEnvAppRoleCreds))
			})
			It("should fail if read json creds from env fails", func() {
				os.Setenv(vault.EnvVaultAppRoleCreds, "abc")
				creds := &vault.AppRoleCredentials{}
				err := creds.ReadEnvironment()
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid character 'a' looking for beginning of value"))
			})
			It("should fail if read file creds from env fails", func() {
				creds := &vault.AppRoleCredentials{}
				os.Setenv(vault.EnvVaultAppRoleCreds, "file:./testdata/approle_invalid_format.json")
				err := creds.ReadEnvironment()
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("json: cannot unmarshal number into Go struct field AppRoleCredentials.role_id of type string"))
			})
			It("should fail if read kubernetes creds from env fails", func() {
				creds := &vault.AppRoleCredentials{}
				os.Setenv(vault.EnvVaultAppRoleCreds, "kubernetes:./testdata/invalid")
				err := creds.ReadEnvironment()
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("open testdata/invalid/role_id: no such file or directory"))
			})
		})
	})
	Describe("ReadFile", func() {
		Context("Success", func() {
			It("should parse file app role credentials from environment", func() {
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadFile("./testdata/approle_valid.json")).To(BeNil())
				Expect(creds.SecretID).To(Equal("supersecretsecretid"))
				Expect(creds.RoleID).To(Equal("supersecretroleid"))
			})
		})

		Context("Failure", func() {
			It("should fail if the file cannot be found", func() {
				creds := &vault.AppRoleCredentials{}
				err := creds.ReadFile("../unknown/invalid/file")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("open ../unknown/invalid/file: no such file or directory"))
			})
			It("should fail if the json is invalidly formatted", func() {
				creds := &vault.AppRoleCredentials{}
				err := creds.ReadFile("./testdata/approle_invalid_format.json")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("json: cannot unmarshal number into Go struct field AppRoleCredentials.role_id of type string"))
			})
			It("should fail if the json is invalid", func() {
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadFile("./testdata/approle_invalid_content.json")).To(MatchError(vault.ErrInvalidRoleID))
			})
		})
	})
	Describe("ReadKubernetes", func() {
		Context("Success", func() {
			It("should parse kubernetes app role credentials from files", func() {
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadKubernetes("./testdata/approle_kubernetes")).To(BeNil())
				Expect(creds.RoleID).To(Equal("supersecretkuberoleid"))
				Expect(creds.SecretID).To(Equal("supersecretkubesecretid"))
			})
		})

		Context("Failure", func() {
			It("should fail if app role credentials cannot be found", func() {
				creds := &vault.AppRoleCredentials{}
				err := creds.ReadKubernetes("./testdata/approle_kubernetes_unknown")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("open testdata/approle_kubernetes_unknown/role_id: no such file or directory"))
			})
			It("should fail if the credentials are invalid", func() {
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadKubernetes("./testdata/approle_kubernetes_invalid_secret")).To(Equal(vault.ErrInvalidSecretID))
			})
			It("should fail if the credentials are invalid", func() {
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadKubernetes("./testdata/approle_kubernetes_invalid_role")).To(Equal(vault.ErrInvalidRoleID))
			})
		})
	})
	Describe("ReadJSON", func() {
		Context("Success", func() {
			It("should successfully read credentials from json", func() {
				expCreds := &vault.AppRoleCredentials{
					SecretID: "secretjsonid",
					RoleID:   "rolejsonid",
				}
				content, err := json.Marshal(expCreds)
				Expect(err).To(BeNil())
				creds := &vault.AppRoleCredentials{}
				Expect(creds.ReadJSON(string(content))).To(BeNil())
				Expect(creds).To(Equal(expCreds))
			})
		})

		Context("Failure", func() {
			It("should fail if the json is invalidly formatted", func() {
				creds := &vault.AppRoleCredentials{}
				err := creds.ReadJSON("abc")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid character 'a' looking for beginning of value"))
			})
			It("should fail if the json is invalid", func() {
				creds := &vault.AppRoleCredentials{}
				credsJSON, _ := json.Marshal(&vault.AppRoleCredentials{
					SecretID: "abc",
				})
				Expect(creds.ReadJSON(string(credsJSON))).To(Equal(vault.ErrInvalidRoleID))
			})
		})
	})

	Describe("Map", func() {
		It("should return the creds as map[string]interface", func() {
			creds := &vault.AppRoleCredentials{}
			credsMap := map[string]interface{}{
				"secret_id": creds.SecretID,
				"role_id":   creds.RoleID,
			}
			Expect(creds.Map()).To(Equal(credsMap))
		})
	})
})
