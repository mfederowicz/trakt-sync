package writer

import (
	"os"
	"trakt-sync/str"
)

func WriteJson(options *str.Options, results []byte) {
	os.WriteFile(options.Output, results, os.ModePerm)
}
