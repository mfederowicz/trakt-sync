package uri

// SmartListItemsOptions query options for smart-lists/{list_id}/items; parental ranges are min-max, from 0 (none) to 3 (severe)
type SmartListItemsOptions struct {
	Certifications         string `url:"certifications,omitempty"`
	Countries              string `url:"countries,omitempty"`
	Extended               string `url:"extended,omitempty"`
	Genres                 string `url:"genres,omitempty"`
	IgnoreWatched          string `url:"ignore_watched,omitempty"`
	IgnoreWatchlisted      string `url:"ignore_watchlisted,omitempty"`
	Limit                  int    `url:"limit,omitempty"`
	Page                   int    `url:"page,omitempty"`
	ParentalAlcohol        string `url:"parental_alcohol,omitempty"`
	ParentalFrightening    string `url:"parental_frightening,omitempty"`
	ParentalIncludeUnrated bool   `url:"parental_include_unrated,omitempty"`
	ParentalNudity         string `url:"parental_nudity,omitempty"`
	ParentalProfanity      string `url:"parental_profanity,omitempty"`
	ParentalViolence       string `url:"parental_violence,omitempty"`
	Ratings                string `url:"ratings,omitempty"`
	Runtimes               string `url:"runtimes,omitempty"`
	Subgenres              string `url:"subgenres,omitempty"`
	WatchNow               string `url:"watchnow,omitempty"`
	WatchNowCountry        string `url:"watchnow_country,omitempty"`
	Years                  string `url:"years,omitempty"`
}
