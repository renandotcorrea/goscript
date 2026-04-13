package file

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJson opens a JSON file and unmarshals its content into dest.
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
