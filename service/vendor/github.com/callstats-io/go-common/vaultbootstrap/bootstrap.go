package vaultbootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"
)

// Vault secret backends
const (
	SecretBackendPostgres = "postgresql"
	SecretBackendMongo    = "mongodb"
	SecretBackendGeneric  = "generic"
	SecretBackendTransit  = "transit"
	SecretBackendAWS      = "aws"
)

// Environment variable names
const (
	EnvEnv                    = "ENV"
	EnvServiceName            = "SERVICE_NAME"
	EnvVaultRootToken         = "VAULT_TOKEN"
	EnvVaultPostgresRootURL   = "VAULT_POSTGRES_ROOT_URL"
	EnvVaultPostgresName      = "VAULT_POSTGRES_NAME"
	EnvVaultMongoRootURL      = "VAULT_MONGO_ROOT_URL"
	EnvVaultMongoClusterName  = "VAULT_MONGO_CLUSTER_NAME"
	EnvVaultTLSCertName       = "VAULT_TLS_CERT_NAME"
	EnvTLSCertFilePath        = "TLS_CERT_FILE"
	EnvTLSCertKeyFilePath     = "TLS_CERT_KEY_FILE"
	EnvVaultAWSRegion         = "VAULT_AWS_REGION"
	EnvVaultAWSRootAccessKey  = "VAULT_AWS_ROOT_ACCESS_KEY"
	EnvVaultAWSRootSecretKey  = "VAULT_AWS_ROOT_SECRET_KEY"
	EnvVaultAWSRoleName       = "VAULT_AWS_ROLE_NAME"
	EnvVaultAWSRolePolicyJSON = "VAULT_AWS_ROLE_POLICY_JSON"
	EnvVaultAWSAssumeRole     = "VAULT_AWS_ASSUME_ROLE"
	EnvGenericDataPath        = "VAULT_GENERIC_DATA_FILE"
	EnvTransitDataPath        = "VAULT_TRANSIT_DATA_FILE"
)

// BootstrapConfig contains all configuration values that can be set in vault to modify bootstrapping config.
// Generally this is built by calling the With*-methods on the bootstrap client.
type BootstrapConfig struct {
	Env                     string // current environment, one of ("test", "dev", "prod")
	ServiceName             string // unique name for this service, e.g. "service_access_management"
	VaultRootToken          string // vault root token to use for bootstrapping
	AuthTokenTTL            string // TTL for generated auth tokens
	AuthTokenMaxTTL         string // Max TTL for generated auth tokens
	PostgresName            string // Name to use in vault role setup, e.g. dashboard -> /dev/postgresql/dashboard/...
	PostgresRootConnURL     string // root pg url, used to create new roles by vault
	PostgresRoleLeaseTTL    string // TTL for generated postgres roles
	PostgresRoleLeaseMaxTTL string // Max TTL for generated postgres roles
	MongoClusterName        string // Mongodb cluster to set policies for
	MongoRootConnURL        string // root mongodb url, used to create new roles by vault
	MongoRoleLeaseTTL       string // TTL for generated mongodb roles
	MongoRoleLeaseMaxTTL    string // Max TTL for generated mongodb roles
	TLSCertName             string // certificate name, e.g. x509
	TLSCertData             string // certificate data, e.g. single cert or cert chain
	TLSCertKeyData          string // certificate private key
	TLSCertLeaseTTL         string // certificate lease ttl
	GenericData             string // generic data JSON
	TransitData             string // transit data JSON
	AWSRegion               string // aws region to use
	AWSRootAccessKey        string // aws access key to use for configuring root access
	AWSRootSecretKey        string // aws secret key to use for configuring root access
	AWSRoleName             string // aws role name
	AWSRolePolicyJSON       string // aws roles policy json content, used to write the policy for the STS creds
	AWSAssumedRole          string // aws role to assume
}

type internalBootstrapConfig struct {
	servicePolicyName             string                 // vault service policy name, e.g. test/go_common
	authPath                      string                 // vault auth path, e.g. /dev/approle
	authType                      string                 // vault auth type, e.g. approle
	authDesc                      string                 // vault auth description
	authPolicyPath                string                 // vault auth approle policy path
	authPolicyRoleIDPath          string                 // vault approle auth policy role id path
	authPolicySecretIDPath        string                 // vault approle auth policy secret id path
	authPolicyData                map[string]interface{} // vault auth policy data
	postgresPolicyCredsPath       string                 // vault postgres read policy path
	postgresPolicyMountPath       string                 // vault postgres policy root mount path
	postgresPolicyRoleCreatePath  string                 // vault postgres policy path for new role create config
	postgresPolicyRoleCreateSQL   string                 // vault postgres policy role creation sql string
	postgresPolicyRoleLeasePath   string                 // vault postgres policy lease configuration path (e.g. lease time)
	postgresPolicyRoleConnPath    string                 // vault postgres policy connection config for creating new postgres roles
	postgresPolicyRoleLeaseTTL    string                 // vault postgres policy role lease ttl
	postgresPolicyRoleLeaseMaxTTL string                 // vault postgres policy role lease max ttl
	postgresRootConnURL           string                 // postgres connection url to be used to create new roles
	mongoPolicyCredsPath          string                 // vault mongodb read policy path
	mongoPolicyMountPath          string                 // vault mongodb policy root mount path
	mongoPolicyRoleCreatePath     string                 // vault mongodb policy path for new role create config
	mongoPolicyRoleLeasePath      string                 // vault mongodb policy lease configuration path (e.g. lease time)
	mongoPolicyRoleConnPath       string                 // vault mongodb policy connection config for creating new mongo roles
	mongoPolicyRoleLeaseTTL       string                 // vault mongodb policy role lease ttl
	mongoPolicyRoleLeaseMaxTTL    string                 // vault mongodb policy role lease max ttl
	mongoRootConnURL              string                 // mongodb connection url to be used to create new roles
	tlsCertPolicyCertPath         string                 // vault certificate policy cert data read path
	tlsCertPolicyKeyPath          string                 // vault certificate policy key data read path
	tlsCertData                   string                 // vault certificate policy certificate chain data
	tlsCertKeyData                string                 // vault certificate policy certificate private key data
	tlsCertPolicyLeaseTTL         string                 // vault certificate policy lease ttl
	awsMountPath                  string                 // vault aws mount path
	awsConfigRootPath             string                 // vault aws config root path
	awsConfigRootData             map[string]interface{} // vault aws config root data
	awsPolicyRolesPath            string                 // vault aws policy roles path
	awsPolicyRolesPolicy          string                 // vault aws policy roles policy json
	awsPolicyCredsPath            string                 // vault aws policy creds (sts) path
	awsAssumedRoleRolesPath       string                 // vault aws assumed role roles path
	awsAssumedRoleCredsPath       string                 // vault aws assumed role creds (sts) path
	awsAssumedRoleArn             string                 // vault aws assumed role arn
	genericData                   []GenericSecret        // vault generic data
	transitData                   []TransitSecret        // vault transit data
	transitMountPath              string                 // transit mount path
}

