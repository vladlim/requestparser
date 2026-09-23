// Package requestparser parses typed path and query parameters from HTTP requests.
//
// Path functions read Request.PathValue populated by a router or SetPathValue.
// Query functions decode the full RawQuery on each call and reject malformed queries.
// Or functions use defaults only for absent keys; Slice functions use repeated keys.
// Use errors.Is to inspect exported error categories. Value validation such as
// minimums, maximums, and allowed strings belongs to the caller.
package requestparser
