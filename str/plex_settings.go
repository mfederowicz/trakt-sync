// Package str used for structs
package str

// PlexSettings represents JSON Plex settings: connection, scrobbler webhook, sync selection and toggles
type PlexSettings struct {
	Connection *PlexConnection `json:"connection,omitempty"`
	Webhook    *PlexWebhook    `json:"webhook,omitempty"`
	Sync       *PlexSync       `json:"sync,omitempty"`
	Scrobbler  *PlexScrobbler  `json:"scrobbler,omitempty"`
}

func (p PlexSettings) String() string {
	return Stringify(p)
}
