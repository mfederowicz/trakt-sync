// Package writer used for write operations
package writer

import (
	"github.com/mfederowicz/trakt-sync/str"
	"os"
)

// WriteJSON write results to file
func WriteJSON(options *str.Options, results []byte) {
	os.WriteFile(options.Output, results, os.ModePerm)
}
