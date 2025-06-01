// Package str used for structs
package str

// UserLastActivities represents JSON user activities object
type UserLastActivities struct {
	All          *Timestamp   `json:"all,omitempty"`
	Movies       *Movies      `json:"movies,omitempty"`
	Episodes     *Episodes    `json:"episodes,omitempty"`
	Shows        *Shows       `json:"shows,omitempty"`
	Seasons      *Seasons     `json:"seasons,omitempty"`
	Comments     *Comments    `json:"comments,omitempty"`
	Lists        *Lists       `json:"lists,omitempty"`
	Watchlist    *Watchlist   `json:"watchlist,omitempty"`
	Favorites    *Favorites   `json:"favorites,omitempty"`
	Account      *Account     `json:"account,omitempty"`
	SavedFilters *SavedFilter `json:"saved_filters,omitempty"`
	Notes        *Notes       `json:"notes,omitempty"`
	Connections  *Connections `json:"connections,omitempty"`
	SharingText  *SharingText `json:"sharing_text,omitempty"`
	Limits       *Limits      `json:"limits,omitempty"`
}

func (u UserLastActivities) String() string {
	return Stringify(u)
}
