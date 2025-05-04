// Package str used for structs
package str

// UserWatched represents JSON user watched object
type UserWatched struct {
	Plays         *int       `json:"plays,omitempty"`
	LastWatchedAt *Timestamp `json:"last_watched_at,omitempty"`
	LastUpdatedAt *Timestamp `json:"last_updated_at,omitempty"`
	ResetAt       *Timestamp `json:"reset_at,omitempty"`
	Movie         *Movie     `json:"movie,omitempty"`
	Show          *Show      `json:"show,omitempty"`
	Seasons       *[]Season  `json:"seasons,omitempty"`
}

func (u UserWatched) String() string {
	return Stringify(u)
}
