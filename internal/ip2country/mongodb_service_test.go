package ip2country

import (
	"testing"
)

func TestNewMongoDBService(t *testing.T) {
	// Test creating a new MongoDB service
	service, err := NewMongoDBService("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to create MongoDB service: %v", err)
	}
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	if service.uri != "mongodb://localhost:27017" {
		t.Errorf("Expected uri to be 'mongodb://localhost:27017', got '%s'", service.uri)
	}
}

func TestMongoDBServiceLookupIP(t *testing.T) {
	// Since the MongoDB implementation is a stub, we just test the interface compliance
	service, _ := NewMongoDBService("mongodb://localhost:27017")

	// Test the LookupIP function
	result, err := service.LookupIP("192.168.1.1")

	// Both should be nil as per the stub implementation
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}
