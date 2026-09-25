// Package str used for structs
package str

// PlexSettingsUpdate represents JSON Plex settings update; omitted keys are left unchanged, trigger_sync enqueues a sync
type PlexSettingsUpdate struct {
	Sync        *PlexSyncUpdate    `json:"sync,omitempty"`
	Scrobbler   *PlexScrobbler     `json:"scrobbler,omitempty"`
	Webhook     *PlexWebhookUpdate `json:"webhook,omitempty"`
	TriggerSync *PlexTriggerSync   `json:"trigger_sync,omitempty"`
}

func (p PlexSettingsUpdate) String() string {
	return Stringify(p)
}
