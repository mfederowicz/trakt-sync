package uri

// MediaFilters are the media filters shared by the routes that list movies, shows and episodes.
// A filter that takes several values is a comma separated string.
type MediaFilters struct {
	Certifications string `url:"certifications,omitempty"`
	Countries      string `url:"countries,omitempty"`
	EndDate        string `url:"end_date,omitempty"`
	Genres         string `url:"genres,omitempty"`
	Ratings        string `url:"ratings,omitempty"`
	Runtimes       string `url:"runtimes,omitempty"`
	StartDate      string `url:"start_date,omitempty"`
	Subgenres      string `url:"subgenres,omitempty"`
	WatchNow       string `url:"watchnow,omitempty"`
	Years          string `url:"years,omitempty"`
}
