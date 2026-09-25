// Package str used for structs
package str

// PlexWebhook represents JSON Plex real-time scrobbler webhook; url is null unless the user is VIP
type PlexWebhook struct {
	URL         *string    `json:"url,omitempty"`
	LastEventAt *Timestamp `json:"last_event_at,omitempty"`
	EventCount  *int       `json:"event_count,omitempty"`
	HomeUsers   *string    `json:"home_users,omitempty"`
}

func (p PlexWebhook) String() string {
	return Stringify(p)
}
