// Package uri used for url operations
package uri

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
)

// RatingRangeFloat represents min/max float parameters
type RatingRangeFloat struct {
	Min float32 `url:"min,omitempty"`
	Max float32 `url:"max,omitempty"`
}

func (rr RatingRangeFloat) String() string {
	if rr.Min <= consts.RatingRageMinFloat || rr.Max > consts.RatingRangeMaxFloat {
		return consts.EmptyString
	}

	if rr.Min > rr.Max {
		return consts.EmptyString
	}
	return fmt.Sprintf(consts.RangeFormatFloats, rr.Min, rr.Max)
}

