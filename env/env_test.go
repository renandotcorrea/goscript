package env

import "testing"

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
