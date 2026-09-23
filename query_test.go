package requestparser

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestParseQueryInt64(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    int64
		wantErr error
	}{
		{name: "success", query: "id=123", want: 123},
		{name: "zero", query: "id=0", want: 0},
		{name: "invalid value", query: "id=abc", wantErr: ErrInvalidValue},
		{name: "overflow", query: "id=9223372036854775808", wantErr: ErrInvalidValue},
		{name: "empty value", query: "id=", wantErr: ErrInvalidValue},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "multiple values", query: "id=1&id=2", wantErr: ErrMultipleValues},
		{name: "malformed repeated value", query: "id=1&id=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "id=1&other=%ZZ", wantErr: ErrInvalidQuery},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/example", nil)
			req.URL.RawQuery = tt.query
			got, err := ParseQueryInt64(req, "id")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ParseQueryInt64() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if got != tt.want {
					t.Errorf("ParseQueryInt64() = %d, want %d", got, tt.want)
				}
				return
			}
			if tt.wantErr != ErrInvalidQuery && !strings.Contains(err.Error(), "key = id") {
				t.Errorf("error does not identify parameter: %v", err)
			}
			if tt.wantErr == ErrInvalidValue {
				if !strings.Contains(err.Error(), "strconv.ParseInt") {
					t.Errorf("error does not retain conversion details: %v", err)
				}
				var numErr *strconv.NumError
				if errors.As(err, &numErr) {
					t.Errorf("error must not wrap strconv.NumError: %v", err)
				}
			}
		})
	}
}

