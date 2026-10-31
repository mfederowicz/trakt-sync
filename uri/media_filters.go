package uri

// MediaFilters are the media filters shared by the routes that list movies, shows and episodes.
// A filter that takes several values is a comma separated string.
type MediaFilters struct {
	Certifications string `url:"certifications,omitempty"`
	Countries      string `url:"countries,omitempty"`
	EndDate        string `url:"end_date,omitempty"`
	Genres         string `url:"genres,omitempty"`
	ImdbRatings    string `url:"imdb_ratings,omitempty"`
	Languages      string `url:"languages,omitempty"`
	// ParentalAlcohol and the other parental guide filters take a severity range from 0 (none) to 3 (severe), ie: 0-1
	ParentalAlcohol        string `url:"parental_alcohol,omitempty"`
	ParentalFrightening    string `url:"parental_frightening,omitempty"`
	ParentalIncludeUnrated bool   `url:"parental_include_unrated,omitempty"`
	ParentalNudity         string `url:"parental_nudity,omitempty"`
	ParentalProfanity      string `url:"parental_profanity,omitempty"`
	ParentalViolence       string `url:"parental_violence,omitempty"`
	Ratings                string `url:"ratings,omitempty"`
	RtMeters               string `url:"rt_meters,omitempty"`
	RtUserMeters           string `url:"rt_user_meters,omitempty"`
	Runtimes               string `url:"runtimes,omitempty"`
	StartDate              string `url:"start_date,omitempty"`
	Subgenres              string `url:"subgenres,omitempty"`
	WatchNow               string `url:"watchnow,omitempty"`
	Years                  string `url:"years,omitempty"`
}
