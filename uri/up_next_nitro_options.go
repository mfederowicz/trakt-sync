// Package uri used for url operations
package uri

// UpNextNitroOptions query options for sync/progress/up_next_nitro
type UpNextNitroOptions struct {
	Certifications string `url:"certifications,omitempty"`
	Countries      string `url:"countries,omitempty"`
	EndDate        string `url:"end_date,omitempty"`
	Genres         string `url:"genres,omitempty"`
	Intent         string `url:"intent,omitempty"`
	Limit          int    `url:"limit,omitempty"`
	Page           int    `url:"page,omitempty"`
	Ratings        string `url:"ratings,omitempty"`
	Runtimes       string `url:"runtimes,omitempty"`
	SortBy         string `url:"sort_by,omitempty"`
	SortHow        string `url:"sort_how,omitempty"`
	StartDate      string `url:"start_date,omitempty"`
	Subgenres      string `url:"subgenres,omitempty"`
	WatchNow       string `url:"watchnow,omitempty"`
	Years          string `url:"years,omitempty"`
}
