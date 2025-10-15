package config

import (
	"os"
)

// SetMockEnv sets environment variables for testing
func SetMockEnv(envs map[string]string) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
}

// UnsetMockEnv unsets environment variables after testing
func UnsetMockEnv(keys []string) {
	for _, k := range keys {
		os.Unsetenv(k)
	}
}
