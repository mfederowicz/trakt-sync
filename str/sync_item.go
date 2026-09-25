// Package str used for structs
package str

import (
	"encoding/json"

	"github.com/mfederowicz/trakt-sync/consts"
)

// SyncItem represents JSON paused or skipped item of a data sync. The API passes the raw stored item through,
// so keys outside the contract are kept in Extra and written back by MarshalJSON.
type SyncItem struct {
	Kind         *string                    `json:"kind,omitempty"`
	Type         *string                    `json:"type,omitempty"`
	TraktItem    *SyncTraktItem             `json:"trakt_item,omitempty"`
	ServiceID    *string                    `json:"service_id,omitempty"`
	ContentID    *string                    `json:"content_id,omitempty"`
	ProfileID    *string                    `json:"profile_id,omitempty"`
	TmdbID       *int                       `json:"tmdb_id,omitempty"`
	TmdbSeriesID *int                       `json:"tmdb_series_id,omitempty"`
	WatchedAt    *Timestamp                 `json:"watched_at,omitempty"`
	RatedAt      *Timestamp                 `json:"rated_at,omitempty"`
	Progress     *float64                   `json:"progress,omitempty"`
	RatingType   *string                    `json:"rating_type,omitempty"`
	RatingValue  *int                       `json:"rating_value,omitempty"`
	Extra        map[string]json.RawMessage `json:"-"`
}

// syncItemFields is SyncItem without its methods, so (un)marshalling it does not recurse.
type syncItemFields SyncItem

// UnmarshalJSON decodes the contract fields and keeps every other key in Extra.
func (s *SyncItem) UnmarshalJSON(data []byte) error {
	fields := syncItemFields{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	all := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &all); err != nil {
		return err
	}
	for _, key := range syncItemKeys {
		delete(all, key)
	}
	*s = SyncItem(fields)
	if len(all) > consts.ZeroValue {
		s.Extra = all
	}
	return nil
}

// MarshalJSON writes the contract fields together with the passed-through keys.
func (s SyncItem) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(syncItemFields(s))
	if err != nil || len(s.Extra) == consts.ZeroValue {
		return data, err
	}
	merged := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &merged); err != nil {
		return nil, err
	}
	for key, value := range s.Extra {
		if _, ok := merged[key]; !ok {
			merged[key] = value
		}
	}
	return json.Marshal(merged)
}

// syncItemKeys are the JSON keys of the contract fields.
var syncItemKeys = []string{"kind", "type", "trakt_item", "service_id", "content_id", "profile_id", "tmdb_id",
	"tmdb_series_id", "watched_at", "rated_at", "progress", "rating_type", "rating_value"}

func (s SyncItem) String() string {
	return Stringify(s)
}
