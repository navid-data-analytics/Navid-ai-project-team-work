package config_test

import (
	"fmt"
	"os"

	"testing"

	"github.com/callstats-io/ai-decision/service/src/config"
	"github.com/stretchr/testify/require"
)

func TestConfigFromEnv(t *testing.T) {
	// validate correct test env
	t.Run("valid test environment", func(t *testing.T) {
		assert := require.New(t)
		_, err := config.FromEnv()
		assert.Nil(err)
	})

	// Default tests
	type envTestCase struct {
		EnvVariableName          string
		EnvVariableInvalidValues []string
		EnvVariableValidValues   []string
	}
	testCases := []envTestCase{
		envTestCase{
			EnvVariableName:          config.EnvEnv,
			EnvVariableInvalidValues: []string{""},
			EnvVariableValidValues:   []string{"test"},
		},
		envTestCase{
			EnvVariableName:          config.EnvServiceName,
			EnvVariableInvalidValues: []string{""},
			EnvVariableValidValues:   []string{"myawesomename"},
		},
		envTestCase{
			EnvVariableName:          config.EnvGRPCPort,
			EnvVariableInvalidValues: []string{"unknown", ""},
			EnvVariableValidValues:   []string{"12345"},
		},
		envTestCase{
			EnvVariableName:          config.EnvStatusPort,
			EnvVariableInvalidValues: []string{"unknown", ""},
			EnvVariableValidValues:   []string{"12345"},
		},
		envTestCase{
			EnvVariableName:          config.EnvVaultPostgresCredsPath,
			EnvVariableInvalidValues: []string{""},
			EnvVariableValidValues:   []string{"anything"},
		},
		envTestCase{
			EnvVariableName:          config.EnvPostgresConnectionTemplate,
			EnvVariableInvalidValues: []string{""},
			EnvVariableValidValues:   []string{"anything"},
		},
		envTestCase{
			EnvVariableName:          config.EnvPostgresRootRole,
			EnvVariableInvalidValues: []string{""},
			EnvVariableValidValues:   []string{"anything"},
		},
		envTestCase{
			EnvVariableName:          config.EnvPostgresReadOnlyRole,
			EnvVariableInvalidValues: []string{""},
			EnvVariableValidValues:   []string{"anything"},
		},
	}
	for idx := range testCases {
		testCase := testCases[idx]
		if testCase.EnvVariableInvalidValues != nil {
			t.Run(fmt.Sprintf("should fail if %s is invalid", testCase.EnvVariableName), func(t *testing.T) {
				assert := require.New(t)
				prev := os.Getenv(testCase.EnvVariableName)
				defer os.Setenv(testCase.EnvVariableName, prev)
				for _, val := range testCase.EnvVariableInvalidValues {
					os.Setenv(testCase.EnvVariableName, val)
					_, err := config.FromEnv()
					assert.NotNil(err)
				}
			})
		}
		if testCase.EnvVariableValidValues != nil {
			t.Run(fmt.Sprintf("should succeed if %s is valid", testCase.EnvVariableName), func(t *testing.T) {
				assert := require.New(t)
				prev := os.Getenv(testCase.EnvVariableName)
				defer os.Setenv(testCase.EnvVariableName, prev)
				for _, val := range testCase.EnvVariableValidValues {
					os.Setenv(testCase.EnvVariableName, val)
					_, err := config.FromEnv()
					assert.Nil(err)
				}
			})
		}
	}
}
