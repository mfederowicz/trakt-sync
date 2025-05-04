// Package test used for process tests
package test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// MapsStringBoolEqual  check if two maps are equal
func MapsStringBoolEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, ok := b[key]; !ok || bValue != value {
			return false
		}
	}
	return true
}

// AssertMethod check request method
func AssertMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

// Ptr is a helper routine that allocates a new T value
// to store v and returns a pointer to it.
func Ptr[T any](v T) *T {
	return &v
}

// AssertType check if v is the targetType
func AssertType(t *testing.T, v any, targetType string) {
	t.Helper()

	to := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)

	if vo.IsNil() {
		t.Errorf("unexpected type: nil, expected:%s", targetType)
		return
	}

	if to.Kind() == reflect.Ptr {
		to = to.Elem()
	}

	if to.Name() != targetType {
		t.Errorf("unexpected type: %v", to.Name())
	}
}

// SafeFprint prints msg and handle error if exists
func SafeFprint(w io.Writer, msg string) {
	if _, err := fmt.Fprint(w, msg); err != nil {
		log.Printf("write error: %v", err)
	}
}

// AssertNilError checks if error is not nil
func AssertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// AssertNoDiff check if vars different
func AssertNoDiff(t *testing.T, want, got any) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diff mismatch (-want +got):\n%v", diff)
	}
}

// AssertWrite check if writer write something
func AssertWrite(t *testing.T, w io.Writer, data []byte) {
	t.Helper()
	_, err := w.Write(data)
	AssertNilError(t, err)
}

// AssertHeader helper check if header exists
func AssertHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

// AssertURLParseError helper to check url parse error
func AssertURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

type values map[string]string

// AssertFormValues checks form values
func AssertFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	AssertNilError(t, r.ParseForm())
	if got := r.Form; !cmp.Equal(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}
