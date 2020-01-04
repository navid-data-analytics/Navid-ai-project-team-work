package mongo

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Env variable names
const (
	EnvConnectionTemplate = "MONGO_CONN_TMPL"
	EnvDialTimeout        = "MONGO_DIAL_TIMEOUT"
)

// Errors
var (
	ErrEmptyConnectionTemplate = errors.New("Connection template cannot be empty")
)

// Options implements all the possible/required custom options a mongo client requires
type Options struct {
	ConnectionTemplate string        // connection template to use when connecting with vault based credentials
	DialTimeout        time.Duration // connection dial timeout to use when connecting
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
	rawDialTimeout := os.Getenv(EnvDialTimeout)
	dialTimeout, err := time.ParseDuration(rawDialTimeout)
	if rawDialTimeout != "" && err != nil {
		return nil, fmt.Errorf("Invalid value for %s, error: %s", EnvDialTimeout, err)
	}
	if dialTimeout == 0 {
		dialTimeout = 5 * time.Second // default timeout
	}
	o := &Options{
		ConnectionTemplate: os.Getenv(EnvConnectionTemplate),
		DialTimeout:        dialTimeout,
	}

	if err := o.Validate(); err != nil {
		return nil, err
	}

	return o, nil
}
