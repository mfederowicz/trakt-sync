// Package str used for structs
package str

// SmartListFilters represents JSON smart list filters object; ranges are [min, max]
type SmartListFilters struct {
	Genres              []string  `json:"genres,omitempty"`
	GenresOperator      *string   `json:"genres_operator,omitempty"`
	Subgenres           []string  `json:"subgenres,omitempty"`
	Certifications      []string  `json:"certifications,omitempty"`
	Languages           []string  `json:"languages,omitempty"`
	Countries           []string  `json:"countries,omitempty"`
	Statuses            []string  `json:"statuses,omitempty"`
	Networks            []string  `json:"networks,omitempty"`
	Keywords            []string  `json:"keywords,omitempty"`
	KeywordsOperator    *string   `json:"keywords_operator,omitempty"`
	WatchNow            []string  `json:"watchnow,omitempty"`
	Years               []int     `json:"years,omitempty"`
	Ratings             []int     `json:"ratings,omitempty"`
	Runtimes            []int     `json:"runtimes,omitempty"`
	ImdbRatings         []float64 `json:"imdb_ratings,omitempty"`
	RtMeters            []int     `json:"rt_meters,omitempty"`
	RtUserMeters        []int     `json:"rt_user_meters,omitempty"`
	LetterboxdRatings   []float64 `json:"letterboxd_ratings,omitempty"`
	MalRatings          []float64 `json:"mal_ratings,omitempty"`
	IgnoreWatched       *bool     `json:"ignore_watched,omitempty"`
	IgnoreWatchlisted   *bool     `json:"ignore_watchlisted,omitempty"`
	IgnoreWatching      *bool     `json:"ignore_watching,omitempty"`
	IgnoreUnreleased    *bool     `json:"ignore_unreleased,omitempty"`
	IgnoreReleased      *bool     `json:"ignore_released,omitempty"`
	IgnoreEnded         *bool     `json:"ignore_ended,omitempty"`
	IgnoreAiring        *bool     `json:"ignore_airing,omitempty"`
	IgnoreNoReleaseDate *bool     `json:"ignore_no_release_date,omitempty"`
}

func (s SmartListFilters) String() string {
	return Stringify(s)
}
