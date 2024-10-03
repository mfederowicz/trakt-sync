// Package writer used for write operations
package writer

import (
	"log"
	"os"

	"github.com/mfederowicz/trakt-sync/str"
)

// WriteJSON write results to file
func WriteJSON(options *str.Options, results []byte) {
	err := os.WriteFile(options.Output, results, os.ModePerm)
	if err != nil {
		log.Println("write error")
	}
}
