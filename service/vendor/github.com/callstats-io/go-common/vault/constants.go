package vault

import (
	"errors"
	"fmt"
	"time"
)

const (
	// MinLeaseExpireDuration is the minimum duration for leases that should be used.
	// If the duration left on a renewed lease is less than this it should be re-fetched as a new secret lease.
	MinLeaseExpireDuration = 5 * time.Minute
)

// AWS Secret data access keys
const (
	SecretDataKeyUsername      = "username"
	SecretDataKeyPassword      = "password"
	SecretDataKeyData          = "data"
	SecretDataKeyAccessKey     = "access_key"
	SecretDataKeySecretKey     = "secret_key"
	SecretDataKeySecurityToken = "security_token"
)

// Log messages
const (
	LogErrFailedReadCredentials  = "Failed to read auth credentials"
	LogErrInvalidCredentials     = "Invalid credentials"
	LogErrFailedRenewAuth        = "Failed to renew auth token"
	LogErrFailedAuth             = "Failed to get auth token from vault"
	LogErrAutoRenewFetchFailure  = "Failed to fetch secret on automatic renew"
	LogAutoRenewFetchSuccess     = "Fetched secret on automatic renew"
	LogAuthSuccessMsg            = "Authenticated successfully"
	LogAuthRenewSuccessMsg       = "Renewed authentication token successfully"
	LogPostgresCredsReadSuccess  = "Read postgres creds from vault"
	LogPostgresCredsReadFailure  = "Failed to read postgres creds"
	LogPostgresCredsRenewSuccess = "Renewed postgres creds from vault"
	LogPostgresCredsRenewFailure = "Failed to renew postgres creds"
	LogMongoCredsReadSuccess     = "Read mongodb creds from vault"
	LogMongoCredsReadFailure     = "Failed to read mongodb creds"
	LogMongoCredsRenewSuccess    = "Renewed mongodb creds from vault"
	LogMongoCredsRenewFailure    = "Failed to renew mongodb creds"
	LogMsgErrFailedSecretRenew   = "Failed to renew secret"
	LogTLSCertKeyReadSuccess     = "Read tls cert key from vault"
	LogTLSCertKeyReadFailure     = "Failed to read tls cert key"
	LogTLSCertKeyRenewSuccess    = "Renewed tls cert key from vault"
	LogTLSCertKeyRenewFailure    = "Failed to renew tls cert key"
	LogTLSCertReadSuccess        = "Read tls cert from vault"
	LogTLSCertReadFailure        = "Failed to read tls cert"
	LogTLSCertRenewSuccess       = "Renewed tls cert from vault"
	LogTLSCertRenewFailure       = "Failed to renew tls cert"
	LogAWSCredsReadSuccess       = "Read aws creds from vault"
	LogAWSCredsReadFailure       = "Failed to read aws creds"
	LogAWSCredsRenewSuccess      = "Renewed aws creds from vault"
	LogAWSCredsRenewFailure      = "Failed to renew aws creds"
)

// Env variables
const (
	EnvEnv                    = "ENV"                       // env key for env
	EnvVaultAppRoleCreds      = "VAULT_AUTHCREDENTIALS"     // env key from where vault creds are resolved, matches dashboard
	EnvVaultApproleAuthPath   = "VAULT_APPROLE_AUTH_PATH"   // env key for reading approle auth path
	EnvVaultEnableMongo       = "VAULT_ENABLE_MONGO"        // env key for enabling mongo credentials reading
	EnvVaultMongoCredsPath    = "VAULT_MONGO_CREDS_PATH"    // env key for reading mongo credentials from vault
	EnvVaultEnablePostgres    = "VAULT_ENABLE_POSTGRES"     // env key for enabling postgres credentials reading
	EnvVaultPostgresCredsPath = "VAULT_POSTGRES_CREDS_PATH" // env key for reading postgres credentials from vault
	EnvVaultEnableTLSCert     = "VAULT_ENABLE_TLS_CERT"     // env key for enabling TLS key reading
	EnvVaultTLSCertPath       = "VAULT_TLS_CERT_PATH"       // env key for reading tls cert from vault
	EnvVaultTLSCertKeyPath    = "VAULT_TLS_CERT_KEY_PATH"   // env key for reading tls cert key from vault
	EnvVaultEnableAWS         = "VAULT_ENABLE_AWS"          // env key for enabling AWS creds reading
	EnvVaultAWSCredsPath      = "VAULT_AWS_CREDS_PATH"      // env key for reading aws creds from vault
)

