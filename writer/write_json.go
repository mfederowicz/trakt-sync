package writer

import (
	"os"
	"github.com/mfederowicz/trakt-sync/str"
)

func WriteJson(options *str.Options, results []byte) {
	os.WriteFile(options.Output, results, os.ModePerm)
}
