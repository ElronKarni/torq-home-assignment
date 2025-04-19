package ip2country

import (
	"testing"
)

func TestNewRedisService(t *testing.T) {
	// Test creating a new Redis service
	service, err := NewRedisService("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to create Redis service: %v", err)
	}
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	if service.addr != "localhost:6379" {
		t.Errorf("Expected addr to be 'localhost:6379', got '%s'", service.addr)
	}
}

func TestRedisServiceLookupIP(t *testing.T) {
	// Since the Redis implementation is a stub, we just test the interface compliance
	service, _ := NewRedisService("localhost:6379")

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