// BootstrapClient is the container for bootstrapping client
type BootstrapClient struct {
	VaultClient     *api.Client             // vault client used for bootstrapping
	BootstrapConfig BootstrapConfig         // bootstrap config
	config          internalBootstrapConfig // internal parsed config
	credsJSON       []byte                  // the creds json (available after MountAppRoleAuth)
}

// GenericSecret is a Vault secret with path + values
type GenericSecret struct {
	Path   string                 `json:"path"`
	Values map[string]interface{} `json:"values"`
}

// TransitSecret is a GenericSecret with specific keyring
type TransitSecret struct {
	EncryptPath    string `json:"encryptPath"`
	DecryptPath    string `json:"decryptPath"`
	KeyRingPath    string `json:"keyRingPath"`
	Plaintext      string `json:"plaintext"`
	OutputPath     string `json:"outputPath"`
	OutputTemplate string `json:"outputTemplate"`
}

// NewBootstrapClient creates a new client for bootstrapping vault from environment
func NewBootstrapClient() *BootstrapClient {
	config := api.DefaultConfig()
	config.ReadEnvironment()
	vaultClient, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	postgresName := os.Getenv(EnvVaultPostgresName)
	if postgresName == "" {
		postgresName = os.Getenv(EnvServiceName)
	}
	mongoClusterName := os.Getenv(EnvVaultMongoClusterName)
	if mongoClusterName == "" {
		mongoClusterName = os.Getenv(EnvServiceName)
	}
	awsRegion := os.Getenv(EnvVaultAWSRegion)
	if awsRegion == "" {
		awsRegion = "eu-west-1"
	}

	bootstrapClient := &BootstrapClient{
		VaultClient: vaultClient,
		BootstrapConfig: BootstrapConfig{
			Env:                     os.Getenv(EnvEnv),
			ServiceName:             os.Getenv(EnvServiceName),
			VaultRootToken:          os.Getenv(EnvVaultRootToken),
			AuthTokenTTL:            "20m",
			AuthTokenMaxTTL:         "30m",
			PostgresName:            postgresName,
			PostgresRootConnURL:     os.Getenv(EnvVaultPostgresRootURL),
			PostgresRoleLeaseTTL:    "1h",
			PostgresRoleLeaseMaxTTL: "24h",
			MongoClusterName:        mongoClusterName,
			MongoRootConnURL:        os.Getenv(EnvVaultMongoRootURL),
			MongoRoleLeaseTTL:       "1h",
			MongoRoleLeaseMaxTTL:    "24h",
			TLSCertName:             os.Getenv(EnvVaultTLSCertName),
			TLSCertLeaseTTL:         "24h",
			AWSRegion:               awsRegion,
			AWSRootAccessKey:        os.Getenv(EnvVaultAWSRootAccessKey),
			AWSRootSecretKey:        os.Getenv(EnvVaultAWSRootSecretKey),
			AWSRoleName:             os.Getenv(EnvVaultAWSRoleName),
			AWSAssumedRole:          os.Getenv(EnvVaultAWSAssumeRole),
		},
		config: internalBootstrapConfig{},
	}
	certData, err := ioutil.ReadFile(os.Getenv(EnvTLSCertFilePath))
	if err == nil {
		bootstrapClient.BootstrapConfig.TLSCertData = string(certData)
	}
	certKeyData, err := ioutil.ReadFile(os.Getenv(EnvTLSCertKeyFilePath))
	if err == nil {
		bootstrapClient.BootstrapConfig.TLSCertKeyData = string(certKeyData)
	}
	genericDataJSON, err := ioutil.ReadFile(os.Getenv(EnvGenericDataPath))
	if err == nil {
		bootstrapClient.BootstrapConfig.GenericData = string(genericDataJSON)
	}
	transitDataJSON, err := ioutil.ReadFile(os.Getenv(EnvTransitDataPath))
	if err == nil {
		bootstrapClient.BootstrapConfig.TransitData = string(transitDataJSON)
	}
	awsRolePolicyJSONFilename := os.Getenv(EnvVaultAWSRolePolicyJSON)
	if awsRolePolicyJSONFilename != "" {
		awsRolePolicyJSON, err := ioutil.ReadFile(awsRolePolicyJSONFilename)
		if err == nil {
			bootstrapClient.BootstrapConfig.AWSRolePolicyJSON = string(awsRolePolicyJSON)
		} else {
			fmt.Println(fmt.Errorf("Error %v", err))
		}
	}

	return bootstrapClient
}

// Config returns this clients BootstrapConfig
func (b *BootstrapClient) Config() *BootstrapConfig {
	return &b.BootstrapConfig
}

// WithVaultRootToken sets the vault root token to the specified value. Defaults to ENV["VAULT_TOKEN"].
func (b *BootstrapClient) WithVaultRootToken(token string) *BootstrapClient {
	b.BootstrapConfig.VaultRootToken = token
	return b
}

// WithEnv sets the env to the specified value. Defaults to ENV["ENV"].
func (b *BootstrapClient) WithEnv(env string) *BootstrapClient {
	b.BootstrapConfig.Env = env
	return b
}

// WithServiceName sets the service name to the specified value. Defaults to ENV["SERVICE_NAME"].
func (b *BootstrapClient) WithServiceName(serviceName string) *BootstrapClient {
	b.BootstrapConfig.ServiceName = serviceName
	return b
}

// WithAuthTokenTTL sets the auth token ttl to the specified value. Defaults to "20m".
func (b *BootstrapClient) WithAuthTokenTTL(ttl string) *BootstrapClient {
	b.BootstrapConfig.AuthTokenTTL = ttl
	return b
}

// WithAuthTokenMaxTTL sets the auth token max ttl to the specified value. Defaults to "30m".
func (b *BootstrapClient) WithAuthTokenMaxTTL(maxTTL string) *BootstrapClient {
	b.BootstrapConfig.AuthTokenMaxTTL = maxTTL
	return b
}

