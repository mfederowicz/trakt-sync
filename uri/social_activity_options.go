package uri

// SocialActivityOptions query options for users/{id}/{type}/activities
type SocialActivityOptions struct {
	Certifications    string `url:"certifications,omitempty"`
	Countries         string `url:"countries,omitempty"`
	EndDate           string `url:"end_date,omitempty"`
	Extended          string `url:"extended,omitempty"`
	Genres            string `url:"genres,omitempty"`
	IgnoreCollected   string `url:"ignore_collected,omitempty"`
	IgnoreWatched     string `url:"ignore_watched,omitempty"`
	IgnoreWatchlisted string `url:"ignore_watchlisted,omitempty"`
	Limit             int    `url:"limit,omitempty"`
	Page              int    `url:"page,omitempty"`
	// ParentalAlcohol and the other parental guide filters take a severity range from 0 (none) to 3 (severe), ie: 0-1
	ParentalAlcohol        string `url:"parental_alcohol,omitempty"`
	ParentalFrightening    string `url:"parental_frightening,omitempty"`
	ParentalIncludeUnrated bool   `url:"parental_include_unrated,omitempty"`
	ParentalNudity         string `url:"parental_nudity,omitempty"`
	ParentalProfanity      string `url:"parental_profanity,omitempty"`
	ParentalViolence       string `url:"parental_violence,omitempty"`
	Ratings                string `url:"ratings,omitempty"`
	Runtimes               string `url:"runtimes,omitempty"`
	StartDate              string `url:"start_date,omitempty"`
	Subgenres              string `url:"subgenres,omitempty"`
	WatchNow               string `url:"watchnow,omitempty"`
	Years                  string `url:"years,omitempty"`
}
