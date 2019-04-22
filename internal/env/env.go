package env

import "os"

// GetEnv get environment variable or use default value
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
