package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	// ImageBuildTime contains the build time of the container with currently installed version
	ImageBuildTime = "No build time provided"
	// ServiceVersion contains currently installed version (based on Git)
	ServiceVersion = "No version provided"
)

// Config contains the runtime configuration of the service
type Config struct {
	Env            string
	ServiceName    string
	GRPCPort       int
	HTTPStatusPort int

	VaultPostgresCredsPath     string
	PostgresConnectionTemplate string
	PostgresRootRole           string
	PostgresReadOnlyRole       string

	FlowdockToken string
}

// FromEnv reads the service settings from environment variables
func FromEnv() (config *Config, err error) {

	defer func() {
		if r := recover(); r != nil {
			config = nil
			switch r.(type) {
			case string:
				err = fmt.Errorf("%s", r)
			case error:
				err = r.(error)
			default:
				return
			}
		}
	}()

	config = &Config{
		Env:                        mustRead(EnvEnv),
		ServiceName:                mustRead(EnvServiceName),
		GRPCPort:                   mustReadInt(EnvGRPCPort),
		HTTPStatusPort:             mustReadInt(EnvStatusPort),
		VaultPostgresCredsPath:     mustRead(EnvVaultPostgresCredsPath),
		PostgresConnectionTemplate: mustRead(EnvPostgresConnectionTemplate),
		PostgresRootRole:           mustRead(EnvPostgresRootRole),
		PostgresReadOnlyRole:       mustRead(EnvPostgresReadOnlyRole),
		FlowdockToken:              os.Getenv(EnvFlowdockToken),
	}

	return
}

func mustRead(envVar string) string {
	s := os.Getenv(envVar)
	if s == "" {
		panic(fmt.Errorf("missing a mandatory environment variable %s", envVar))
	}
	return s
}

func mustReadInt(envVar string) int {
	s := mustRead(envVar)
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("invalid integer %s for environment variable %s", s, envVar))
	}
	return i
}
