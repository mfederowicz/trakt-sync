// Package internal used for client and services
package internal

import (
	"testing"
	"time"
)

func TestDateFormat(t *testing.T) {
	in := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-01T00:00Z")

	// test date format
	if got, want := in, "2025-01-01T00:00Z"; got != want {
		t.Errorf("in is %v, want %v", in, want)
	}
}
