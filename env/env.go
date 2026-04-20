// Package env provides utilities for reading and loading environment variables.
package env

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// MustGet returns the value of the environment variable named by the key.
// It panics if the variable is not set or is set to an empty string.
func MustGet(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		panic("missing required environment variable: " + key)
	}

	return value
}

// GetOr returns the value of the environment variable named by the key.
// If the variable is not set or is empty, it returns the default value def.
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

// GetIntOr returns the int value of the environment variable for key.
// If the variable is not set, empty, or cannot be parsed as an int, it returns def.
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

// GetBoolOr returns the bool value of the environment variable for key.
// If the variable is not set, empty, or cannot be parsed as a bool, it returns def.
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

// GetDurationOr returns the duration value of the environment variable for key.
// If the variable is not set, empty, or cannot be parsed as a duration, it returns def.
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

// LoadFile reads environment variables from file and sets them in the current process.
// If no file path is provided, it uses ".env" by default.
// Lines that are empty, comments, or malformed (without '=') are ignored.
func LoadFile(filePath ...string) error {
	path := ".env"
	if len(filePath) > 0 && strings.TrimSpace(filePath[0]) != "" {
		path = filePath[0]
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read env file %q: %w", path, err)
	}

	lines := strings.Split(string(content), "\n")
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexRune(line, '=')
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(line[idx+1:])
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %q from %q: %w", key, path, err)
		}
	}

	return nil
}
