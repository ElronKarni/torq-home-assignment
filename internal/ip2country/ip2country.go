package ip2country

import (
	"fmt"
	"ip2country-api/internal/config"
)

// Service defines the interface for the IP-to-country service
type Service interface {
	LookupIP(ip string) (*Result, error)
}

// NewService creates a new IP-to-country lookup service based on the configuration.
func NewService(config config.BackendConfig) (Service, error) {

	switch config.Type {
	case "csv":
		return NewCSVService(config.CSVPath)
	case "mongodb":
		return NewMongoDBService(config.MongoURI)
	case "redis":
		return NewRedisService(config.RedisAddr)

	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}
