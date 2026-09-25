// Package str used for structs
package str

// PlexSyncRequest represents JSON Plex sync now request; no server_id syncs every selected server
type PlexSyncRequest struct {
	ServerID *string `json:"server_id,omitempty"`
	AllData  *bool   `json:"all_data,omitempty"`
}

func (p PlexSyncRequest) String() string {
	return Stringify(p)
}
