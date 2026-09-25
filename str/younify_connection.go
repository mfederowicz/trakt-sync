// Package str used for structs
package str

// YounifyConnection represents JSON streaming service connection object
type YounifyConnection struct {
	ID           *string               `json:"id,omitempty"`
	Name         *string               `json:"name,omitempty"`
	Vip          *bool                 `json:"vip,omitempty"`
	Color        *string               `json:"color,omitempty"`
	Images       *WatchNowSourceImages `json:"images,omitempty"`
	Connectable  *bool                 `json:"connectable,omitempty"`
	Connected    *bool                 `json:"connected,omitempty"`
	Active       *bool                 `json:"active,omitempty"`
	Profile      *string               `json:"profile,omitempty"`
	LastSyncedAt *Timestamp            `json:"last_synced_at,omitempty"`
}

func (y YounifyConnection) String() string {
	return Stringify(y)
}