// WithPostgresName sets the postgres name to use to the specified value. Defaults to ENV["VAULT_POSTGRES_NAME"], and if that's not present ENV["SERVICE_NAME"].
func (b *BootstrapClient) WithPostgresName(name string) *BootstrapClient {
	b.BootstrapConfig.PostgresName = name
	return b
}

// WithPostgresRootConnURL sets the postgres root connection url to use to the specified value. Defaults to ENV["VAULT_POSTGRES_ROOT_URL"].
func (b *BootstrapClient) WithPostgresRootConnURL(rootURL string) *BootstrapClient {
	b.BootstrapConfig.PostgresRootConnURL = rootURL
	return b
}

// WithPostgresLeaseTTL sets the postgres lease ttl to the specified value. Defaults to "1h".
func (b *BootstrapClient) WithPostgresLeaseTTL(ttl string) *BootstrapClient {
	b.BootstrapConfig.PostgresRoleLeaseTTL = ttl
	return b
}

// WithPostgresLeaseMaxTTL sets the postgres lease ttl to the specified value. Defaults to "24h".
func (b *BootstrapClient) WithPostgresLeaseMaxTTL(maxTTL string) *BootstrapClient {
	b.BootstrapConfig.PostgresRoleLeaseMaxTTL = maxTTL
	return b
}

// WithMongoClusterName sets the mongo name to use to the specified value. Defaults to ENV["VAULT_MONGO_CLUSTER_NAME"], and if that's not present ENV["SERVICE_NAME"].
func (b *BootstrapClient) WithMongoClusterName(name string) *BootstrapClient {
	b.BootstrapConfig.MongoClusterName = name
	return b
}

// WithMongoRootConnURL sets the mongo root connection url to use to the specified value. Defaults to ENV["VAULT_MONGO_ROOT_URL"].
func (b *BootstrapClient) WithMongoRootConnURL(rootURL string) *BootstrapClient {
	b.BootstrapConfig.MongoRootConnURL = rootURL
	return b
}

// WithMongoRoleLeaseTTL sets the mongo lease ttl to the specified value. Defaults to "1h".
func (b *BootstrapClient) WithMongoRoleLeaseTTL(ttl string) *BootstrapClient {
	b.BootstrapConfig.MongoRoleLeaseTTL = ttl
	return b
}

// WithMongoRoleLeaseMaxTTL sets the mongo lease ttl to the specified value. Defaults to "24h".
func (b *BootstrapClient) WithMongoRoleLeaseMaxTTL(maxTTL string) *BootstrapClient {
	b.BootstrapConfig.MongoRoleLeaseMaxTTL = maxTTL
	return b
}

// WithTLSCertName sets the tls cert name to use to the specified value. Defaults to ENV["VAULT_TLS_CERT_NAME"].
func (b *BootstrapClient) WithTLSCertName(certName string) *BootstrapClient {
	b.BootstrapConfig.TLSCertName = certName
	return b
}

// WithTLSCertData sets the tls cert data to use to the specified value.
func (b *BootstrapClient) WithTLSCertData(data string) *BootstrapClient {
	b.BootstrapConfig.TLSCertData = data
	return b
}

// WithTLSCertKeyData sets the tls cert key data to use to the specified value.
func (b *BootstrapClient) WithTLSCertKeyData(data string) *BootstrapClient {
	b.BootstrapConfig.TLSCertKeyData = data
	return b
}

// WithTLSCertLeaseTTL sets the postgres lease ttl to the specified value. Defaults to "1h".
func (b *BootstrapClient) WithTLSCertLeaseTTL(ttl string) *BootstrapClient {
	b.BootstrapConfig.TLSCertLeaseTTL = ttl
	return b
}

// WithTestTLSCertData sets the tls cert data to use to the specified value. Required to be called for TLS setup.
func (b *BootstrapClient) WithTestTLSCertData() *BootstrapClient {
	certData, err := ioutil.ReadFile("../vault/testdata/test_cert.pem")
	if err != nil {
		panic(fmt.Errorf("Failed create read TLS cert file from path pointed to by %s. Error: %v", EnvTLSCertFilePath, err))
	}
	b.BootstrapConfig.TLSCertData = string(certData)
	return b
}

// WithTestTLSCertKeyData sets the tls cert key data to use to the specified value. Required to be called for TLS setup.
func (b *BootstrapClient) WithTestTLSCertKeyData() *BootstrapClient {
	certKeyData, err := ioutil.ReadFile("../vault/testdata/test_cert_key.pem")
	if err != nil {
		panic(fmt.Errorf("Failed create read TLS cert key file from path pointed to by %s. Error: %v", EnvTLSCertKeyFilePath, err))
	}
	b.BootstrapConfig.TLSCertKeyData = string(certKeyData)
	return b
}

// WithTestAWSIAMRolePolicy reads the AWS IAM role policy from Vault package's test data. Required for AWS setup.
// Note: This is an **IAM** role policy, other AWS roles need different policy JSONs
func (b *BootstrapClient) WithTestAWSIAMRolePolicy() *BootstrapClient {
	policyData, err := ioutil.ReadFile("../vault/testdata/aws_invite_email_role_policy.json")
	if err != nil {
		panic(fmt.Errorf("Failed to read AWS IAM role policy file from path from ../vault/testdata/aws_invite_email_role_policy.json. Error: %v", err))
	}
	b.BootstrapConfig.AWSRolePolicyJSON = string(policyData) // TODO: what if we need multiple role policies?
	return b
}

// WithGenericData sets the generic data to use to the specified value.
func (b *BootstrapClient) WithGenericData(data string) *BootstrapClient {
	b.BootstrapConfig.GenericData = data
	return b
}

// WithTransitData sets the generic data to use to the specified value.
func (b *BootstrapClient) WithTransitData(data string) *BootstrapClient {
	b.BootstrapConfig.TransitData = data
	return b
}

// WithTestGenericData sets the tls cert data to use to the specified value. Required to be called for TLS setup.
func (b *BootstrapClient) WithTestGenericData() *BootstrapClient {
	genericData, err := ioutil.ReadFile("../vault/testdata/genericdata.json")
	if err != nil {
		panic(fmt.Errorf("Failed to read generic data from ../vault/testdata/genericdata.json Error: %v", err))
	}
	b.BootstrapConfig.GenericData = string(genericData)
	return b
}

