package try

import (
	"errors"
	"testing"
)

var errTest = errors.New("test error")

func TestTry_NoError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic, got: %v", r)
		}
	}()
	Try(nil)
}

func TestTry_WithError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got none")
		}
		if r != errTest {
			t.Fatalf("unexpected panic value: got %v, want %v", r, errTest)
		}
	}()
	Try(errTest)
}

func TestTry1_NoError(t *testing.T) {
	got := Try1(42, nil)
	if got != 42 {
		t.Fatalf("unexpected value: got %d, want 42", got)
	}
}

func TestTry1_WithError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	Try1(0, errTest)
}

func TestTry2_NoError(t *testing.T) {
	v1, v2 := Try2(1, "hello", nil)
	if v1 != 1 {
		t.Fatalf("unexpected v1: got %d, want 1", v1)
	}
	if v2 != "hello" {
		t.Fatalf("unexpected v2: got %s, want hello", v2)
	}
}

func TestTry2_WithError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	Try2(0, "", errTest)
}

func TestTry3_NoError(t *testing.T) {
	v1, v2, v3 := Try3(1, "hello", true, nil)
	if v1 != 1 || v2 != "hello" || !v3 {
		t.Fatalf("unexpected values: got %d, %s, %v", v1, v2, v3)
	}
}

func TestTry3_WithError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	Try3(0, "", false, errTest)
}

func TestTry4_NoError(t *testing.T) {
	v1, v2, v3, v4 := Try4(1, "hello", true, 3.14, nil)
	if v1 != 1 || v2 != "hello" || !v3 || v4 != 3.14 {
		t.Fatalf("unexpected values: got %d, %s, %v, %f", v1, v2, v3, v4)
	}
}

func TestTry4_WithError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	Try4(0, "", false, 0.0, errTest)
}

func TestHandle_WithError(t *testing.T) {
	var handled error

	func() {
		defer Handle(func(err error) {
			handled = err
		})
		Try(errTest)
	}()

	if handled != errTest {
		t.Fatalf("unexpected handled error: got %v, want %v", handled, errTest)
	}
}

func TestHandle_NoError(t *testing.T) {
	called := false

	func() {
		defer Handle(func(err error) {
			called = true
		})
		Try(nil)
	}()

	if called {
		t.Fatal("handler should not be called when there is no panic")
	}
}

func TestHandle_NonErrorPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic to propagate")
		}
		if r != "not an error" {
			t.Fatalf("unexpected panic value: got %v", r)
		}
	}()

	func() {
		defer Handle(func(err error) {
			t.Fatal("handler should not be called for non-error panics")
		})
		panic("not an error")
	}()
}
