package str

import (
	"encoding/json"
	"fmt"
	"strings"
)

// sourceSeparator splits collected_at from the source in a minimal collection value
const sourceSeparator = "|"

// MinimalCollected is one value of a minimal collection: collected_at, followed by "|source" when the request
// sent available_on (source is plex, jellyfin or empty). It is encoded back in the form it was decoded from.
type MinimalCollected struct {
	CollectedAt *Timestamp
	// Source is nil for a plain collected_at; with available_on it is plex, jellyfin or empty
	Source *string
}

func (m MinimalCollected) String() string {
	return Stringify(m)
}

// UnmarshalJSON decodes "collected_at" or "collected_at|source"
func (m *MinimalCollected) UnmarshalJSON(b []byte) error {
	// A JSON null leaves the value as it is, like Timestamp does
	if string(b) == jsonNull {
		return nil
	}
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return fmt.Errorf("minimal collection value: %w", err)
	}
	collectedAt, source, found := strings.Cut(value, sourceSeparator)
	quoted, err := json.Marshal(collectedAt)
	if err != nil {
		return fmt.Errorf("minimal collection value: %w", err)
	}
	ts := new(Timestamp)
	if err := ts.UnmarshalJSON(quoted); err != nil {
		return fmt.Errorf("minimal collection value %q: %w", value, err)
	}
	m.CollectedAt = ts
	m.Source = nil
	if found {
		m.Source = String(source)
	}
	return nil
}

// MarshalJSON encodes collected_at, followed by "|source" when the value had a source part
func (m MinimalCollected) MarshalJSON() ([]byte, error) {
	if m.CollectedAt == nil {
		return []byte(jsonNull), nil
	}
	collectedAt, err := m.CollectedAt.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if m.Source == nil {
		return collectedAt, nil
	}
	var value string
	if err := json.Unmarshal(collectedAt, &value); err != nil {
		return nil, err
	}
	return json.Marshal(value + sourceSeparator + *m.Source)
}