// WithTestTransitData sets the transit backend data. Required to be called for transit setup.
func (b *BootstrapClient) WithTestTransitData() *BootstrapClient {
	transitData, err := ioutil.ReadFile("../vault/testdata/transitdata.json")
	if err != nil {
		panic(fmt.Errorf("Failed to read generic data from ../vault/testdata/transitdata.json Error: %v", err))
	}
	return b.WithTransitData(string(transitData))
}

// WriteCredentialsFile write the bootstrapped credentials to file
func (b *BootstrapClient) WriteCredentialsFile(path string) *BootstrapClient {
	if err := ioutil.WriteFile(path, b.credsJSON, 0644); err != nil {
		panic(fmt.Errorf("Failed to write credentials to file at %s. Error: %v", path, err))
	}
	return b
}

// WriteCredentialsEnv write the bootstrapped credentials to file
func (b *BootstrapClient) WriteCredentialsEnv() *BootstrapClient {
	os.Setenv(vault.EnvVaultAppRoleCreds, string(b.credsJSON))
	return b
}

// UnmountAll removes all current vault configuration
func (b *BootstrapClient) UnmountAll() *BootstrapClient {
	//ignore error as this may fail if there is no old configuration
	b.usingVaultRootToken(func() error {
		// Remove all old configuration, these may fail if there aren't old configurations to remove
		b.DisableAuth()
		b.unmountPostgres()
		b.unmountMongo()
		b.unmountTLSCert()
		b.unmountAWS()
		b.unmountGeneric()
		b.unmountTransit()
		return nil
	})
	return b
}

// MountAppRoleAuth sets up vault authentication with approle
func (b *BootstrapClient) MountAppRoleAuth() *BootstrapClient {
	b.must(b.configureAppRoleAuth())
	b.DisableAuth()
	b.EnableAuth()
	b.must(b.writeAuthAppRolePolicy())
	creds, err := b.readAppRoleCredentials()
	b.must(err)
	credsJSON, err := json.Marshal(creds)
	b.must(err)
	b.credsJSON = credsJSON
	return b
}

// MountMongo sets up all mongo connection related policies
func (b *BootstrapClient) MountMongo() *BootstrapClient {
	b.must(b.configureMongo())
	b.unmountMongo()
	b.must(b.writeMongoPolicy())
	b.must(b.mountMongoRole())
	return b
}

// MountPostgres sets up all postgres connection related policies
func (b *BootstrapClient) MountPostgres() *BootstrapClient {
	b.must(b.configurePostgres())
	b.unmountPostgres()
	b.must(b.writePostgresPolicy())
	b.must(b.mountPostgresRole())
	return b
}

// MountTLSCert sets up all TLS certificate related policies
func (b *BootstrapClient) MountTLSCert() *BootstrapClient {
	b.must(b.configureTLSCert())
	b.unmountTLSCert()
	b.must(b.writeTLSCertPolicy())
	b.must(b.mountTLSCert())
	return b
}

// MountGeneric sets up all generic related policies
func (b *BootstrapClient) MountGeneric() *BootstrapClient {
	b.must(b.configureGeneric())
	b.unmountGeneric()
	b.must(b.writeGenericPolicy())
	b.must(b.mountGenericData())
	return b
}

// MountTransit sets up all transit related policies
func (b *BootstrapClient) MountTransit() *BootstrapClient {
	b.must(b.configureTransit())
	b.unmountTransit()
	b.must(b.writeTransitPolicy())
	b.must(b.mountTransitData())
	return b
}

// MountAWS sets up all postgres connection related policies
func (b *BootstrapClient) MountAWS() *BootstrapClient {
	b.must(b.configureAWS())
	b.unmountAWS()
	b.must(b.writeAWSPolicy())
	b.must(b.mountAWSRole())
	return b
}

// EnableAuth enables auth based on config
func (b *BootstrapClient) EnableAuth() *BootstrapClient {
	b.must(b.configureAppRoleAuth())
	b.must(b.usingVaultRootToken(func() error {
		return b.VaultClient.Sys().EnableAuth(b.config.authPath, b.config.authType, b.config.authDesc)
	}))
	return b
}

// DisableAuth deletes auth based on config
func (b *BootstrapClient) DisableAuth() *BootstrapClient {
	b.must(b.configureAppRoleAuth())
	b.must(b.usingVaultRootToken(func() error {
		return b.VaultClient.Sys().DisableAuth(b.config.authPath)
	}))
	return b
}

// readAppRoleCredentials reads both role id and secret id
func (b *BootstrapClient) readAppRoleCredentials() (*vault.AppRoleCredentials, error) {
	creds := &vault.AppRoleCredentials{}
	err := b.usingVaultRootToken(func() error {
		roleSecret, err := b.VaultClient.Logical().Read(b.config.authPolicyRoleIDPath)
		if err != nil {
			return err
		}
		secret, err := b.VaultClient.Logical().Write(b.config.authPolicySecretIDPath, nil)
		if err != nil {
			return err
		}
		creds.RoleID = roleSecret.Data["role_id"].(string)
		creds.SecretID = secret.Data["secret_id"].(string)
		return nil
	})
	return creds, err
}

// WritePostgresPolicy creates the policy rule for reading postgres user/password
func (b *BootstrapClient) writePostgresPolicy() error {
	return b.usingVaultRootToken(func() error {
		return b.writeReadPolicy(b.config.postgresPolicyCredsPath)
	})
}

// WriteMongoPolicy creates the policy rule for reading mongodb user/password
func (b *BootstrapClient) writeMongoPolicy() error {
	return b.usingVaultRootToken(func() error {
		return b.writeReadPolicy(b.config.mongoPolicyCredsPath)
	})
}

// WriteTLSCertPolicy creates the policy rule for reading certificate cert/key data
func (b *BootstrapClient) writeTLSCertPolicy() error {
	return b.usingVaultRootToken(func() error {
		if err := b.writeReadPolicy(b.config.tlsCertPolicyCertPath); err != nil {
			return err
		}
		return b.writeReadPolicy(b.config.tlsCertPolicyKeyPath)
	})
}

// writeGenericPolicy creates the policy rule for reading generic data
func (b *BootstrapClient) writeGenericPolicy() error {
	return b.usingVaultRootToken(func() error {
		for i := range b.config.genericData {
			if err := b.writeReadPolicy(b.config.genericData[i].Path); err != nil {
				return err
			}
		}
		return nil
	})
}