func TestParseQueryInt64_InvalidRequest(t *testing.T) {
	for _, tt := range []struct {
		name    string
		req     *http.Request
		wantErr error
	}{
		{name: "nil request", wantErr: ErrNilRequest},
		{name: "nil URL", req: &http.Request{}, wantErr: ErrNilURL},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseQueryInt64(tt.req, "id"); !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseQueryInt64() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// queryCase describes observable parser behavior for one raw query.
type queryCase[T any] struct {
	name    string
	query   string
	want    T
	wantErr error
}

func checkQueryCases[T any](t *testing.T, parse func(*http.Request, string) (T, error), cases []queryCase[T]) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/example", nil)
			req.URL.RawQuery = tt.query
			got, err := parse(req, "value")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("parse(%q) error = %v, want %v", tt.query, err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parse(%q) = %#v, want %#v", tt.query, got, tt.want)
				}
				return
			}
			if tt.wantErr != ErrInvalidQuery && !strings.Contains(err.Error(), "key = value") {
				t.Errorf("error does not identify parameter: %v", err)
			}
			if tt.wantErr == ErrInvalidValue {
				var numErr *strconv.NumError
				if errors.As(err, &numErr) {
					t.Errorf("error must not wrap strconv.NumError: %v", err)
				}
			}
		})
	}
	for _, tt := range []struct {
		name    string
		req     *http.Request
		wantErr error
	}{
		{name: "nil request", wantErr: ErrNilRequest},
		{name: "nil URL", req: &http.Request{}, wantErr: ErrNilURL},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parse(tt.req, "value"); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseQueryInt64Or(t *testing.T) {
	checkQueryCases(t, func(r *http.Request, key string) (int64, error) { return ParseQueryInt64Or(r, key, 99) }, []queryCase[int64]{
		{name: "success", query: "value=42", want: 42},
		{name: "explicit zero", query: "value=0", want: 0},
		{name: "missing", query: "other=1", want: 99},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "invalid", query: "value=abc", wantErr: ErrInvalidValue},
		{name: "empty", query: "value=", wantErr: ErrInvalidValue},
		{name: "bare key", query: "value", wantErr: ErrInvalidValue},
		{name: "overflow", query: "value=9223372036854775808", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=42&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=42&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryInt64Slice(t *testing.T) {
	checkQueryCases(t, ParseQueryInt64Slice, []queryCase[[]int64]{
		{name: "ordered repeated values", query: "value=42&value=0&value=42", want: []int64{42, 0, 42}},
		{name: "single value", query: "value=42", want: []int64{42}},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "invalid later element", query: "value=42&value=abc", wantErr: ErrInvalidValue},
		{name: "empty later element", query: "value=42&value=", wantErr: ErrInvalidValue},
		{name: "comma is not separator", query: "value=1,2", wantErr: ErrInvalidValue},
		{name: "overflow", query: "value=9223372036854775808", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=42&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=42&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryFloat64(t *testing.T) {
	checkQueryCases(t, ParseQueryFloat64, []queryCase[float64]{
		{name: "success", query: "value=3.25", want: 3.25},
		{name: "explicit zero", query: "value=0", want: 0},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "invalid", query: "value=abc", wantErr: ErrInvalidValue},
		{name: "empty", query: "value=", wantErr: ErrInvalidValue},
		{name: "bare key", query: "value", wantErr: ErrInvalidValue},
		{name: "overflow", query: "value=1e9999", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=3.25&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=3.25&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryFloat64Or(t *testing.T) {
	checkQueryCases(t, func(r *http.Request, key string) (float64, error) { return ParseQueryFloat64Or(r, key, 9.5) }, []queryCase[float64]{
		{name: "success", query: "value=3.25", want: 3.25},
		{name: "explicit zero", query: "value=0", want: 0},
		{name: "missing", query: "other=1", want: 9.5},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "invalid", query: "value=abc", wantErr: ErrInvalidValue},
		{name: "empty", query: "value=", wantErr: ErrInvalidValue},
		{name: "bare key", query: "value", wantErr: ErrInvalidValue},
		{name: "overflow", query: "value=1e9999", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=3.25&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=3.25&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryFloat64Slice(t *testing.T) {
	checkQueryCases(t, ParseQueryFloat64Slice, []queryCase[[]float64]{
		{name: "ordered repeated values", query: "value=3.25&value=0&value=3.25", want: []float64{3.25, 0, 3.25}},
		{name: "single value", query: "value=3.25", want: []float64{3.25}},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "invalid later element", query: "value=3.25&value=abc", wantErr: ErrInvalidValue},
		{name: "empty later element", query: "value=3.25&value=", wantErr: ErrInvalidValue},
		{name: "comma is not separator", query: "value=1,2", wantErr: ErrInvalidValue},
		{name: "overflow", query: "value=1e9999", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=3.25&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=3.25&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryBool(t *testing.T) {
	checkQueryCases(t, ParseQueryBool, []queryCase[bool]{
		{name: "success", query: "value=true", want: true},
		{name: "explicit zero", query: "value=false", want: false},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "invalid", query: "value=maybe", wantErr: ErrInvalidValue},
		{name: "empty", query: "value=", wantErr: ErrInvalidValue},
		{name: "bare key", query: "value", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=true&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=true&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryBoolOr(t *testing.T) {
	checkQueryCases(t, func(r *http.Request, key string) (bool, error) { return ParseQueryBoolOr(r, key, true) }, []queryCase[bool]{
		{name: "success", query: "value=true", want: true},
		{name: "explicit zero", query: "value=false", want: false},
		{name: "missing", query: "other=1", want: true},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "invalid", query: "value=maybe", wantErr: ErrInvalidValue},
		{name: "empty", query: "value=", wantErr: ErrInvalidValue},
		{name: "bare key", query: "value", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=true&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=true&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryBoolSlice(t *testing.T) {
	checkQueryCases(t, ParseQueryBoolSlice, []queryCase[[]bool]{
		{name: "ordered repeated values", query: "value=true&value=false&value=true", want: []bool{true, false, true}},
		{name: "single value", query: "value=true", want: []bool{true}},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "invalid later element", query: "value=true&value=maybe", wantErr: ErrInvalidValue},
		{name: "empty later element", query: "value=true&value=", wantErr: ErrInvalidValue},
		{name: "comma is not separator", query: "value=1,2", wantErr: ErrInvalidValue},
		{name: "malformed repeated value", query: "value=true&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=true&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryString(t *testing.T) {
	checkQueryCases(t, ParseQueryString, []queryCase[string]{
		{name: "success", query: "value=hello", want: "hello"},
		{name: "explicit zero", query: "value=", want: ""},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "bare key", query: "value", want: ""},
		{name: "decoded", query: "value=a+b%2Bc%26d", want: "a b+c&d"},
		{name: "malformed repeated value", query: "value=hello&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=hello&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryStringOr(t *testing.T) {
	checkQueryCases(t, func(r *http.Request, key string) (string, error) { return ParseQueryStringOr(r, key, "fallback") }, []queryCase[string]{
		{name: "success", query: "value=hello", want: "hello"},
		{name: "explicit zero", query: "value=", want: ""},
		{name: "missing", query: "other=1", want: "fallback"},
		{name: "duplicate", query: "value=1&value=2", wantErr: ErrMultipleValues},
		{name: "bare key", query: "value", want: ""},
		{name: "decoded", query: "value=a+b%2Bc%26d", want: "a b+c&d"},
		{name: "malformed repeated value", query: "value=hello&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=hello&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQueryStringSlice(t *testing.T) {
	checkQueryCases(t, ParseQueryStringSlice, []queryCase[[]string]{
		{name: "ordered repeated values", query: "value=hello&value=&value=hello", want: []string{"hello", "", "hello"}},
		{name: "single value", query: "value=hello", want: []string{"hello"}},
		{name: "missing", query: "other=1", wantErr: ErrMissingParameter},
		{name: "empty string", query: "value=", want: []string{""}},
		{name: "comma preserved", query: "value=a%2Cb&value=c", want: []string{"a,b", "c"}},
		{name: "decode once", query: "value=%252C&value=a+b&value=%2B", want: []string{"%2C", "a b", "+"}},
		{name: "malformed repeated value", query: "value=hello&value=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed unrelated value", query: "value=hello&other=%ZZ", wantErr: ErrInvalidQuery},
		{name: "malformed query with missing key", query: "other=%ZZ", wantErr: ErrInvalidQuery},
	})
}

func TestParseQuerySlices_InvalidElement(t *testing.T) {
	parsers := []struct {
		name  string
		parse func(*testing.T, *http.Request, string) error
	}{
		{"int64", func(t *testing.T, r *http.Request, key string) error {
			got, err := ParseQueryInt64Slice(r, key)
			if got != nil {
				t.Errorf("int64: partial result = %v, want nil", got)
			}
			return err
		}},
		{"float64", func(t *testing.T, r *http.Request, key string) error {
			got, err := ParseQueryFloat64Slice(r, key)
			if got != nil {
				t.Errorf("float64: partial result = %v, want nil", got)
			}
			return err
		}},
		{"bool", func(t *testing.T, r *http.Request, key string) error {
			got, err := ParseQueryBoolSlice(r, key)
			if got != nil {
				t.Errorf("bool: partial result = %v, want nil", got)
			}
			return err
		}},
	}
	for _, parser := range parsers {
		t.Run(parser.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/example?value=1&value=invalid&value=0", nil)
			err := parser.parse(t, req, "value")
			if !errors.Is(err, ErrInvalidValue) {
				t.Fatalf("error = %v, want ErrInvalidValue", err)
			}
			if !strings.Contains(err.Error(), "index = 1") {
				t.Errorf("error does not identify element: %v", err)
			}
		})
	}
}
