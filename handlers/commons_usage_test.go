// Package handlers used to handle module actions
package handlers

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// captureStdout returns what fn writes to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r) // a read error only shortens the captured output, which the asserts catch
		done <- buf.String()
	}()
	fn()
	assert.NoError(t, w.Close())
	return <-done
}

func TestGenUsage(t *testing.T) {
	c := CommonLogic{}

	tests := []struct {
		name string
		gen  func()
		want string
	}{
		{
			name: "actions",
			gen:  func() { c.GenActionsUsage("people", []string{"lists", "movies"}) },
			want: "Usage: ./trakt-sync people -a [action]\nAvailable actions:\n  - lists\n  - movies\n",
		},
		{
			name: "no actions",
			gen:  func() { c.GenActionsUsage("people", nil) },
			want: "Usage: ./trakt-sync people -a [action]\nAvailable actions:\n",
		},
		{
			name: "types",
			gen:  func() { c.GenTypeUsage("search", []string{"movie", "show"}) },
			want: "Usage: ./trakt-sync search -t [type]\nAvailable types:\n  - movie\n  - show\n",
		},
		{
			name: "no types",
			gen:  func() { c.GenTypeUsage("search", nil) },
			want: "Usage: ./trakt-sync search -t [type]\nAvailable types:\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, captureStdout(t, tt.gen))
		})
	}
}
