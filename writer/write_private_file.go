// Package writer used for write operations
package writer

import (
	"os"

	"github.com/mfederowicz/trakt-sync/consts"
)

// WritePrivateFile writes data readable only by the owner (0600), also tightening a file that already exists
func WritePrivateFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, consts.X600); err != nil {
		return err
	}

	// os.WriteFile keeps the mode of an existing file, so an old 0644 file is fixed here
	return os.Chmod(path, consts.X600)
}
