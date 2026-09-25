// Package str used for structs
package str

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestSyncItemKeepsPassthroughKeys checks that keys outside the contract survive a decode and encode
// (timestamps use the Z form, as str.Timestamp writes them).
func TestSyncItemKeepsPassthroughKeys(t *testing.T) {
	in := `{"kind":"history","type":"movie","trakt_item":{"type":"movie","title":"Arrival","year":2016,"ids":{"trakt":1,"slug":"arrival-2016"}},` +
		`"service_id":"netflix","tmdb_id":329865,"watched_at":"2026-09-01T20:00:00Z","progress":87.5,"netflix_title":"Arrival","device":{"name":"TV"}}`

	item := new(SyncItem)
	if err := json.Unmarshal([]byte(in), item); err != nil {
		t.Fatal(err)
	}
	if item.Kind == nil || *item.Kind != "history" || item.TmdbID == nil || *item.TmdbID != 329865 || item.TraktItem == nil || *item.TraktItem.Title != "Arrival" {
		t.Fatalf("contract fields not decoded: %v", item)
	}
	if len(item.Extra) != 2 {
		t.Fatalf("Extra has %d keys, want 2 (netflix_title, device): %v", len(item.Extra), item.Extra)
	}

	out, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	var want, got map[string]any
	if err := json.Unmarshal([]byte(in), &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("round trip changed the item (-want +got):\n%s", diff)
	}
}

// TestSyncItemWithoutExtra checks that an item with only contract keys has no Extra map.
func TestSyncItemWithoutExtra(t *testing.T) {
	item := new(SyncItem)
	if err := json.Unmarshal([]byte(`{"kind":"rating","rating_value":8}`), item); err != nil {
		t.Fatal(err)
	}
	if item.Extra != nil {
		t.Errorf("Extra is %v, want nil", item.Extra)
	}
	out, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{"kind":"rating","rating_value":8}` {
		t.Errorf("marshal is %s", out)
	}
}