// writeTransitPolicy creates the policy rule for writing transit data
func (b *BootstrapClient) writeTransitPolicy() error {
	return b.usingVaultRootToken(func() error {
		for i := range b.config.transitData {
			if err := b.writeWritePolicy(b.config.transitData[i].EncryptPath); err != nil {
				return err
			}
			if err := b.writeWritePolicy(b.config.transitData[i].DecryptPath); err != nil {
				return err
			}
			if err := b.writeReadPolicy(b.config.transitData[i].OutputPath); err != nil {
				return err
			}
		}
		return nil
	})
}

// WriteAWSPolicy creates the policy rule for reading AWS credentials
func (b *BootstrapClient) writeAWSPolicy() error {
	return b.usingVaultRootToken(func() error {
		if b.config.awsPolicyCredsPath != "" {
			return b.writeReadPolicy(b.config.awsPolicyCredsPath)
		} else if b.config.awsAssumedRoleCredsPath != "" {
			return b.writeReadPolicy(b.config.awsAssumedRoleCredsPath)
		} else {
			return fmt.Errorf("Neither AWS role policy or assumed role is configured")
		}
	})
}

func (b *BootstrapClient) writeReadPolicy(path string) error {
	policy, err := b.VaultClient.Sys().GetPolicy(b.config.servicePolicyName)
	if err != nil {
		return err
	}
	policy += fmt.Sprintf("\npath \"%s\" {\n    policy = \"read\"\n}", path)
	return b.VaultClient.Sys().PutPolicy(b.config.servicePolicyName, policy)
}

func (b *BootstrapClient) writeWritePolicy(path string) error {
	policy, err := b.VaultClient.Sys().GetPolicy(b.config.servicePolicyName)
	if err != nil {
		return err
	}
	policy += fmt.Sprintf("\npath \"%s\" {\n    policy = \"write\"\n}", path)
	return b.VaultClient.Sys().PutPolicy(b.config.servicePolicyName, policy)
}

// MountPostgresRoleConfig configures postgres role, connection and lease information to vault
func (b *BootstrapClient) mountPostgresRole() error {
	return b.usingVaultRootToken(func() error {
		mountInfo := &api.MountInput{
			Type: SecretBackendPostgres,
		}
		if err := b.VaultClient.Sys().Mount(b.config.postgresPolicyMountPath, mountInfo); err != nil {
			return err
		}

		// set connection data used by vault to connect to postgres
		connectionData := map[string]interface{}{
			"connection_url": b.config.postgresRootConnURL,
		}
		if _, err := b.VaultClient.Logical().Write(b.config.postgresPolicyRoleConnPath, connectionData); err != nil {
			return err
		}

		// create lease configuration
		leaseData := map[string]interface{}{
			"lease":     b.config.postgresPolicyRoleLeaseTTL,
			"lease_max": b.config.postgresPolicyRoleLeaseMaxTTL,
		}
		if _, err := b.VaultClient.Logical().Write(b.config.postgresPolicyRoleLeasePath, leaseData); err != nil {
			return err
		}

		// create pg configuration for creating roles
		roleData := map[string]interface{}{
			"sql": b.config.postgresPolicyRoleCreateSQL,
		}
		if _, err := b.VaultClient.Logical().Write(b.config.postgresPolicyRoleCreatePath, roleData); err != nil {
			return err
		}

		return nil
	})
}

// MountMongoRoleConfig configures mongo role, connection and lease information to vault
func (b *BootstrapClient) mountMongoRole() error {
	return b.usingVaultRootToken(func() error {
		mountInfo := &api.MountInput{
			Type: SecretBackendMongo,
		}
		if err := b.VaultClient.Sys().Mount(b.config.mongoPolicyMountPath, mountInfo); err != nil {
			return err
		}

		// set connection data used by vault to connect to mongo
		connectionData := map[string]interface{}{
			"uri": b.config.mongoRootConnURL,
		}
		if _, err := b.VaultClient.Logical().Write(b.config.mongoPolicyRoleConnPath, connectionData); err != nil {
			return err
		}

		// create lease configuration
		leaseData := map[string]interface{}{
			"ttl":     b.config.mongoPolicyRoleLeaseTTL,
			"ttl_max": b.config.mongoPolicyRoleLeaseMaxTTL,
		}
		if _, err := b.VaultClient.Logical().Write(b.config.mongoPolicyRoleLeasePath, leaseData); err != nil {
			return err
		}

		mongoRoleData := map[string]interface{}{
			"roles": `[{"role":"readWriteAnyDatabase","db":"admin"},"clusterMonitor"]`,
			"db":    "admin",
		}
		if _, err := b.VaultClient.Logical().Write(b.config.mongoPolicyRoleCreatePath, mongoRoleData); err != nil {
			return err
		}

		return nil
	})
}

// MountTLSCertConfig configures certificate cert and key data reading from vault
func (b *BootstrapClient) mountTLSCert() error {
	return b.usingVaultRootToken(func() error {
		configs := map[string](map[string]interface{}){
			// policy path and data for cert
			b.config.tlsCertPolicyCertPath: map[string]interface{}{
				"data": b.config.tlsCertData,
				"ttl":  b.config.tlsCertPolicyLeaseTTL,
			},
			// policy path and data for cert private key
			b.config.tlsCertPolicyKeyPath: map[string]interface{}{
				"data": b.config.tlsCertKeyData,
				"ttl":  b.config.tlsCertPolicyLeaseTTL,
			},
		}
		for path, data := range configs {
			mountInfo := &api.MountInput{
				Type: SecretBackendGeneric,
			}
			if err := b.VaultClient.Sys().Mount(path, mountInfo); err != nil {
				return err
			}
			if _, err := b.VaultClient.Logical().Write(path, data); err != nil {
				return err
			}
		}
		return nil
	})
}

// MountGenericData configures generic data reading from vault
func (b *BootstrapClient) mountGenericData() error {
	return b.usingVaultRootToken(func() error {
		mountInfo := &api.MountInput{
			Type: SecretBackendGeneric,
		}

		for _, secret := range b.config.genericData {
			tmp := make(map[string]interface{}, len(secret.Values))
			for i := range secret.Values {
				tmp[i] = secret.Values[i]
			}
			if err := b.VaultClient.Sys().Mount(secret.Path, mountInfo); err != nil {
				return err
			}
			if _, err := b.VaultClient.Logical().Write(secret.Path, tmp); err != nil {
				return err
			}
		}
		return nil
	})
}

