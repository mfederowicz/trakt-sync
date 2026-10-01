// Package writer used for write operations
package writer

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// WritePrivateJSON encodes v as JSON and writes it with WritePrivateFile. On an encoding error the file is
// left untouched, so a stored token or settings file is never replaced by an empty one.
func WritePrivateJSON(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encode %s: %w", filepath.Base(path), err)
	}
	return WritePrivateFile(path, data)
}
