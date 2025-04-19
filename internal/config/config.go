package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type BackendConfig struct {
	Type      string // "csv", "mongo", "redis", etc.
	CSVPath   string
	MongoURI  string
	RedisAddr string
}

// Config holds the application-wide settings.
type Config struct {
	Port           int
	RateLimit      int
	IP2Country     BackendConfig
	AllowedOrigins []string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading it:", err)
	}

	// Read RATE_LIMIT
	rateLimit := 100
	if rateLimitStr := os.Getenv("RATE_LIMIT"); rateLimitStr != "" {
		rateLimitInt, err := strconv.Atoi(rateLimitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT value: %v", err)
		}
		rateLimit = rateLimitInt
	}

	// Read PORT
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value: %v", err)
		}
		port = portInt
	}

	// Read IP2Country DB Type
	dbType := "csv"
	if dbTypeStr := os.Getenv("IP2COUNTRY_DB_TYPE"); dbTypeStr != "" {
		dbType = dbTypeStr
	}

	// Read CSV Path
	dataPath := "data/ip2country.csv"
	if dataPathStr := os.Getenv("CSV_DATA_PATH"); dataPathStr != "" {
		dataPath = dataPathStr
	}

	// Read Mongo URI
	MongoURI := "mongodb://localhost:27017"
	if mongoURI := os.Getenv("MONGO_URI"); mongoURI != "" {
		MongoURI = mongoURI
	}

	// Read Redis Address
	RedisAddr := "localhost:6379"
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		RedisAddr = redisAddr
	}

	// Read allowed origins for CORS
	allowedOrigins := []string{"http://localhost:3000"}
	if originsStr := os.Getenv("ALLOWED_ORIGINS"); originsStr != "" {
		allowedOrigins = strings.Split(originsStr, ",")
	}

	config := &Config{
		Port:           port,
		RateLimit:      rateLimit,
		AllowedOrigins: allowedOrigins,
		IP2Country: BackendConfig{
			Type:      dbType,
			CSVPath:   dataPath,
			MongoURI:  MongoURI,
			RedisAddr: RedisAddr,
		},
	}

	return config, nil
}
