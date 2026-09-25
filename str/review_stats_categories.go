// Package str used for structs
package str

// ReviewStatsCategories represents JSON review stats per category; lists_counts is sent for all media only
type ReviewStatsCategories struct {
	Minutes         *ReviewStat `json:"minutes,omitempty"`
	PlayCounts      *ReviewStat `json:"play_counts,omitempty"`
	CollectedCounts *ReviewStat `json:"collected_counts,omitempty"`
	RatingsCounts   *ReviewStat `json:"ratings_counts,omitempty"`
	CommentsCounts  *ReviewStat `json:"comments_counts,omitempty"`
	ListsCounts     *ReviewStat `json:"lists_counts,omitempty"`
}

func (r ReviewStatsCategories) String() string {
	return Stringify(r)
}
