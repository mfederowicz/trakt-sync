// Package writer used for write operations
package writer

import (
	"log"

	"github.com/mfederowicz/trakt-sync/str"
)

// WriteJSON write results to file
func WriteJSON(options *str.Options, results []byte) {
	err := WritePrivateFile(options.Output, results)
	if err != nil {
		log.Println("write error")
	}
}