// MountTransitData configures transit data reading from vault
func (b *BootstrapClient) mountTransitData() error {
	return b.usingVaultRootToken(func() error {
		mountInfo := &api.MountInput{
			Type: SecretBackendTransit,
		}
		if err := b.VaultClient.Sys().Mount(b.config.transitMountPath, mountInfo); err != nil {
			return err
		}

		for _, secret := range b.config.transitData {
			if _, err := b.VaultClient.Logical().Write(secret.KeyRingPath, nil); err != nil {
				return err
			}
			encryptedSecret, err := b.VaultClient.Logical().Write(secret.EncryptPath, map[string]interface{}{"plaintext": secret.Plaintext})
			if err != nil {
				return err
			}

			// write payload for output to a generic backend, used e.g. for obfuscation keys
			ciphertext, ok := encryptedSecret.Data["ciphertext"].(string)
			if !ok {
				return errors.New("Invalid ciphertext")
			}
			outputData := map[string]interface{}{}
			if err := json.Unmarshal([]byte(fmt.Sprintf(secret.OutputTemplate, ciphertext)), &outputData); err != nil {
				return err
			}
			mountInfo := &api.MountInput{
				Type: SecretBackendGeneric,
			}
			if err := b.VaultClient.Sys().Mount(secret.OutputPath, mountInfo); err != nil {
				return err
			}
			if _, err := b.VaultClient.Logical().Write(secret.OutputPath, outputData); err != nil {
				return err
			}
		}
		return nil
	})
}

// MountAWSRoleConfig configures AWS role, connection and lease information to vault
func (b *BootstrapClient) mountAWSRole() error {
	return b.usingVaultRootToken(func() error {
		mountInfo := &api.MountInput{
			Type: SecretBackendAWS,
		}
		if err := b.VaultClient.Sys().Mount(b.config.awsMountPath, mountInfo); err != nil {
			return err
		}

		// set root config data used by vault to connect to AWS
		if _, err := b.VaultClient.Logical().Write(b.config.awsConfigRootPath, b.config.awsConfigRootData); err != nil {
			return err
		}

		var rolePath string
		var roleData map[string]interface{}

		// create aws role configuration for generating sts tokens
		if b.config.awsPolicyRolesPath != "" {
			rolePath = b.config.awsPolicyRolesPath
			roleData = map[string]interface{}{
				"policy": b.config.awsPolicyRolesPolicy,
			}
		} else if b.config.awsAssumedRoleRolesPath != "" {
			rolePath = b.config.awsAssumedRoleRolesPath
			roleData = map[string]interface{}{
				"arn": b.config.awsAssumedRoleArn,
			}
		}

		if rolePath == "" || roleData == nil {
			return fmt.Errorf("AWS role configuration is missing")
		}

		if _, err := b.VaultClient.Logical().Write(rolePath, roleData); err != nil {
			return err
		}

		return nil
	})
}

// UnmountPostgres disables postgres role auth
func (b *BootstrapClient) unmountPostgres() error {
	return b.usingVaultRootToken(func() error {
		return b.unmount(b.config.postgresPolicyMountPath)
	})
}

// UnmountMongo disables mongodb role auth
func (b *BootstrapClient) unmountMongo() error {
	return b.usingVaultRootToken(func() error {
		return b.unmount(b.config.mongoPolicyMountPath)
	})
}

// UnmountTLSCert disables certificate reading
func (b *BootstrapClient) unmountTLSCert() error {
	return b.usingVaultRootToken(func() error {
		if err := b.unmount(b.config.tlsCertPolicyCertPath); err != nil {
			return err
		}
		return b.unmount(b.config.tlsCertPolicyKeyPath)
	})
}

// UnmountGeneric disables generic data reading
func (b *BootstrapClient) unmountGeneric() error {
	return b.usingVaultRootToken(func() error {
		for _, entry := range b.config.genericData {
			if err := b.unmount(entry.Path); err != nil {
				return err
			}
		}
		return nil
	})
}

// UnmountTransit disables transit data reading
func (b *BootstrapClient) unmountTransit() error {
	return b.usingVaultRootToken(func() error {
		for _, entry := range b.config.transitData {
			if err := b.unmount(entry.OutputPath); err != nil {
				return err
			}
		}
		return b.unmount(b.config.transitMountPath)
	})
}

// UnmountAWSConfig disables AWS auth
func (b *BootstrapClient) unmountAWS() error {
	return b.usingVaultRootToken(func() error {
		return b.unmount(b.config.awsMountPath)
	})
}

func (b *BootstrapClient) unmount(path string) error {
	return b.VaultClient.Sys().Unmount(path)
}

// WriteAuthAppRolePolicy writes an auth policy with the provided config
func (b *BootstrapClient) writeAuthAppRolePolicy() error {
	return b.usingVaultRootToken(func() error {
		if _, err := b.VaultClient.Logical().Write(b.config.authPolicyPath, b.config.authPolicyData); err != nil {
			return err
		}
		return nil
	})
}
func (b *BootstrapClient) configurePolicyName() error {
	if b.BootstrapConfig.Env == "" {
		return fmt.Errorf("BootstrapConfig.Env cannot be empty when bootstrapping")
	}
	if b.BootstrapConfig.ServiceName == "" {
		return fmt.Errorf("BootstrapConfig.ServiceName cannot be empty when bootstrapping")
	}

	b.config.servicePolicyName = fmt.Sprintf("%s/%s", b.BootstrapConfig.Env, b.BootstrapConfig.ServiceName)
	return nil
}

