package ip2country

import (
	"ip2country-api/internal/config"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		config      config.BackendConfig
		expectError bool
	}{
		{
			name: "CSV Service",
			config: config.BackendConfig{
				Type:    "csv",
				CSVPath: createTestCSVFile(t),
			},
			expectError: false,
		},
		{
			name: "MongoDB Service",
			config: config.BackendConfig{
				Type:     "mongodb",
				MongoURI: "mongodb://localhost:27017",
			},
			expectError: false,
		},
		{
			name: "Redis Service",
			config: config.BackendConfig{
				Type:      "redis",
				RedisAddr: "localhost:6379",
			},
			expectError: false,
		},
		{
			name: "Unsupported Service",
			config: config.BackendConfig{
				Type: "unsupported",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service, err := NewService(tc.config)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if service == nil {
					t.Error("Expected non-nil service")
				}
			}

			// Cleanup test file if it was created
			if tc.config.Type == "csv" {
				os.Remove(tc.config.CSVPath)
			}
		})
	}
}

// Helper function to create a test CSV file
func createTestCSVFile(t *testing.T) string {
	testFile := "test_service.csv"
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.WriteString("192.168.1.1,New York,USA\n")
	f.Close()
	return testFile
}
