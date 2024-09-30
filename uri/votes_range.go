// Package uri used for url operations
package uri

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
)

// VotesRange represents min/max int votes parameters
type VotesRange struct {
	Min int `url:"min,omitempty"`
	Max int `url:"max,omitempty"`
}

func (r VotesRange) String() string {
	if r.Min <= consts.VotesRangeMin || r.Max > consts.VotesRangeMax {
		return consts.EmptyString
	}

	if r.Min > r.Max {
		return consts.EmptyString
	}

	return fmt.Sprintf(consts.RangeFormatDigits, r.Min, r.Max)
}
