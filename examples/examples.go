// Run these examples with: go run ./examples
package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/vladlim/requestparser"
)

func main() {
	// ServeMux populates PathValue before invoking the handler.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/users/42", nil))
	fmt.Print(recorder.Body.String())

	// Repeated keys produce ordered slices; commas remain part of strings.
	r := httptest.NewRequest(http.MethodGet, "/?tag=go&tag=hello%2Cworld", nil)
	tags, err := requestparser.ParseQueryStringSlice(r, "tag")
	if err != nil {
		panic(err)
	}
	fmt.Printf("tags=%q\n", tags)

	// An explicitly supplied false value does not trigger the default.
	r = httptest.NewRequest(http.MethodGet, "/?active=false", nil)
	active, err := requestparser.ParseQueryBoolOr(r, "active", true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("active=%t\n", active)

	// Defaults do not hide invalid input. Inspect the package error category.
	r = httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
	_, err = requestparser.ParseQueryInt64Or(r, "limit", 20)
	fmt.Printf("invalid limit=%t\n", errors.Is(err, requestparser.ErrInvalidValue))
}
