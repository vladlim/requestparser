package requestparser

import (
	"fmt"
	"net/http"
	"strconv"
)

// ParsePathInt64 reads the named path parameter using r.PathValue.
//
// The router or caller must populate PathValue; this function does not match URLs.
//
// A nil request returns ErrNilRequest. Missing or empty values return ErrMissingParameter.
//
// Parses a signed base-10 integer; invalid syntax and overflow return ErrInvalidValue.
//
// On error, the value result is -1. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParsePathInt64(r *http.Request, key string) (int64, error) {
	if r == nil {
		return -1, ErrNilRequest
	}

	raw := r.PathValue(key)
	if raw == "" {
		return -1, fmt.Errorf("path parameter %q: %w", key, ErrMissingParameter)
	}

	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParsePathFloat64 reads the named path parameter using r.PathValue.
//
// The router or caller must populate PathValue; this function does not match URLs.
//
// A nil request returns ErrNilRequest. Missing or empty values return ErrMissingParameter.
//
// Uses strconv.ParseFloat with 64-bit precision, including NaN and infinities.
// Invalid syntax and overflow return ErrInvalidValue; finite values are not required.
//
// On error, the value result is -1. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParsePathFloat64(r *http.Request, key string) (float64, error) {
	if r == nil {
		return -1, ErrNilRequest
	}

	raw := r.PathValue(key)
	if raw == "" {
		return -1, fmt.Errorf("path parameter %q: %w", key, ErrMissingParameter)
	}

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return -1, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParsePathBool reads the named path parameter using r.PathValue.
//
// The router or caller must populate PathValue; this function does not match URLs.
//
// A nil request returns ErrNilRequest. Missing or empty values return ErrMissingParameter.
//
// Accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, and False.
// Other values return ErrInvalidValue.
//
// On error, the value result is false. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParsePathBool(r *http.Request, key string) (bool, error) {
	if r == nil {
		return false, ErrNilRequest
	}

	raw := r.PathValue(key)
	if raw == "" {
		return false, fmt.Errorf("path parameter %q: %w", key, ErrMissingParameter)
	}

	val, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParsePathString reads the named path parameter using r.PathValue.
//
// The router or caller must populate PathValue; this function does not match URLs.
//
// A nil request returns ErrNilRequest. Missing or empty values return ErrMissingParameter.
//
// Returns the stored string unchanged, including whitespace.
//
// On error, the value result is an empty string. Check the error before using the result.
func ParsePathString(r *http.Request, key string) (string, error) {
	if r == nil {
		return "", ErrNilRequest
	}

	str := r.PathValue(key)
	if str == "" {
		return "", fmt.Errorf("path parameter %q: %w", key, ErrMissingParameter)
	}

	return str, nil
}
