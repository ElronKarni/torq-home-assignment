package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment and restore after test
	origDataPath := os.Getenv("IP2COUNTRY_DATA_PATH")
	origRateLimit := os.Getenv("RATE_LIMIT")
	origPort := os.Getenv("PORT")
	origDBType := os.Getenv("IP2COUNTRY_DB_TYPE")
	defer func() {
		os.Setenv("IP2COUNTRY_DATA_PATH", origDataPath)
		os.Setenv("RATE_LIMIT", origRateLimit)
		os.Setenv("PORT", origPort)
		os.Setenv("IP2COUNTRY_DB_TYPE", origDBType)
	}()

	testCases := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name:    "Default values",
			envVars: map[string]string{},
			expectedConfig: &Config{
				DataPath:         "data/ip2country.csv",
				RateLimit:        100,
				Port:             8080,
				IP2CountryDBType: "csv",
			},
			expectError: false,
		},
		{
			name: "Custom values",
			envVars: map[string]string{
				"IP2COUNTRY_DATA_PATH": "custom/path.csv",
				"RATE_LIMIT":           "200",
				"PORT":                 "9090",
				"IP2COUNTRY_DB_TYPE":   "mmdb",
			},
			expectedConfig: &Config{
				DataPath:         "custom/path.csv",
				RateLimit:        200,
				Port:             9090,
				IP2CountryDBType: "mmdb",
			},
			expectError: false,
		},
		{
			name: "Invalid RATE_LIMIT",
			envVars: map[string]string{
				"RATE_LIMIT": "not-a-number",
			},
			expectedConfig: nil,
			expectError:    true,
		},
		{
			name: "Invalid PORT",
			envVars: map[string]string{
				"PORT": "not-a-number",
			},
			expectedConfig: nil,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear relevant environment variables first
			os.Unsetenv("IP2COUNTRY_DATA_PATH")
			os.Unsetenv("RATE_LIMIT")
			os.Unsetenv("PORT")
			os.Unsetenv("IP2COUNTRY_DB_TYPE")

			// Set environment variables for this test case
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			// Call the function being tested
			config, err := Load()

			// Check error
			if tc.expectError && err == nil {
				t.Error("Expected an error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Did not expect an error but got: %v", err)
			}

			// If we don't expect a config, we're done
			if tc.expectedConfig == nil {
				return
			}

			// Check config values
			if config.DataPath != tc.expectedConfig.DataPath {
				t.Errorf("DataPath: expected %q, got %q", tc.expectedConfig.DataPath, config.DataPath)
			}
			if config.RateLimit != tc.expectedConfig.RateLimit {
				t.Errorf("RateLimit: expected %d, got %d", tc.expectedConfig.RateLimit, config.RateLimit)
			}
			if config.Port != tc.expectedConfig.Port {
				t.Errorf("Port: expected %d, got %d", tc.expectedConfig.Port, config.Port)
			}
			if config.IP2CountryDBType != tc.expectedConfig.IP2CountryDBType {
				t.Errorf("IP2CountryDBType: expected %q, got %q", tc.expectedConfig.IP2CountryDBType, config.IP2CountryDBType)
			}
		})
	}
}
