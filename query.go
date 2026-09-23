package requestparser

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ParseQueryInt64 reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Parses a signed base-10 integer; invalid syntax and overflow return ErrInvalidValue.
//
// On error, the value result is -1. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryInt64(r *http.Request, key string) (int64, error) {
	if err := requestCheck(r); err != nil {
		return -1, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return -1, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	if len(raw) != 1 {
		return -1, fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	val, err := strconv.ParseInt(raw[0], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParseQueryInt64Or reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// Returns defaultValue only when the key is absent. Invalid input never uses the default.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Parses a signed base-10 integer; invalid syntax and overflow return ErrInvalidValue.
//
// On error, the value result is -1. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryInt64Or(
	r *http.Request,
	key string,
	defaultValue int64,
) (int64, error) {
	if err := requestCheck(r); err != nil {
		return -1, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return defaultValue, nil
	}

	if len(raw) != 1 {
		return -1, fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	val, err := strconv.ParseInt(raw[0], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParseQueryInt64Slice reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys supply elements in order, preserving duplicates. Commas are not separators.
// Returns nil on error; conversion errors include the zero-based element index.
//
// Parses a signed base-10 integer; invalid syntax and overflow return ErrInvalidValue.
//
// On error, the value result is nil. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryInt64Slice(r *http.Request, key string) ([]int64, error) {
	if err := requestCheck(r); err != nil {
		return nil, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return nil, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	res := make([]int64, len(raw))

	for i, v := range raw {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: key = %s: index = %d", ErrInvalidValue, err, key, i)
		}

		res[i] = val
	}

	return res, nil
}

// ParseQueryFloat64 reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Uses strconv.ParseFloat with 64-bit precision, including NaN and infinities.
// Invalid syntax and overflow return ErrInvalidValue; finite values are not required.
//
// On error, the value result is -1. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryFloat64(r *http.Request, key string) (float64, error) {
	if err := requestCheck(r); err != nil {
		return -1, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return -1, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	if len(raw) != 1 {
		return -1, fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	val, err := strconv.ParseFloat(raw[0], 64)
	if err != nil {
		return -1, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParseQueryFloat64Or reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// Returns defaultValue only when the key is absent. Invalid input never uses the default.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Uses strconv.ParseFloat with 64-bit precision, including NaN and infinities.
// Invalid syntax and overflow return ErrInvalidValue; finite values are not required.
//
// On error, the value result is -1. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryFloat64Or(
	r *http.Request,
	key string,
	defaultValue float64,
) (float64, error) {
	if err := requestCheck(r); err != nil {
		return -1, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return defaultValue, nil
	}

	if len(raw) != 1 {
		return -1, fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	val, err := strconv.ParseFloat(raw[0], 64)
	if err != nil {
		return -1, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParseQueryFloat64Slice reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys supply elements in order, preserving duplicates. Commas are not separators.
// Returns nil on error; conversion errors include the zero-based element index.
//
// Uses strconv.ParseFloat with 64-bit precision, including NaN and infinities.
// Invalid syntax and overflow return ErrInvalidValue; finite values are not required.
//
// On error, the value result is nil. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryFloat64Slice(r *http.Request, key string) ([]float64, error) {
	if err := requestCheck(r); err != nil {
		return nil, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return nil, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	res := make([]float64, len(raw))

	for i, v := range raw {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: key = %s: index = %d", ErrInvalidValue, err, key, i)
		}

		res[i] = val
	}

	return res, nil
}

// ParseQueryBool reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, and False.
// Other values return ErrInvalidValue.
//
// On error, the value result is false. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryBool(r *http.Request, key string) (bool, error) {
	if err := requestCheck(r); err != nil {
		return false, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return false, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	if len(raw) != 1 {
		return false, fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	val, err := strconv.ParseBool(raw[0])
	if err != nil {
		return false, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParseQueryBoolOr reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// Returns defaultValue only when the key is absent. Invalid input never uses the default.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, and False.
// Other values return ErrInvalidValue.
//
// On error, the value result is false. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryBoolOr(
	r *http.Request,
	key string,
	defaultValue bool,
) (bool, error) {
	if err := requestCheck(r); err != nil {
		return false, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return defaultValue, nil
	}

	if len(raw) != 1 {
		return false, fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	val, err := strconv.ParseBool(raw[0])
	if err != nil {
		return false, fmt.Errorf("%w: %v: key = %s", ErrInvalidValue, err, key)
	}

	return val, nil
}

// ParseQueryBoolSlice reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys supply elements in order, preserving duplicates. Commas are not separators.
// Returns nil on error; conversion errors include the zero-based element index.
//
// Accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, and False.
// Other values return ErrInvalidValue.
//
// On error, the value result is nil. Check the error before using the result.
// Conversion errors wrap ErrInvalidValue; the strconv error is included only as text.
func ParseQueryBoolSlice(r *http.Request, key string) ([]bool, error) {
	if err := requestCheck(r); err != nil {
		return nil, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return nil, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	res := make([]bool, len(raw))

	for i, v := range raw {
		val, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: key = %s: index = %d", ErrInvalidValue, err, key, i)
		}

		res[i] = val
	}

	return res, nil
}

// ParseQueryString reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Preserves whitespace and accepts an explicitly empty value.
//
// On error, the value result is an empty string. Check the error before using the result.
func ParseQueryString(r *http.Request, key string) (string, error) {
	if err := requestCheck(r); err != nil {
		return "", err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return "", fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	if len(raw) != 1 {
		return "", fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	return raw[0], nil
}

// ParseQueryStringOr reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// Returns defaultValue only when the key is absent. Invalid input never uses the default.
// Repeated keys return ErrMultipleValues, even if their values are equal.
//
// Preserves whitespace and accepts an explicitly empty value.
//
// On error, the value result is an empty string. Check the error before using the result.
func ParseQueryStringOr(
	r *http.Request,
	key string,
	defaultValue string,
) (string, error) {
	if err := requestCheck(r); err != nil {
		return "", err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return defaultValue, nil
	}

	if len(raw) != 1 {
		return "", fmt.Errorf("%w: key = %s", ErrMultipleValues, key)
	}

	return raw[0], nil
}

// ParseQueryStringSlice reads the named query parameter from r.URL.RawQuery.
//
// A nil request returns ErrNilRequest; a nil URL returns ErrNilURL.
// Malformed encoding anywhere in the query returns ErrInvalidQuery.
//
// An absent key returns ErrMissingParameter.
// Repeated keys supply elements in order, preserving duplicates. Commas are not separators.
// Returns nil on error; conversion errors include the zero-based element index.
//
// Preserves whitespace and accepts an explicitly empty value.
//
// On error, the value result is nil. Check the error before using the result.
func ParseQueryStringSlice(r *http.Request, key string) ([]string, error) {
	if err := requestCheck(r); err != nil {
		return nil, err
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidQuery, err)
	}

	raw, ok := queryValues[key]
	if !ok {
		return nil, fmt.Errorf("%w: key = %s", ErrMissingParameter, key)
	}

	res := make([]string, len(raw))
	copy(res, raw)

	return res, nil
}

// requestCheck verifies that query parsing can safely access the request URL.
func requestCheck(r *http.Request) error {
	if r == nil {
		return ErrNilRequest
	}

	if r.URL == nil {
		return ErrNilURL
	}

	return nil
}
