package requestparser

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestParsePathInt64_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "123")

	got, err := ParsePathInt64(req, "key")
	if err != nil {
		t.Fatalf("ParsePathInt64() returned unexpected error: %v", err)
	}

	want := int64(123)
	if got != want {
		t.Errorf("ParsePathInt64() = %v, want %v", got, want)
	}
}

func TestParsePathInt64_InvalidValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "abc")

	_, err := ParsePathInt64(req, "key")
	if err == nil {
		t.Fatal("ParsePathInt64() returned nil error for invalid value")
	}

	if !errors.Is(err, ErrInvalidValue) {
		t.Errorf("error = %v, want ErrInvalidValue", err)
	}

	if errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("error must not wrap strconv.ErrSyntax: %v", err)
	}

	if !strings.Contains(err.Error(), "key") {
		t.Errorf("error %q does not identify parameter %q", err, "key")
	}
}

func TestParsePathInt64_MissingValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)

	_, err := ParsePathInt64(req, "missing_key")
	if err == nil {
		t.Fatal("ParsePathInt64() returned nil error for missing parameter")
	}

	if !errors.Is(err, ErrMissingParameter) {
		t.Errorf("error = %v, want ErrMissingParameter", err)
	}

	if !strings.Contains(err.Error(), "missing_key") {
		t.Errorf("error %q does not identify parameter %q", err, "missing_key")
	}
}

func TestParsePathFloat64_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "3.14")

	got, err := ParsePathFloat64(req, "key")
	if err != nil {
		t.Fatalf("ParsePathFloat64() returned unexpected error: %v", err)
	}

	want := 3.14
	if got != want {
		t.Errorf("ParsePathFloat64() = %v, want %v", got, want)
	}
}

func TestParsePathFloat64_InvalidValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "abc")

	_, err := ParsePathFloat64(req, "key")
	if err == nil {
		t.Fatal("ParsePathFloat64() returned nil error for invalid value")
	}

	if !errors.Is(err, ErrInvalidValue) {
		t.Errorf("error = %v, want ErrInvalidValue", err)
	}

	if errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("error must not wrap strconv.ErrSyntax: %v", err)
	}

	if !strings.Contains(err.Error(), "key") {
		t.Errorf("error %q does not identify parameter %q", err, "key")
	}
}

func TestParsePathFloat64_MissingValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)

	_, err := ParsePathFloat64(req, "missing_key")
	if err == nil {
		t.Fatal("ParsePathFloat64() returned nil error for missing parameter")
	}

	if !errors.Is(err, ErrMissingParameter) {
		t.Errorf("error = %v, want ErrMissingParameter", err)
	}

	if !strings.Contains(err.Error(), "missing_key") {
		t.Errorf("error %q does not identify parameter %q", err, "missing_key")
	}
}

func TestParsePathBool_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "true")

	got, err := ParsePathBool(req, "key")
	if err != nil {
		t.Fatalf("ParsePathBool() returned unexpected error: %v", err)
	}

	want := true
	if got != want {
		t.Errorf("ParsePathBool() = %v, want %v", got, want)
	}
}

func TestParsePathBool_InvalidValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "not-a-bool")

	_, err := ParsePathBool(req, "key")
	if err == nil {
		t.Fatal("ParsePathBool() returned nil error for invalid value")
	}

	if !errors.Is(err, ErrInvalidValue) {
		t.Errorf("error = %v, want ErrInvalidValue", err)
	}

	if errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("error must not wrap strconv.ErrSyntax: %v", err)
	}

	if !strings.Contains(err.Error(), "key") {
		t.Errorf("error %q does not identify parameter %q", err, "key")
	}
}

func TestParsePathBool_MissingValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)

	_, err := ParsePathBool(req, "missing_key")
	if err == nil {
		t.Fatal("ParsePathBool() returned nil error for missing parameter")
	}

	if !errors.Is(err, ErrMissingParameter) {
		t.Errorf("error = %v, want ErrMissingParameter", err)
	}

	if !strings.Contains(err.Error(), "missing_key") {
		t.Errorf("error %q does not identify parameter %q", err, "missing_key")
	}
}

func TestParsePathString_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "example")

	got, err := ParsePathString(req, "key")
	if err != nil {
		t.Fatalf("ParsePathString() returned unexpected error: %v", err)
	}

	want := "example"
	if got != want {
		t.Errorf("ParsePathString() = %v, want %v", got, want)
	}
}

func TestParsePathString_InvalidValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req.SetPathValue("key", "")

	_, err := ParsePathString(req, "key")
	if err == nil {
		t.Fatal("ParsePathString() returned nil error for invalid value")
	}

	if !strings.Contains(err.Error(), "key") {
		t.Errorf("error %q does not identify parameter %q", err, "key")
	}
}

func TestParsePathString_MissingValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)

	_, err := ParsePathString(req, "missing_key")
	if err == nil {
		t.Fatal("ParsePathString() returned nil error for missing parameter")
	}

	if !errors.Is(err, ErrMissingParameter) {
		t.Errorf("error = %v, want ErrMissingParameter", err)
	}

	if !strings.Contains(err.Error(), "missing_key") {
		t.Errorf("error %q does not identify parameter %q", err, "missing_key")
	}
}
