// Package uri used for url operations
package uri

// SmartListItemsOptions query options for smart-lists/{list_id}/items
type SmartListItemsOptions struct {
	Certifications    string `url:"certifications,omitempty"`
	Countries         string `url:"countries,omitempty"`
	Extended          string `url:"extended,omitempty"`
	Genres            string `url:"genres,omitempty"`
	IgnoreWatched     string `url:"ignore_watched,omitempty"`
	IgnoreWatchlisted string `url:"ignore_watchlisted,omitempty"`
	Limit             int    `url:"limit,omitempty"`
	Page              int    `url:"page,omitempty"`
	Ratings           string `url:"ratings,omitempty"`
	Runtimes          string `url:"runtimes,omitempty"`
	Subgenres         string `url:"subgenres,omitempty"`
	WatchNow          string `url:"watchnow,omitempty"`
	Years             string `url:"years,omitempty"`
}
