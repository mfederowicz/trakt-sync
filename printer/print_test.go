// Package printer is replacement for fmt.* functions
package printer

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// failingWriter fails every write, to reach the error branches.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

// captureLog returns what fn logs through the log package.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	orig := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(orig)
	fn()
	return buf.String()
}

// captureStdout returns what fn writes to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(out)
}

func TestStdoutPrinters(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
		want string
	}{
		{name: "Println", fn: func() { Println("a", 1) }, want: "a 1\n"},
		{name: "Printf", fn: func() { Printf("%s=%d\n", "n", 2) }, want: "n=2\n"},
		{name: "Print", fn: func() { Print("x", "y") }, want: "xy"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := captureStdout(t, tt.fn); got != tt.want {
				t.Errorf("stdout = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriterPrinters(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(w io.Writer)
		want    string
		wantLog string
	}{
		{name: "Fprint", fn: func(w io.Writer) { Fprint(w, "a", "b") }, want: "ab", wantLog: "Fprint error: disk full"},
		{name: "Fprintf", fn: func(w io.Writer) { Fprintf(w, "%d-%s", 3, "c") }, want: "3-c", wantLog: "Fprintf error: disk full"},
		{name: "Fprintln", fn: func(w io.Writer) { Fprintln(w, "d", 4) }, want: "d 4\n", wantLog: "Fprintln error: disk full"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.fn(&buf)
			if got := buf.String(); got != tt.want {
				t.Errorf("output = %q, want %q", got, tt.want)
			}
			logged := captureLog(t, func() { tt.fn(failingWriter{}) })
			if !strings.Contains(logged, tt.wantLog) {
				t.Errorf("log = %q, want it to contain %q", logged, tt.wantLog)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	inner := errors.New("inner")
	err := Errorf("wrap %d: %w", 5, inner)
	if err.Error() != "wrap 5: inner" || !errors.Is(err, inner) {
		t.Errorf("Errorf = %v, want a wrapped error", err)
	}
}
