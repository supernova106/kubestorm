package utils

import "os"

// Get environment variable
// Otherwise, fallback to default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
