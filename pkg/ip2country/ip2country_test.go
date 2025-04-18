package ip2country

import (
	"os"
	"testing"
)

func TestNewCSVService(t *testing.T) {
	// Create a temporary test CSV file
	testCSVContent := "1.1.1.1,Sydney,Australia\n8.8.8.8,Mountain View,United States"
	tempFile, err := os.CreateTemp("", "ip2country_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(testCSVContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test valid initialization
	service, err := NewCSVService(tempFile.Name())
	if err != nil {
		t.Errorf("NewCSVService failed with valid file: %v", err)
	}
	if service == nil {
		t.Error("NewCSVService returned nil service with valid file")
	}

	// Test initialization with non-existent file
	service, err = NewCSVService("nonexistent_file.csv")
	if err == nil {
		t.Error("NewCSVService should fail with non-existent file")
	}
	if service != nil {
		t.Error("NewCSVService should return nil service with non-existent file")
	}
}

func TestLoadData(t *testing.T) {
	// Test with valid CSV file
	testCSVContent := "1.1.1.1,Sydney,Australia\n8.8.8.8,Mountain View,United States"
	tempFile, err := os.CreateTemp("", "ip2country_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(testCSVContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	service := &CSVService{
		filePath: tempFile.Name(),
		data:     make(map[string]*Result),
	}

	// Test successful data loading
	err = service.loadData()
	if err != nil {
		t.Errorf("loadData failed with valid file: %v", err)
	}

	// Verify data was loaded correctly
	if len(service.data) != 2 {
		t.Errorf("Expected 2 entries in data map, got %d", len(service.data))
	}

	// Check specific entries
	result, exists := service.data["1.1.1.1"]
	if !exists {
		t.Error("Expected entry for 1.1.1.1 not found")
	} else {
		if result.Country != "Australia" || result.City != "Sydney" {
			t.Errorf("Incorrect data for 1.1.1.1: got %+v", result)
		}
	}

	// Test with invalid CSV format
	invalidCSVContent := "1.1.1.1,Sydney" // Missing third column
	invalidTempFile, err := os.CreateTemp("", "ip2country_invalid_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(invalidTempFile.Name())

	if _, err := invalidTempFile.WriteString(invalidCSVContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := invalidTempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	invalidService := &CSVService{
		filePath: invalidTempFile.Name(),
		data:     make(map[string]*Result),
	}

	// Test data loading with invalid format
	err = invalidService.loadData()
	if err == nil {
		t.Error("loadData should fail with invalid CSV format")
	}
}

func TestIsValidIP(t *testing.T) {
	testCases := []struct {
		ip       string
		expected bool
	}{
		{"1.1.1.1", true},
		{"8.8.8.8", true},
		{"192.168.1.1", true},
		{"256.256.256.256", false}, // Invalid IP address
		{"not-an-ip", false},       // Not an IP address
		{"", false},                // Empty string
		{"2001:db8::1", true},      // IPv6 address
	}

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			result := isValidIP(tc.ip)
			if result != tc.expected {
				t.Errorf("isValidIP(%q) = %v, expected %v", tc.ip, result, tc.expected)
			}
		})
	}
}

func TestLookupIP(t *testing.T) {
	// Create a service with test data
	service := &CSVService{
		data: map[string]*Result{
			"1.1.1.1": {Country: "Australia", City: "Sydney"},
			"8.8.8.8": {Country: "United States", City: "Mountain View"},
		},
	}

	// Test cases
	testCases := []struct {
		name        string
		ip          string
		wantResult  *Result
		wantErr     error
		description string
	}{
		{
			name:        "Valid IP - Found",
			ip:          "1.1.1.1",
			wantResult:  &Result{Country: "Australia", City: "Sydney"},
			wantErr:     nil,
			description: "Should return data for a valid IP that exists in the database",
		},
		{
			name:        "Valid IP - Not Found",
			ip:          "9.9.9.9",
			wantResult:  nil,
			wantErr:     ErrIPNotFound,
			description: "Should return ErrIPNotFound for a valid IP that doesn't exist in the database",
		},
		{
			name:        "Invalid IP",
			ip:          "not-an-ip",
			wantResult:  nil,
			wantErr:     ErrInvalidIP,
			description: "Should return ErrInvalidIP for an invalid IP address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.LookupIP(tc.ip)

			// Check error
			if err != tc.wantErr {
				t.Errorf("LookupIP(%q) error = %v, wantErr %v", tc.ip, err, tc.wantErr)
				return
			}

			// Check result
			if tc.wantResult == nil {
				if result != nil {
					t.Errorf("LookupIP(%q) = %v, want nil", tc.ip, result)
				}
			} else if result == nil {
				t.Errorf("LookupIP(%q) = nil, want %v", tc.ip, tc.wantResult)
			} else if result.Country != tc.wantResult.Country || result.City != tc.wantResult.City {
				t.Errorf("LookupIP(%q) = %v, want %v", tc.ip, result, tc.wantResult)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	// Create a temporary test CSV file
	testCSVContent := "1.1.1.1,Sydney,Australia\n8.8.8.8,Mountain View,United States"
	tempFile, err := os.CreateTemp("", "ip2country_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(testCSVContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	testCases := []struct {
		name        string
		dbType      string
		expectError bool
	}{
		{
			name:        "CSV Database Type",
			dbType:      "csv",
			expectError: false,
		},
		{
			name:        "MMDB Database Type",
			dbType:      "mmdb",
			expectError: true, // Currently not implemented
		},
		{
			name:        "MySQL Database Type",
			dbType:      "mysql",
			expectError: true, // Currently not implemented
		},
		{
			name:        "Unsupported Database Type",
			dbType:      "unknown",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, err := NewService(tempFile.Name(), tc.dbType)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for database type %q, got nil", tc.dbType)
				}
				if service != nil {
					t.Errorf("Expected nil service for database type %q, got non-nil", tc.dbType)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error for database type %q, got: %v", tc.dbType, err)
				}
				if service == nil {
					t.Errorf("Expected non-nil service for database type %q, got nil", tc.dbType)
				}
			}
		})
	}
}
