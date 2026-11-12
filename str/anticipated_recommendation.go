package str

// AnticipatedRecommendation represents JSON anticipated recommendation: an unreleased movie or show ranked for the user, with its release date (YYYY-MM-DD) when known
type AnticipatedRecommendation struct {
	Score    *float64 `json:"score,omitempty"`
	Released *string  `json:"released,omitempty"`
	Movie    *Movie   `json:"movie,omitempty"`
	Show     *Show    `json:"show,omitempty"`
}

func (a AnticipatedRecommendation) String() string {
	return Stringify(a)
}
