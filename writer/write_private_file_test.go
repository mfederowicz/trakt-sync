// Package writer used for write operations
package writer

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

func assertPrivateFile(t *testing.T, path string, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != want {
		t.Errorf("content = %q, want %q", data, want)
	}
	if runtime.GOOS == "windows" {
		return
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if got := info.Mode().Perm(); got != consts.X600 {
		t.Errorf("mode = %o, want %o", got, consts.X600)
	}
}

// TestWritePrivateFile checks new and existing (0644) files end up 0600.
func TestWritePrivateFile(t *testing.T) {
	tests := []struct {
		name     string
		existing bool
	}{
		{name: "new file"},
		{name: "existing 0644 file", existing: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "token.json")
			if tt.existing {
				if err := os.WriteFile(path, []byte("old"), consts.X644); err != nil {
					t.Fatal(err)
				}
			}
			if err := WritePrivateFile(path, []byte("{}")); err != nil {
				t.Fatal(err)
			}
			assertPrivateFile(t, path, "{}")
		})
	}
}

// TestWriteJSONPrivate checks exports are written 0600.
func TestWriteJSONPrivate(t *testing.T) {
	options := &str.Options{Output: filepath.Join(t.TempDir(), "export.json")}
	WriteJSON(options, []byte("[]"))
	assertPrivateFile(t, options.Output, "[]")
}
