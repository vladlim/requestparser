package requestparser

import "errors"

var (
	// ErrNilRequest indicates a nil request.
	ErrNilRequest = errors.New("request is nil")
	// ErrInvalidValue identifies a conversion failure. The underlying conversion error is text only.
	ErrInvalidValue = errors.New("invalid parameter value")
	// ErrMissingParameter indicates an absent parameter, or an empty path parameter.
	ErrMissingParameter = errors.New("missing parameter")
	// ErrMultipleValues indicates repeated values where a single query value is required.
	ErrMultipleValues = errors.New("multiple values for parameter")
	// ErrNilURL indicates that the request has no URL.
	ErrNilURL = errors.New("request URL is nil")
	// ErrInvalidQuery indicates malformed query encoding; no partial result is returned.
	ErrInvalidQuery = errors.New("invalid query")
)
