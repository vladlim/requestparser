# requestparser

[![CI](https://github.com/vladlim/requestparser/actions/workflows/ci.yaml/badge.svg?branch=main&event=push)](https://github.com/vladlim/requestparser/actions/workflows/ci.yaml)

Typed path and query parameter parsing for Go HTTP handlers, with defaults,
repeated-key slices, and errors you can inspect with `errors.Is`.
Uses only the standard library. Requires Go 1.25.3 or later.

## Install

```sh
go get github.com/vladlim/requestparser
```

## Usage

```go
import (
    "fmt"
    "net/http"

    "github.com/vladlim/requestparser"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
    id, err := requestparser.ParsePathInt64(r, "id")
    if err != nil {
        http.Error(w, "invalid user ID", http.StatusBadRequest)
        return
    }
    limit, err := requestparser.ParseQueryInt64Or(r, "limit", 20)
    if err != nil {
        http.Error(w, "invalid limit", http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "user=%d limit=%d\n", id, limit)
}
```

Register the handler with `mux.HandleFunc("GET /users/{id}", userHandler)` on an
`http.NewServeMux()`. For `/users/42`, the result is `user=42 limit=20`.
Path functions read `r.PathValue`: a router must populate it, or tests/adapters
can call `r.SetPathValue`. They do not extract parameter names from URLs themselves.

Run the complete [examples](examples/examples.go) without starting a server:

```sh
go run ./examples
```

## API and behavior

Replace `<Type>` with `Int64`, `Float64`, `Bool`, or `String`:

| Function | Behavior |
| --- | --- |
| `ParsePath<Type>(r, key)` | Required, nonempty path value |
| `ParseQuery<Type>(r, key)` | Required single query value |
| `ParseQuery<Type>Or(r, key, defaultValue)` | Default only when the key is absent |
| `ParseQuery<Type>Slice(r, key)` | Required list from repeated query keys |

- `?id=1&id=2` produces `[1, 2]` with `ParseQueryInt64Slice`. Order and duplicates are preserved; comma-separated lists are not supported.
- Single-value query functions reject repeated keys, even equal values.
- `?name=` and `?name` are present empty strings, not missing parameters. String query functions accept them; numeric and boolean functions reject them. Empty path values are treated as missing.
- Defaults never replace explicit zero, false, empty strings, or invalid input.
- Integers use base 10. Booleans accept the spellings supported by `strconv.ParseBool`, including `1` and `0`. Floats follow `strconv.ParseFloat`, including `NaN` and infinities.
- Query values are URL-decoded once: `+` becomes a space, `%2B` becomes `+`. Whitespace is not trimmed.
- Each query call parses the entire raw query. Malformed encoding anywhere returns an error, even for an unrelated key.
- Slice conversion stops at the first invalid element and returns no partial result. The error includes its zero-based index.
- Business rules such as positive IDs, finite floats, and maximum limits belong to the caller.

Always check the error before using the result. Error results are `-1` for numbers,
`false` for booleans, `""` for strings, and `nil` for slices.

## Errors

Use `errors.Is(err, requestparser.ErrInvalidValue)`, not equality or message matching.

| Error | Meaning |
| --- | --- |
| `ErrNilRequest` | Request is nil |
| `ErrNilURL` | Query request has a nil URL |
| `ErrMissingParameter` | Required key is absent, or a path value is empty |
| `ErrMultipleValues` | Repeated key passed to a single-value query function |
| `ErrInvalidValue` | Value cannot be converted to the requested type |
| `ErrInvalidQuery` | Raw query cannot be decoded correctly |

Parameter errors include the key. Underlying conversion errors are included as
text only; `errors.Is`/`errors.As` do not expose the original `strconv` error.

## Development and CI

```sh
make check  # formatting check, go vet, and tests
make fmt    # format Go files locally
```

CI runs for pull requests targeting `main` and pushes to `main`.
The badge above shows the latest push workflow status for `main`: green when
passing and red when failing. It needs a workflow run on GitHub before it can
display a result. Click it to see the workflow logs.
