package database

import "time"

// Option defines an database client option
type Option func(*VaultSQLClient)

// OptionCredentialsVaultPath configures the vault secret path on the client
func OptionCredentialsVaultPath(path string) Option {
	return func(c *VaultSQLClient) {
		c.vaultSecretPath = path
	}
}

// OptionSQLConnectionTemplate configures sql connection template on the client.
// This is used to resolve the actual connection url with username and password fetched from vault.
func OptionSQLConnectionTemplate(tmpl string) Option {
	return func(c *VaultSQLClient) {
		c.sqlConnTemplate = tmpl
	}
}

// OptionSQLDriver configures the driver on sql client
func OptionSQLDriver(driver string) Option {
	return func(c *VaultSQLClient) {
		c.sqlDriver = driver
	}
}

// OptionSQLConnectionPingInterval configures sql connection ping interval on the client
func OptionSQLConnectionPingInterval(interval time.Duration) Option {
	return func(c *VaultSQLClient) {
		c.pingInterval = interval
	}
}