func (b *BootstrapClient) configureMongo() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	if b.BootstrapConfig.Env == "" {
		return fmt.Errorf("BootstrapConfig.Env cannot be empty when bootstrapping mongo")
	}
	if b.BootstrapConfig.ServiceName == "" {
		return fmt.Errorf("BootstrapConfig.ServiceName cannot be empty when bootstrapping mongo")
	}
	if b.BootstrapConfig.MongoRootConnURL == "" {
		return fmt.Errorf("BootstrapConfig.MongoRootConnURL cannot be empty when bootstrapping mongo")
	}
	if b.BootstrapConfig.MongoClusterName == "" {
		return fmt.Errorf("BootstrapConfig.MongoClusterName cannot be empty when bootstrapping mongo")
	}
	if b.BootstrapConfig.MongoRoleLeaseTTL == "" {
		return fmt.Errorf("BootstrapConfig.MongoRoleLeaseTTL cannot be empty when bootstrapping mongo")
	}
	if b.BootstrapConfig.MongoRoleLeaseMaxTTL == "" {
		return fmt.Errorf("BootstrapConfig.MongoRoleLeaseMaxTTL cannot be empty when bootstrapping mongo")
	}

	b.config.mongoRootConnURL = b.BootstrapConfig.MongoRootConnURL
	b.config.mongoPolicyMountPath = fmt.Sprintf("%s/%s/%s", b.BootstrapConfig.Env, SecretBackendMongo, b.BootstrapConfig.MongoClusterName)
	b.config.mongoPolicyCredsPath = fmt.Sprintf("%s/%s/%s/creds/%s", b.BootstrapConfig.Env, SecretBackendMongo, b.BootstrapConfig.MongoClusterName, b.BootstrapConfig.ServiceName)
	b.config.mongoPolicyRoleCreatePath = fmt.Sprintf("%s/%s/%s/roles/%s", b.BootstrapConfig.Env, SecretBackendMongo, b.BootstrapConfig.MongoClusterName, b.BootstrapConfig.ServiceName)
	b.config.mongoPolicyRoleLeasePath = fmt.Sprintf("%s/%s/%s/config/lease", b.BootstrapConfig.Env, SecretBackendMongo, b.BootstrapConfig.MongoClusterName)
	b.config.mongoPolicyRoleConnPath = fmt.Sprintf("%s/%s/%s/config/connection", b.BootstrapConfig.Env, SecretBackendMongo, b.BootstrapConfig.MongoClusterName)
	b.config.mongoPolicyRoleLeaseTTL = b.BootstrapConfig.MongoRoleLeaseTTL
	b.config.mongoPolicyRoleLeaseMaxTTL = b.BootstrapConfig.MongoRoleLeaseMaxTTL

	return nil
}

func (b *BootstrapClient) configurePostgres() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	if b.BootstrapConfig.Env == "" {
		return fmt.Errorf("BootstrapConfig.Env cannot be empty when bootstrapping postgres")
	}
	if b.BootstrapConfig.ServiceName == "" {
		return fmt.Errorf("BootstrapConfig.ServiceName cannot be empty when bootstrapping postgres")
	}
	if b.BootstrapConfig.PostgresRootConnURL == "" {
		return fmt.Errorf("BootstrapConfig.PostgresRootConnURL cannot be empty when bootstrapping postgres")
	}
	if b.BootstrapConfig.PostgresName == "" {
		return fmt.Errorf("BootstrapConfig.PostgresName cannot be empty when bootstrapping postgres")
	}
	if b.BootstrapConfig.PostgresRoleLeaseTTL == "" {
		return fmt.Errorf("BootstrapConfig.PostgresRoleLeaseTTL cannot be empty when bootstrapping postgres")
	}
	if b.BootstrapConfig.PostgresRoleLeaseMaxTTL == "" {
		return fmt.Errorf("BootstrapConfig.PostgresRoleLeaseMaxTTL cannot be empty when bootstrapping postgres")
	}

	pgRoleTmpl := "CREATE ROLE \"{{name}}\" WITH CREATEDB LOGIN PASSWORD '{{password}}' IN ROLE \"%s\" VALID UNTIL '{{expiration}}';"
	b.config.postgresRootConnURL = b.BootstrapConfig.PostgresRootConnURL
	b.config.postgresPolicyMountPath = fmt.Sprintf("%s/%s/%s", b.BootstrapConfig.Env, SecretBackendPostgres, b.BootstrapConfig.PostgresName)
	b.config.postgresPolicyCredsPath = fmt.Sprintf("%s/%s/%s/creds/%s", b.BootstrapConfig.Env, SecretBackendPostgres, b.BootstrapConfig.PostgresName, b.BootstrapConfig.ServiceName)
	b.config.postgresPolicyRoleCreatePath = fmt.Sprintf("%s/%s/%s/roles/%s", b.BootstrapConfig.Env, SecretBackendPostgres, b.BootstrapConfig.PostgresName, b.BootstrapConfig.ServiceName)
	b.config.postgresPolicyRoleCreateSQL = fmt.Sprintf(pgRoleTmpl, b.BootstrapConfig.ServiceName)
	b.config.postgresPolicyRoleLeasePath = fmt.Sprintf("%s/%s/%s/config/lease", b.BootstrapConfig.Env, SecretBackendPostgres, b.BootstrapConfig.PostgresName)
	b.config.postgresPolicyRoleConnPath = fmt.Sprintf("%s/%s/%s/config/connection", b.BootstrapConfig.Env, SecretBackendPostgres, b.BootstrapConfig.PostgresName)
	b.config.postgresPolicyRoleLeaseTTL = b.BootstrapConfig.PostgresRoleLeaseTTL
	b.config.postgresPolicyRoleLeaseMaxTTL = b.BootstrapConfig.PostgresRoleLeaseMaxTTL

	return nil
}

func (b *BootstrapClient) configureAppRoleAuth() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	if b.BootstrapConfig.Env == "" {
		return fmt.Errorf("BootstrapConfig.Env cannot be empty when bootstrapping app role auth")
	}
	if b.BootstrapConfig.ServiceName == "" {
		return fmt.Errorf("BootstrapConfig.ServiceName cannot be empty when bootstrapping app role auth")
	}
	if b.BootstrapConfig.AuthTokenTTL == "" {
		return fmt.Errorf("BootstrapConfig.AuthTokenTTL cannot be empty when bootstrapping app role auth")
	}
	if b.BootstrapConfig.AuthTokenMaxTTL == "" {
		return fmt.Errorf("BootstrapConfig.AuthTokenMaxTTL cannot be empty when bootstrapping app role auth")
	}
	b.config.authPath = fmt.Sprintf("%s/approle", b.BootstrapConfig.Env)
	b.config.authType = "approle"
	b.config.authDesc = ""
	b.config.authPolicyPath = fmt.Sprintf("auth/%s/approle/role/%s", b.BootstrapConfig.Env, b.BootstrapConfig.ServiceName)
	b.config.authPolicyRoleIDPath = fmt.Sprintf("auth/%s/approle/role/%s/role-id", b.BootstrapConfig.Env, b.BootstrapConfig.ServiceName)
	b.config.authPolicySecretIDPath = fmt.Sprintf("auth/%s/approle/role/%s/secret-id", b.BootstrapConfig.Env, b.BootstrapConfig.ServiceName)
	b.config.authPolicyData = map[string]interface{}{
		"token_ttl":     b.BootstrapConfig.AuthTokenTTL,
		"token_max_ttl": b.BootstrapConfig.AuthTokenMaxTTL,
		"policies":      b.config.servicePolicyName,
	}

	return nil
}

