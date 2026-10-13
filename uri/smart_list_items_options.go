package uri

// SmartListItemsOptions query options for smart-lists/{list_id}/items; parental ranges are min-max, from 0 (none) to 3 (severe).
// A value of a list filter (genres, studios, people, watchnow and similar) with a - prefix is excluded.
type SmartListItemsOptions struct {
	Certifications         string `url:"certifications,omitempty"`
	Countries              string `url:"countries,omitempty"`
	Extended               string `url:"extended,omitempty"`
	Genres                 string `url:"genres,omitempty"`
	IgnoreWatched          string `url:"ignore_watched,omitempty"`
	IgnoreWatchlisted      string `url:"ignore_watchlisted,omitempty"`
	ImdbVotes              string `url:"imdb_votes,omitempty"`
	Limit                  int    `url:"limit,omitempty"`
	Metascores             string `url:"metascores,omitempty"`
	Page                   int    `url:"page,omitempty"`
	ParentalAlcohol        string `url:"parental_alcohol,omitempty"`
	ParentalFrightening    string `url:"parental_frightening,omitempty"`
	ParentalIncludeUnrated bool   `url:"parental_include_unrated,omitempty"`
	ParentalNudity         string `url:"parental_nudity,omitempty"`
	ParentalProfanity      string `url:"parental_profanity,omitempty"`
	ParentalViolence       string `url:"parental_violence,omitempty"`
	People                 string `url:"people,omitempty"`
	PeopleOperator         string `url:"people_operator,omitempty"`
	PeopleRole             string `url:"people_role,omitempty"`
	Ratings                string `url:"ratings,omitempty"`
	ReleasedWithinDays     int    `url:"released_within_days,omitempty"`
	Runtimes               string `url:"runtimes,omitempty"`
	Studios                string `url:"studios,omitempty"`
	Subgenres              string `url:"subgenres,omitempty"`
	TmdbRatings            string `url:"tmdb_ratings,omitempty"`
	Votes                  string `url:"votes,omitempty"`
	WatchNow               string `url:"watchnow,omitempty"`
	WatchNowCountry        string `url:"watchnow_country,omitempty"`
	Years                  string `url:"years,omitempty"`
}
