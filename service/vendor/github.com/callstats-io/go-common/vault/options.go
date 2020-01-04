package vault

import (
	"fmt"
	"os"
	"strconv"
)

// Options implements options required for connecting to vault
type Options struct {
	AppRoleAuthPath string // auth path to use when connecting to vault

	// Mongo related options. EnableMongo has to be set to true if the client uses mongo connections.
	EnableMongo        bool   // enable mongo integration
	MongoCredsReadPath string // mongo credentials read path

	// Postgres related options. EnablePostgres has to be set to true if the client uses postgres connections.
	EnablePostgres        bool   // enable postgres integration
	PostgresCredsReadPath string // postgres credentials read path

	// TLS certificate related options. EnablePostgres has to be set to true if the client uses tls certificates from vault.
	EnableTLSCert      bool
	TLSCertReadPath    string
	TLSCertKeyReadPath string

	// AWS related options. EnableAWS has to be set to true if the client uses AWS connections.
	EnableAWS        bool   // enable AWS integration
	AWSCredsReadPath string // AWS credentials read path
}

// Validate validates the options are usable for mongo connections
func (o *Options) Validate() error {
	if o.AppRoleAuthPath == "" {
		return ErrEmptyAppRoleAuthPath
	}
	if o.EnableMongo {
		if o.MongoCredsReadPath == "" {
			return ErrEmptyMongoCredsReadPath
		}
	}
	if o.EnablePostgres {
		if o.PostgresCredsReadPath == "" {
			return ErrEmptyPostgresCredsReadPath
		}
	}
	if o.EnableTLSCert {
		if o.TLSCertReadPath == "" {
			return ErrEmptyTLSCertReadPath
		}
		if o.TLSCertKeyReadPath == "" {
			return ErrEmptyTLSCertKeyReadPath
		}
	}
	if o.EnableAWS {
		if o.AWSCredsReadPath == "" {
			return ErrEmptyAWSCredsReadPath
		}
	}
	return nil
}

// OptionsFromEnv parses options from the default environment variables.
func OptionsFromEnv() (*Options, error) {
	enableMongo, err := parseEnvBool(EnvVaultEnableMongo)
	if err != nil {
		return nil, formatEnvError(EnvVaultEnableMongo, err)
	}
	enablePostgres, err := parseEnvBool(EnvVaultEnablePostgres)
	if err != nil {
		return nil, formatEnvError(EnvVaultEnablePostgres, err)
	}
	enableTLSCert, err := parseEnvBool(EnvVaultEnableTLSCert)
	if err != nil {
		return nil, formatEnvError(EnvVaultEnableTLSCert, err)
	}
	enableAWS, err := parseEnvBool(EnvVaultEnableAWS)
	if err != nil {
		return nil, formatEnvError(EnvVaultEnableAWS, err)
	}
	appRoleAuthPath := os.Getenv(EnvVaultApproleAuthPath)
	if appRoleAuthPath == "" {
		appRoleAuthPath = "auth/" + os.Getenv(EnvEnv) + "/approle/login"
	}
	opts := &Options{
		AppRoleAuthPath:       appRoleAuthPath,
		EnableMongo:           enableMongo,
		MongoCredsReadPath:    os.Getenv(EnvVaultMongoCredsPath),
		EnablePostgres:        enablePostgres,
		PostgresCredsReadPath: os.Getenv(EnvVaultPostgresCredsPath),
		EnableTLSCert:         enableTLSCert,
		TLSCertReadPath:       os.Getenv(EnvVaultTLSCertPath),
		TLSCertKeyReadPath:    os.Getenv(EnvVaultTLSCertKeyPath),
		EnableAWS:             enableAWS,
		AWSCredsReadPath:      os.Getenv(EnvVaultAWSCredsPath),
	}
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return opts, nil
}

func parseEnvBool(key string) (bool, error) {
	val := os.Getenv(key)
	if val == "" {
		return false, nil
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return parsed, nil
}

func formatEnvError(key string, err error) error {
	return fmt.Errorf("Failed to parse %s, error: %s", key, err)
}
