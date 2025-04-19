package ip2country

import (
	"testing"
)

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{name: "Valid IPv4", ip: "192.168.1.1", expected: true},
		{name: "Valid IPv6", ip: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", expected: true},
		{name: "Invalid IP", ip: "256.256.256.256", expected: false},
		{name: "Not an IP", ip: "not-an-ip", expected: false},
		{name: "Empty string", ip: "", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidIP(tc.ip)
			if result != tc.expected {
				t.Errorf("isValidIP(%s) = %v, expected %v", tc.ip, result, tc.expected)
			}
		})
	}
}
