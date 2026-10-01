// Package buffer used to write buffer bytes and printout logs if exists
package buffer

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, "hello ")
	Write(&buf, "world")
	if got := buf.String(); got != "hello world" {
		t.Errorf("buffer = %q, want %q", got, "hello world")
	}
}
