// Package str used for structs
package str

// SettingsGenres represents JSON favorite and disliked genre slugs
type SettingsGenres struct {
	Favorites []string `json:"favorites,omitempty"`
	Disliked  []string `json:"disliked,omitempty"`
}

func (s SettingsGenres) String() string {
	return Stringify(s)
}
