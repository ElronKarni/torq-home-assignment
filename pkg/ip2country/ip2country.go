package ip2country

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
)

// Custom errors
var (
	ErrInvalidIP  = errors.New("invalid IP address")
	ErrIPNotFound = errors.New("IP address not found")
)

// Result represents the result of an IP lookup
type Result struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

// Service defines the interface for the IP-to-country service
type Service interface {
	LookupIP(ip string) (*Result, error)
}

// NewService creates a new IP-to-country lookup service based on the configuration.
func NewService(filePath string, dbType string) (Service, error) {
	switch dbType {
	case "csv":
		return NewCSVService(filePath)
	case "mmdb":
		// This would be implemented in the future
		return nil, fmt.Errorf("mmdb database type not implemented yet")
	case "mysql":
		// This would be implemented in the future
		return nil, fmt.Errorf("mysql database type not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// CSVService implements Service by reading data from a CSV file
type CSVService struct {
	filePath string
	data     map[string]*Result
	mu       sync.RWMutex
}

// NewCSVService creates a new CSVService with the given CSV file path
func NewCSVService(filePath string) (*CSVService, error) {
	service := &CSVService{
		filePath: filePath,
		data:     make(map[string]*Result),
	}

	if err := service.loadData(); err != nil {
		return nil, err
	}

	return service, nil
}

// loadData reads the CSV file and loads the data into memory
func (s *CSVService) loadData() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, record := range records {
		if len(record) != 3 {
			return fmt.Errorf("invalid CSV format, expected 3 columns (ip,city,country)")
		}
		ip := record[0]
		city := record[1]
		country := record[2]

		s.data[ip] = &Result{
			Country: country,
			City:    city,
		}
	}

	return nil
}

// isValidIP checks if the provided string is a valid IP address
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// LookupIP returns country information for a given IP address
func (s *CSVService) LookupIP(ip string) (*Result, error) {
	// Validate IP address format
	if !isValidIP(ip) {
		return nil, ErrInvalidIP
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result, found := s.data[ip]
	if !found {
		return nil, ErrIPNotFound
	}

	return result, nil
}