// Log field names
const (
	LogKeyAccessor        = "accessor"
	LogKeyLeaseID         = "leaseID"
	LogKeyLeaseDuration   = "leaseDuration"
	LogKeyCreateTime      = "createdAt"
	LogKeyRenewTime       = "renewAt"
	LogKeyExpireTime      = "expireAt"
	LogKeyBackoff         = "backoff"
	LogKeyCurrentSecretID = "currentSecretID"
	LogKeyPrevSecretID    = "prevSecretID"
)

// Errors
var (
	ErrEmptyEnv                         = errors.New("Env cannot be empty")
	ErrEmptyAppRoleAuthPath             = errors.New("AppRoleAuthPath cannot be empty")
	ErrEmptyMongoCredsReadPath          = errEmptyRequired("MongoCredsReadPath", "EnableMongo")
	ErrEmptyPostgresCredsReadPath       = errEmptyRequired("PostgresCredsReadPath", "EnablePostgres")
	ErrEmptyTLSCertReadPath             = errEmptyRequired("TLSCertReadPath", "EnableTLSCert")
	ErrEmptyTLSCertKeyReadPath          = errEmptyRequired("TLSCertKeyReadPath", "EnableTLSCert")
	ErrEmptyAWSCredsReadPath            = errEmptyRequired("AWSCredsReadPath", "EnableAWS")
	ErrMongoDisabled                    = errDisabled("Mongo", "EnableMongo")
	ErrPostgresDisabled                 = errDisabled("Postgres", "EnablePostgres")
	ErrTLSCertDisabled                  = errDisabled("TLSCert", "EnableTLSCert")
	ErrAWSDisabled                      = errDisabled("AWS", "EnableAWS")
	ErrEmptyEnvAppRoleCreds             = errors.New("ENV VAULT_AUTHCREDENTIALS is empty")
	ErrInvalidRoleID                    = errors.New("Invalid app role authentication role id")
	ErrInvalidSecretID                  = errors.New("Invalid app role authentication secret id")
	ErrVaultSecretNotSet                = errors.New("Cannot renew non-existent vault secret")
	ErrVaultSecretNotRenewable          = errors.New("Vault secret is not renewable")
	ErrVaultSecretRenewFailed           = errors.New("Vault renew request failed")
	ErrVaultMaxLeaseExceeded            = errors.New("Vault secret max lease time exceeded")
	ErrInvalidSecretTLSCertFormat       = errInvalidFormat("tls cert")
	ErrInvalidSecretTLSCertKeyFormat    = errInvalidFormat("tls cert key")
	ErrEmptySecretTLSCertData           = errEmpty("tls cert", SecretDataKeyData)
	ErrEmptySecretTLSCertKeyData        = errEmpty("tls cert key", SecretDataKeyData)
	ErrInvalidSecretUsernameFormat      = errInvalidFormat("username")
	ErrInvalidSecretPasswordFormat      = errInvalidFormat("password")
	ErrInvalidSecretAccessKeyFormat     = errInvalidFormat("access_key")
	ErrInvalidSecretSecretKeyFormat     = errInvalidFormat("secret_key")
	ErrInvalidSecretSecurityTokenFormat = errInvalidFormat("security_token")
	ErrEmptySecretUsernameData          = errEmpty("username", SecretDataKeyUsername)
	ErrEmptySecretPasswordData          = errEmpty("password", SecretDataKeyPassword)
	ErrEmptySecretAccessKeyData         = errEmpty("access_key", SecretDataKeyAccessKey)
	ErrEmptySecretSecretKeyData         = errEmpty("secret_key", SecretDataKeySecretKey)
	ErrEmptySecretSecurityTokenData     = errEmpty("security_token", SecretDataKeySecurityToken)
	ErrUnauthenticated                  = errors.New("Unauthenticated")
	ErrInvalidClientState               = errors.New("Client is incorrectly setup")
	ErrAuthSecretChanged                = errors.New("Auth secret changed")
	ErrSecretNotFound                   = errors.New("Secret was not found (Vault likely returned 404)")
)

func errDisabled(name, required string) error {
	return fmt.Errorf("%s features are disabled. Set options.%s to true and configure the required options", name, required)
}

func errEmptyRequired(name, required string) error {
	return fmt.Errorf("%s cannot be empty when %s is true", name, required)
}

func errEmpty(name, key string) error {
	return fmt.Errorf("Invalid %s data: \"%s\" not present in response data", name, key)
}

func errInvalidFormat(name string) error {
	return fmt.Errorf("Invalid %s format: not a string", name)
}
