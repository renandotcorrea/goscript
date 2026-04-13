package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMustGet(t *testing.T) {
	t.Setenv("ENV_TEST_REQUIRED", "value")

	got := MustGet("ENV_TEST_REQUIRED")
	if got != "value" {
		t.Fatalf("unexpected value: got %s, want value", got)
	}
}

func TestMustGet_PanicsWhenMissing(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when variable is missing")
		}
	}()

	MustGet("ENV_TEST_REQUIRED_MISSING")
}

func TestMustGet_PanicsWhenEmpty(t *testing.T) {
	t.Setenv("ENV_TEST_REQUIRED_EMPTY", "")

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when variable is empty")
		}
	}()

	MustGet("ENV_TEST_REQUIRED_EMPTY")
}

func TestGetOr(t *testing.T) {
	t.Setenv("ENV_TEST_OPTIONAL", "set")

	got := GetOr("ENV_TEST_OPTIONAL", "fallback")
	if got != "set" {
		t.Fatalf("unexpected value: got %s, want set", got)
	}
}

func TestGetOr_ReturnsDefaultWhenMissing(t *testing.T) {
	got := GetOr("ENV_TEST_OPTIONAL_MISSING", "fallback")
	if got != "fallback" {
		t.Fatalf("unexpected value: got %s, want fallback", got)
	}
}

func TestGetOr_ReturnsDefaultWhenEmpty(t *testing.T) {
	t.Setenv("ENV_TEST_OPTIONAL_EMPTY", "")

	got := GetOr("ENV_TEST_OPTIONAL_EMPTY", "fallback")
	if got != "fallback" {
		t.Fatalf("unexpected value: got %s, want fallback", got)
	}
}

func TestLoadFile_DefaultPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "API_TOKEN=abc123\n# comment\nEMPTY_OK=\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(originalWd); chdirErr != nil {
			t.Fatalf("restore cwd failed: %v", chdirErr)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	if err := LoadFile(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("API_TOKEN"); got != "abc123" {
		t.Fatalf("unexpected API_TOKEN: got %q", got)
	}

	if got := os.Getenv("EMPTY_OK"); got != "" {
		t.Fatalf("unexpected EMPTY_OK: got %q", got)
	}
}

func TestLoadFile_MissingDefault(t *testing.T) {
	dir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(originalWd); chdirErr != nil {
			t.Fatalf("restore cwd failed: %v", chdirErr)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	if err := LoadFile(); err == nil {
		t.Fatal("expected error when .env is missing")
	}
}

func TestLoadFile_CustomPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.env")
	if err := os.WriteFile(path, []byte("REGION=us-east-1\n"), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := LoadFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("REGION"); got != "us-east-1" {
		t.Fatalf("unexpected REGION: got %q", got)
	}
}

func TestLoadFile_IgnoresMalformedLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "malformed.env")
	content := "VALID=yes\nINVALID_LINE\n=emptykey\nANOTHER=ok\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := LoadFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := os.Getenv("VALID"); got != "yes" {
		t.Fatalf("unexpected VALID: got %q", got)
	}

	if got := os.Getenv("ANOTHER"); got != "ok" {
		t.Fatalf("unexpected ANOTHER: got %q", got)
	}
}
