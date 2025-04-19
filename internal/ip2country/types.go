package ip2country

import "errors"

// Result represents the result of an IP lookup
type Result struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

// Custom errors
var (
	ErrInvalidIP  = errors.New("invalid IP address")
	ErrIPNotFound = errors.New("IP address not found")
)
