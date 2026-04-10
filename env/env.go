package env

import (
	"fmt"
	"os"
	"time"
)

// MustGet returns the environment variable value for key.
// It panics when the variable is not set or is set to an empty string.
func MustGet(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		panic("missing required environment variable: " + key)
	}

	return value
}

// GetOr returns the environment variable value for key.
// If the variable is not set or is empty, it returns def.
func GetOr(key, def string) string {
	return getOr(key, def)
}

func getOr(key, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return def
	}

	return value
}

func GetIntOr(key string, def int) int {
	value := getOr(key, "")

	if value == "" {
		return def
	}

	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		return def
	}

	return intValue
}

func GetBoolOr(key string, def bool) bool {
	value := getOr(key, "")

	if value == "" {
		return def
	}

	var boolValue bool
	_, err := fmt.Sscanf(value, "%t", &boolValue)
	if err != nil {
		return def
	}

	return boolValue
}

// GetFloatOr returns the float64 value of the environment variable for key.
// If the variable is not set, empty, or cannot be parsed as a float, it returns def.
func GetFloatOr(key string, def float64) float64 {
	value := getOr(key, "")

	if value == "" {
		return def
	}

	var floatValue float64
	_, err := fmt.Sscanf(value, "%f", &floatValue)
	if err != nil {
		return def
	}

	return floatValue
}

func GetDurationOr(key string, def time.Duration) time.Duration {
	value := getOr(key, "")

	if value == "" {
		return def
	}

	var durationValue int64
	_, err := fmt.Sscanf(value, "%d", &durationValue)
	if err != nil {
		return def
	}

	return time.Duration(durationValue)
}
