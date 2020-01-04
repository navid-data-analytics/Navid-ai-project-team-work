package postgres

import (
	"errors"
	"os"
)

// Env variable names
const (
	EnvConnectionTemplate = "POSTGRES_CONN_TMPL"
)

// Errors
var (
	ErrEmptyConnectionTemplate = errors.New("Connection template cannot be empty")
)

// Options implements all the possible/required custom options a postgres client requires
type Options struct {
	ConnectionTemplate string // connection template to use when connecting with vault based credentials
}

// Validate returns an error if one of the options is invalid
func (o *Options) Validate() error {
	if o.ConnectionTemplate == "" {
		return ErrEmptyConnectionTemplate
	}

	return nil
}

// OptionsFromEnv reads the options based on environment variables
func OptionsFromEnv() (*Options, error) {
	opts := &Options{
		ConnectionTemplate: os.Getenv(EnvConnectionTemplate),
	}

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return opts, nil
}
