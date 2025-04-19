package ip2country

import (
	"os"
	"testing"
)

func TestNewCSVService(t *testing.T) {
	// Create a test CSV file
	testFile := "test_data.csv"
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.WriteString("192.168.1.1,New York,USA\n")
	f.WriteString("10.0.0.1,London,UK\n")
	f.Close()
	defer os.Remove(testFile)

	// Test successful creation
	service, err := NewCSVService(testFile)
	if err != nil {
		t.Fatalf("Failed to create CSV service: %v", err)
	}
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Test with non-existent file
	_, err = NewCSVService("non_existent_file.csv")
	if err == nil {
		t.Fatal("Expected error with non-existent file, got nil")
	}
}

func TestCSVServiceLookupIP(t *testing.T) {
	// Create a test CSV file
	testFile := "test_lookup_data.csv"
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.WriteString("192.168.1.1,New York,USA\n")
	f.WriteString("10.0.0.1,London,UK\n")
	f.Close()
	defer os.Remove(testFile)

	service, err := NewCSVService(testFile)
	if err != nil {
		t.Fatalf("Failed to create CSV service: %v", err)
	}

	tests := []struct {
		name        string
		ip          string
		wantCity    string
		wantCountry string
		wantErr     error
	}{
		{name: "Existing IP 1", ip: "192.168.1.1", wantCity: "New York", wantCountry: "USA", wantErr: nil},
		{name: "Existing IP 2", ip: "10.0.0.1", wantCity: "London", wantCountry: "UK", wantErr: nil},
		{name: "Non-existent IP", ip: "8.8.8.8", wantCity: "", wantCountry: "", wantErr: ErrIPNotFound},
		{name: "Invalid IP", ip: "invalid-ip", wantCity: "", wantCountry: "", wantErr: ErrInvalidIP},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.LookupIP(tc.ip)

			if err != tc.wantErr {
				t.Errorf("LookupIP(%s) error = %v, want %v", tc.ip, err, tc.wantErr)
				return
			}

			if err == nil {
				if result.City != tc.wantCity {
					t.Errorf("LookupIP(%s) city = %v, want %v", tc.ip, result.City, tc.wantCity)
				}
				if result.Country != tc.wantCountry {
					t.Errorf("LookupIP(%s) country = %v, want %v", tc.ip, result.Country, tc.wantCountry)
				}
			}
		})
	}
}

func TestCSVServiceInvalidFormat(t *testing.T) {
	// Create a test CSV file with invalid format
	testFile := "test_invalid_format.csv"
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.WriteString("192.168.1.1,New York\n") // Missing the country column
	f.Close()
	defer os.Remove(testFile)

	_, err = NewCSVService(testFile)
	if err == nil {
		t.Fatal("Expected error with invalid CSV format, got nil")
	}
}

func TestCSVServiceReadError(t *testing.T) {
	// Create a directory with the same name as the test file to cause a read error
	testFile := "test_dir_not_file.csv"
	err := os.Mkdir(testFile, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testFile)

	// Attempting to open a directory as a file should cause a read error
	_, err = NewCSVService(testFile)
	if err == nil {
		t.Fatal("Expected error when reading a directory as a file, got nil")
	}
}
