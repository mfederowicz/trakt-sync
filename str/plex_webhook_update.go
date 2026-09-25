// Package str used for structs
package str

// PlexWebhookUpdate represents JSON Plex webhook values to change
type PlexWebhookUpdate struct {
	HomeUsers *string `json:"home_users,omitempty"`
}

func (p PlexWebhookUpdate) String() string {
	return Stringify(p)
}