func (b *BootstrapClient) configureTLSCert() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	if b.BootstrapConfig.Env == "" {
		return fmt.Errorf("BootstrapConfig.Env cannot be empty when bootstrapping TLS cert")
	}
	if b.BootstrapConfig.ServiceName == "" {
		return fmt.Errorf("BootstrapConfig.ServiceName cannot be empty when bootstrapping TLS cert")
	}
	if b.BootstrapConfig.TLSCertName == "" {
		return fmt.Errorf("BootstrapConfig.TLSCertName cannot be empty when bootstrapping TLS cert")
	}
	if b.BootstrapConfig.TLSCertData == "" {
		return fmt.Errorf("BootstrapConfig.TLSCertData cannot be empty when bootstrapping TLS cert")
	}
	if b.BootstrapConfig.TLSCertKeyData == "" {
		return fmt.Errorf("BootstrapConfig.TLSCertKeyData cannot be empty when bootstrapping TLS cert")
	}
	if b.BootstrapConfig.TLSCertLeaseTTL == "" {
		return fmt.Errorf("BootstrapConfig.TLSCertLeaseTTL cannot be empty when bootstrapping TLS cert")
	}
	b.config.tlsCertPolicyCertPath = fmt.Sprintf("%s/%s/%s/certificates/%s/cert", b.BootstrapConfig.Env, SecretBackendGeneric, b.BootstrapConfig.ServiceName, b.BootstrapConfig.TLSCertName)
	b.config.tlsCertPolicyKeyPath = fmt.Sprintf("%s/%s/%s/certificates/%s/key", b.BootstrapConfig.Env, SecretBackendGeneric, b.BootstrapConfig.ServiceName, b.BootstrapConfig.TLSCertName)
	b.config.tlsCertData = b.BootstrapConfig.TLSCertData
	b.config.tlsCertKeyData = b.BootstrapConfig.TLSCertKeyData
	b.config.tlsCertPolicyLeaseTTL = b.BootstrapConfig.TLSCertLeaseTTL
	return nil
}

func (b *BootstrapClient) configureGeneric() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	if b.BootstrapConfig.GenericData == "" {
		return fmt.Errorf("BootstrapConfig.GenericData cannot be empty when bootstrapping generic data")
	}
	return json.Unmarshal([]byte(b.BootstrapConfig.GenericData), &b.config.genericData)

}

func (b *BootstrapClient) configureTransit() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	b.config.transitMountPath = fmt.Sprintf("%s/%s", b.BootstrapConfig.Env, SecretBackendTransit)
	if b.BootstrapConfig.TransitData == "" {
		return fmt.Errorf("BootstrapConfig.TransitData cannot be empty when bootstrapping transit data")
	}
	return json.Unmarshal([]byte(b.BootstrapConfig.TransitData), &b.config.transitData)

}

func (b *BootstrapClient) configureAWS() error {
	if err := b.configurePolicyName(); err != nil {
		return err
	}
	if b.BootstrapConfig.Env == "" {
		return fmt.Errorf("BootstrapConfig.Env cannot be empty when bootstrapping AWS")
	}

	// Required always
	if b.BootstrapConfig.AWSRegion == "" {
		return fmt.Errorf("BootstrapConfig.AWSRegion cannot be empty when bootstrapping AWS")
	}
	if b.BootstrapConfig.AWSRootAccessKey == "" {
		return fmt.Errorf("BootstrapConfig.AWSRootAccessKey cannot be empty when bootstrapping AWS")
	}
	if b.BootstrapConfig.AWSRootSecretKey == "" {
		return fmt.Errorf("BootstrapConfig.AWSRootSecretKey cannot be empty when bootstrapping AWS")
	}

	b.config.awsMountPath = fmt.Sprintf("%s/%s", b.BootstrapConfig.Env, SecretBackendAWS)
	b.config.awsConfigRootPath = fmt.Sprintf("%s/%s/config/root", b.BootstrapConfig.Env, SecretBackendAWS)
	b.config.awsConfigRootData = map[string]interface{}{
		"region":     b.BootstrapConfig.AWSRegion,
		"access_key": b.BootstrapConfig.AWSRootAccessKey,
		"secret_key": b.BootstrapConfig.AWSRootSecretKey,
	}

	if b.BootstrapConfig.AWSRoleName != "" {
		if b.BootstrapConfig.AWSRolePolicyJSON == "" {
			return fmt.Errorf("BootstrapConfig.AWSRolePolicyJSON cannot be empty when BootstrapConfig.AWSRoleName is present")
		}

		b.config.awsPolicyRolesPath = fmt.Sprintf("%s/%s/roles/%s", b.BootstrapConfig.Env, SecretBackendAWS, b.BootstrapConfig.AWSRoleName)
		b.config.awsPolicyRolesPolicy = b.BootstrapConfig.AWSRolePolicyJSON
		b.config.awsPolicyCredsPath = fmt.Sprintf("%s/%s/sts/%s", b.BootstrapConfig.Env, SecretBackendAWS, b.BootstrapConfig.AWSRoleName)
		return nil
	}

	if b.BootstrapConfig.AWSAssumedRole == "" {
		return fmt.Errorf("BootstrapConfig.AWSAssumeRole cannot be empty when bootstrapping AWS")
	}

	assumedRole := strings.SplitN(b.BootstrapConfig.AWSAssumedRole, "|", 2)
	// assumedRole is [name, arn]
	b.config.awsAssumedRoleRolesPath = fmt.Sprintf("%s/%s/roles/%s", b.BootstrapConfig.Env, SecretBackendAWS, assumedRole[0])
	b.config.awsAssumedRoleCredsPath = fmt.Sprintf("%s/%s/sts/%s", b.BootstrapConfig.Env, SecretBackendAWS, assumedRole[0])
	b.config.awsAssumedRoleArn = assumedRole[1]
	return nil
}

func (b *BootstrapClient) usingVaultRootToken(cb func() error) error {
	if b.BootstrapConfig.VaultRootToken == "" {
		panic(fmt.Errorf("BootstrapConfig.VaultRootToken cannot be empty"))
	}
	prevToken := b.VaultClient.Token()
	b.VaultClient.SetToken(b.BootstrapConfig.VaultRootToken)
	defer b.VaultClient.SetToken(prevToken)
	return cb()
}

func (b *BootstrapClient) must(err error) {
	if err != nil {
		panic(err)
	}
}
