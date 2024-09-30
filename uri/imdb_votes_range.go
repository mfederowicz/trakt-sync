// Package uri used for url operations
package uri

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
)

// ImdbVotesRange represents min/max int imdb votes parameters
type ImdbVotesRange struct {
	Min int `url:"min,omitempty"`
	Max int `url:"max,omitempty"`
}

func (r ImdbVotesRange) String() string {
	if r.Min <= consts.ImdbVotesRangeMin || r.Max > consts.ImdbVotesRangeMax {
		return consts.EmptyString
	}

	if r.Min > r.Max {
		return consts.EmptyString
	}

	return fmt.Sprintf(consts.RangeFormatDigits, r.Min, r.Max)
}

