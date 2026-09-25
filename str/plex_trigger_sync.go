// Package str used for structs
package str

// PlexTriggerSync represents JSON Plex sync to start with a settings update; *_all_data re-pulls full data
type PlexTriggerSync struct {
	WatchedAllData    *bool `json:"watched_all_data,omitempty"`
	CollectionAllData *bool `json:"collection_all_data,omitempty"`
	RatingsAllData    *bool `json:"ratings_all_data,omitempty"`
	WatchlistAllData  *bool `json:"watchlist_all_data,omitempty"`
}

func (p PlexTriggerSync) String() string {
	return Stringify(p)
}
