package uri

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
)

// RatingRange represents min/max int parameters
type RatingRange struct {
	Min int `url:"min,omitempty"`
	Max int `url:"max,omitempty"`
}

func (rr RatingRange) String() string {
	if rr.Min <= consts.RatingRageMin || rr.Max > consts.RatingRangeMax {
		return ""
	}

	if rr.Min > rr.Max {
		return ""
	}
	return fmt.Sprintf(consts.RangeFormatDigits, rr.Min, rr.Max)
}
