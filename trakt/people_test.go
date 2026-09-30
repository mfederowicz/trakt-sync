// Package trakt is a Go client for the Trakt API
package trakt

import (
	"strconv"
	"testing"
	"time"
)

func TestDateFormat(t *testing.T) {
	in := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-01T00:00Z")

	// test date format
	if got, want := in, strconv.Itoa(time.Now().Year())+"-01-01T00:00Z"; got != want {
		t.Errorf("in is %v, want %v", in, want)
	}
}
