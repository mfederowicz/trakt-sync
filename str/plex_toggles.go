// Package str used for structs
package str

// PlexToggles represents JSON Plex toggles of one media type; each media type uses a subset
type PlexToggles struct {
	Watching  *bool `json:"watching,omitempty"`
	Watched   *bool `json:"watched,omitempty"`
	Rated     *bool `json:"rated,omitempty"`
	Collected *bool `json:"collected,omitempty"`
	Watchlist *bool `json:"watchlist,omitempty"`
}

func (p PlexToggles) String() string {
	return Stringify(p)
}
