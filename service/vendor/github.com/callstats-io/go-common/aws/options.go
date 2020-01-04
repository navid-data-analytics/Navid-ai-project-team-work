package aws

import (
	"net/http"
	"os"
)

// Options stores AWS related options
type Options struct {
	AWSRegion    string
	MaxRetries   int
	HTTPClient   *http.Client
	DebugLogging bool
}

// OptionsFromEnv reads environment variables to *Options
func OptionsFromEnv() (*Options, error) {
	awsRegion := os.Getenv(EnvAWSRegion)
	if awsRegion == "" {
		awsRegion = DefaultAWSRegion
	}
	opts := &Options{
		AWSRegion:  awsRegion,
		MaxRetries: 3,
		HTTPClient: http.DefaultClient,
	}
	logLevel := os.Getenv(EnvLogLevel)
	if logLevel == "DEBUG" {
		opts.DebugLogging = true
	}
	return opts, nil
}

// WithMaxRetries overrides the max number of retries when a Client session is created
func (o *Options) WithMaxRetries(retries int) *Options {
	o.MaxRetries = retries
	return o
}

// WithHTTPClient overrides the default HTTP client
func (o *Options) WithHTTPClient(hc *http.Client) *Options {
	o.HTTPClient = hc
	return o
}

// WithDebugLogging turns on debug logging
func (o *Options) WithDebugLogging() *Options {
	o.DebugLogging = true
	return o
}
