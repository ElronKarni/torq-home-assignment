package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DataPath         string
	RateLimit        int
	Port             int
	IP2CountryDBType string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading it:", err)
	}

	config := &Config{}

	// Read IP2COUNTRY_DATA_PATH
	dataPath := os.Getenv("IP2COUNTRY_DATA_PATH")
	if dataPath == "" {
		dataPath = "data/ip2country.csv" // Default value
	}
	config.DataPath = dataPath

	// Read IP2COUNTRY_DB_TYPE
	dbType := os.Getenv("IP2COUNTRY_DB_TYPE")
	if dbType == "" {
		dbType = "csv" // Default value
	}
	config.IP2CountryDBType = dbType

	// Read RATE_LIMIT
	rateLimitStr := os.Getenv("RATE_LIMIT")
	log.Println("rateLimitStr", rateLimitStr)
	if rateLimitStr == "" {
		config.RateLimit = 100 // Default value
	} else {
		rateLimit, err := strconv.Atoi(rateLimitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT value: %v", err)
		}
		config.RateLimit = rateLimit
	}

	// Read PORT
	portStr := os.Getenv("PORT")
	if portStr == "" {
		config.Port = 8080 // Default value
	} else {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value: %v", err)
		}
		config.Port = port
	}

	return config, nil
}
