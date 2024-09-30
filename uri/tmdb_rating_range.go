// Package uri used for url operations
package uri

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
)

// TmdbRatingRange represents min/max float tmdb rating parameters
type TmdbRatingRange struct {
	Min float32 `url:"min,omitempty"`
	Max float32 `url:"max,omitempty"`
}

func (r TmdbRatingRange) String() string {
	if r.Min < consts.TmdbRatingRangeMin || r.Max > consts.TmdbRatingRangeMax {
		return consts.EmptyString
	}

	if r.Min > r.Max || r.Min == r.Max {
		return consts.EmptyString
	}

	return fmt.Sprintf(consts.RangeFormatFloats, r.Min, r.Max)
}

