package ip2country

import "net"

// isValidIP checks if the provided string is a valid IP address
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}
