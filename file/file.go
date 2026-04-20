// Package file provides utilities for reading and writing JSON files.
package file

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJson opens a JSON file and unmarshals its content into the dest parameter.
//
// The filePath argument specifies the path to the JSON file to read.
// The dest argument is a pointer to a value where the unmarshaled data will be stored.
// It returns an error if the file cannot be read or if the JSON is invalid.
func ReadJson(filePath string, dest any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read json file %q: %w", filePath, err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("unmarshal json file %q: %w", filePath, err)
	}

	return nil
}

// WriteJson creates or truncates a JSON file and writes the marshaled src value.
//
// The filePath argument specifies the path where the JSON file will be created or overwritten.
// The src argument is the value to be marshaled as JSON.
// The file is created with permissions 0o644 (rw-r--r--).
// It returns an error if the value cannot be marshaled as JSON or if the file cannot be written.
func WriteJson(filePath string, src any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal json for %q: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("write json file %q: %w", filePath, err)
	}

	return nil
}
