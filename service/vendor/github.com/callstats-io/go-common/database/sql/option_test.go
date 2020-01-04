package database

import (
	"testing"
	"time"
)

func TestOptionSQLConnectionPingInterval(t *testing.T) {
	client := &VaultSQLClient{}
	opt := OptionSQLConnectionPingInterval(time.Second)

	// call option
	opt(client)

	if client.pingInterval != time.Second {
		t.Errorf("Expected OptionSQLConnectionPingInterval to set interval to %d got %d", time.Second, client.pingInterval)
	}
}

func TestOptionSQLDriver(t *testing.T) {
	driver := "testdriver"
	client := &VaultSQLClient{}
	opt := OptionSQLDriver(driver)

	// call option
	opt(client)

	if client.sqlDriver != driver {
		t.Errorf("Expected OptionSQLDriver to set SQL driver to %s got %s", driver, client.sqlDriver)
	}
}

func TestOptionCredentialsVaultPath(t *testing.T) {
	path := "/testpath"
	client := &VaultSQLClient{}
	opt := OptionCredentialsVaultPath(path)

	// call option
	opt(client)

	if client.vaultSecretPath != path {
		t.Errorf("Expected OptionCredentialsVaultPath to set Vault path to %s got %s", path, client.vaultSecretPath)
	}
}

func TestOptionSQLConnectionTemplate(t *testing.T) {
	tmpl := "/testtmpl"
	client := &VaultSQLClient{}
	opt := OptionSQLConnectionTemplate(tmpl)

	// call option
	opt(client)

	if client.sqlConnTemplate != tmpl {
		t.Errorf("Expected OptionSQLConnectionTemplate to set SQL connection template to %s got %s", tmpl, client.sqlConnTemplate)
	}
}
