package render

import (
	"os"
	"testing"
)

func TestGetEnvInt(t *testing.T) {
	// Save original environment variables to restore later
	originalEnv := make(map[string]string)
	for _, env := range os.Environ() {
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				key := env[:i]
				value := env[i+1:]
				originalEnv[key] = value
				break
			}
		}
	}
	
	// Restore environment variables after the test
	defer func() {
		for key := range originalEnv {
			os.Unsetenv(key)
		}
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()
	
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "environment variable not set",
			key:          "TEST_ENV_VAR_NOT_SET",
			envValue:     "",
			defaultValue: 42,
			expected:     42,
		},
		{
			name:         "environment variable set to valid integer",
			key:          "TEST_ENV_VAR_VALID",
			envValue:     "123",
			defaultValue: 42,
			expected:     123,
		},
		{
			name:         "environment variable set to invalid integer",
			key:          "TEST_ENV_VAR_INVALID",
			envValue:     "not-an-integer",
			defaultValue: 42,
			expected:     42,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unset any existing value
			os.Unsetenv(tt.key)
			
			// Set the environment variable if needed
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}
			
			result := getEnvInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvInt(%q, %d) = %d, expected %d", 
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}