package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type fileTestPayload struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestReadJson(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.json")
	if err := os.WriteFile(path, []byte(`{"name":"alpha","count":2}`), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	var got fileTestPayload
	if err := ReadJson(path, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Name != "alpha" || got.Count != 2 {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestReadJson_FileNotFound(t *testing.T) {
	var got fileTestPayload
	if err := ReadJson("missing.json", &got); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadJson_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(path, []byte(`{"name":`), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	var got fileTestPayload
	if err := ReadJson(path, &got); err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestWriteJson(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	input := fileTestPayload{Name: "beta", Count: 7}

	if err := WriteJson(path, input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got fileTestPayload
	if err := ReadJson(path, &got); err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}

	if got != input {
		t.Fatalf("unexpected roundtrip payload: got %+v, want %+v", got, input)
	}
}

func TestWriteJson_MarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	if err := WriteJson(path, make(chan int)); err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestWriteJson_WriteError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing", "out.json")
	if err := WriteJson(path, fileTestPayload{Name: "x"}); err == nil {
		t.Fatal("expected write error")
	}
}

func ExampleReadJson() {
	// Example demonstrates reading and unmarshaling a JSON file
	dir := os.TempDir()
	path := filepath.Join(dir, "example.json")
	payload := fileTestPayload{Name: "goscript", Count: 5}
	WriteJson(path, payload)
	defer os.Remove(path)

	var result fileTestPayload
	err := ReadJson(path, &result)
	if err != nil {
		return
	}
	fmt.Printf("%s: %d\n", result.Name, result.Count)
	// Output: goscript: 5
}

func ExampleWriteJson() {
	// Example demonstrates marshaling and writing a JSON file
	dir := os.TempDir()
	path := filepath.Join(dir, "output.json")
	defer os.Remove(path)

	payload := fileTestPayload{Name: "example", Count: 42}
	err := WriteJson(path, payload)
	if err != nil {
		return
	}

	var result fileTestPayload
	ReadJson(path, &result)
	fmt.Printf("Written: %s\n", result.Name)
	// Output: Written: example
}
