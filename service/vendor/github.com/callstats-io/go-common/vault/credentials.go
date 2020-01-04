package vault

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	loadCredsFileOffset       = len("file:")
	loadCredsKubernetesOffset = len("kubernetes:")
	kubernetesRoleFileName    = "role_id"
	kubernetesSecretFileName  = "secret_id"
)

// AppRoleCredentials contains the credentials required for app role authentication to vault
type AppRoleCredentials struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// Validate checks that the role id and secret id are not empty strings
func (arc *AppRoleCredentials) Validate() error {
	if arc.RoleID == "" {
		return ErrInvalidRoleID
	}
	if arc.SecretID == "" {
		return ErrInvalidSecretID
	}
	return nil
}

// ReadEnvironment tries to read the credentials based on the VAULT_AUTHCREDENTIALS environment variable
// This can have one of three forms:
// 1) JSON with secret and role id as string
// 2) file://path/to/file pointing to a file with json of secret and role
// 3) kubernetes://path/to/dir with two files: secret_id and role_id each containing a string for the key
func (arc *AppRoleCredentials) ReadEnvironment() error {
	loc := os.Getenv(EnvVaultAppRoleCreds)
	if loc == "" {
		return ErrEmptyEnvAppRoleCreds
	}

	if strings.HasPrefix(loc, "file:") {
		return arc.ReadFile(loc[loadCredsFileOffset:])
	}

	if strings.HasPrefix(loc, "kubernetes:") {
		return arc.ReadKubernetes(loc[loadCredsKubernetesOffset:])
	}

	return arc.ReadJSON(loc)
}

// ReadFile parses the credentials from a file containing json format of the credentials and validates them
// @params file path/to/file with secret_id + role_id as json
func (arc *AppRoleCredentials) ReadFile(filepath string) error {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(contents, arc); err != nil {
		return err
	}

	return arc.Validate()
}

// ReadKubernetes parses the credentials two files containing the credentials and validates them
// @params dirpath path/to/dir with two files: secret_id and role_id each containing a string for the key
func (arc *AppRoleCredentials) ReadKubernetes(dirpath string) error {
	// get role from file pointed to by kubernetes
	role, err := ioutil.ReadFile(path.Join(dirpath, kubernetesRoleFileName))
	if err != nil {
		return err
	}

	// get secret from file pointed to by kubernetes
	secret, err := ioutil.ReadFile(path.Join(dirpath, kubernetesSecretFileName))
	if err != nil {
		return err
	}

	arc.RoleID = string(role)
	arc.SecretID = string(secret)

	return arc.Validate()
}

// ReadJSON parses the credentials from a json string and validates them
// @params str json object string containing the secret_id and role_id
// Example: {"secret_id":"abc","role_id":"def"}
func (arc *AppRoleCredentials) ReadJSON(str string) error {
	if err := json.Unmarshal([]byte(str), arc); err != nil {
		return err
	}

	return arc.Validate()
}

// Map returns the credentials as a Go map
func (arc *AppRoleCredentials) Map() map[string]interface{} {
	return map[string]interface{}{
		"role_id":   arc.RoleID,
		"secret_id": arc.SecretID,
	}
}

// UserPassCredentials contains Username and password
type UserPassCredentials struct {
	User     string
	Password string
}

// STSCredentials contains an access key, secret key and security token
type STSCredentials struct {
	AccessKey     string
	SecretKey     string
	SecurityToken string
}
