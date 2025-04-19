package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment and restore after test
	origDataPath := os.Getenv("CSV_DATA_PATH")
	origRateLimit := os.Getenv("RATE_LIMIT")
	origPort := os.Getenv("PORT")
	origDBType := os.Getenv("IP2COUNTRY_DB_TYPE")
	origMongoURI := os.Getenv("MONGO_URI")
	origRedisAddr := os.Getenv("REDIS_ADDR")
	defer func() {
		os.Setenv("CSV_DATA_PATH", origDataPath)
		os.Setenv("RATE_LIMIT", origRateLimit)
		os.Setenv("PORT", origPort)
		os.Setenv("IP2COUNTRY_DB_TYPE", origDBType)
		os.Setenv("MONGO_URI", origMongoURI)
		os.Setenv("REDIS_ADDR", origRedisAddr)
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
				IP2Country: BackendConfig{
					Type:      "csv",
					CSVPath:   "data/ip2country.csv",
					MongoURI:  "mongodb://localhost:27017",
					RedisAddr: "localhost:6379",
				},
				RateLimit: 100,
				Port:      8080,
			},
			expectError: false,
		},
		{
			name: "Custom values",
			envVars: map[string]string{
				"CSV_DATA_PATH":      "custom/path.csv",
				"RATE_LIMIT":         "200",
				"PORT":               "9090",
				"IP2COUNTRY_DB_TYPE": "mmdb",
			},
			expectedConfig: &Config{
				IP2Country: BackendConfig{
					Type:      "mmdb",
					CSVPath:   "custom/path.csv",
					MongoURI:  "mongodb://localhost:27017",
					RedisAddr: "localhost:6379",
				},
				RateLimit: 200,
				Port:      9090,
			},
			expectError: false,
		},
		{
			name: "MongoDB configuration",
			envVars: map[string]string{
				"IP2COUNTRY_DB_TYPE": "mongodb",
				"MONGO_URI":          "mongodb://custom-server:27018",
			},
			expectedConfig: &Config{
				IP2Country: BackendConfig{
					Type:      "mongodb",
					CSVPath:   "data/ip2country.csv",
					MongoURI:  "mongodb://custom-server:27018",
					RedisAddr: "localhost:6379",
				},
				RateLimit: 100,
				Port:      8080,
			},
			expectError: false,
		},
		{
			name: "Redis configuration",
			envVars: map[string]string{
				"IP2COUNTRY_DB_TYPE": "redis",
				"REDIS_ADDR":         "custom-redis:6380",
			},
			expectedConfig: &Config{
				IP2Country: BackendConfig{
					Type:      "redis",
					CSVPath:   "data/ip2country.csv",
					MongoURI:  "mongodb://localhost:27017",
					RedisAddr: "custom-redis:6380",
				},
				RateLimit: 100,
				Port:      8080,
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
			os.Unsetenv("CSV_DATA_PATH")
			os.Unsetenv("RATE_LIMIT")
			os.Unsetenv("PORT")
			os.Unsetenv("IP2COUNTRY_DB_TYPE")
			os.Unsetenv("MONGO_URI")
			os.Unsetenv("REDIS_ADDR")

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
			if config.IP2Country.Type != tc.expectedConfig.IP2Country.Type {
				t.Errorf("IP2Country.Type: expected %q, got %q", tc.expectedConfig.IP2Country.Type, config.IP2Country.Type)
			}
			if config.IP2Country.CSVPath != tc.expectedConfig.IP2Country.CSVPath {
				t.Errorf("IP2Country.CSVPath: expected %q, got %q", tc.expectedConfig.IP2Country.CSVPath, config.IP2Country.CSVPath)
			}
			if config.IP2Country.MongoURI != tc.expectedConfig.IP2Country.MongoURI {
				t.Errorf("IP2Country.MongoURI: expected %q, got %q", tc.expectedConfig.IP2Country.MongoURI, config.IP2Country.MongoURI)
			}
			if config.IP2Country.RedisAddr != tc.expectedConfig.IP2Country.RedisAddr {
				t.Errorf("IP2Country.RedisAddr: expected %q, got %q", tc.expectedConfig.IP2Country.RedisAddr, config.IP2Country.RedisAddr)
			}
			if config.RateLimit != tc.expectedConfig.RateLimit {
				t.Errorf("RateLimit: expected %d, got %d", tc.expectedConfig.RateLimit, config.RateLimit)
			}
			if config.Port != tc.expectedConfig.Port {
				t.Errorf("Port: expected %d, got %d", tc.expectedConfig.Port, config.Port)
			}
		})
	}
}
