// Package str used for structs
package str

import (
	"github.com/mfederowicz/trakt-sync/consts"
)

// ItemsList represents JSON items object
type ItemsList struct {
	Movies   *[]ExportlistItem `json:"movies,omitempty"`
	Shows    *[]ExportlistItem `json:"shows,omitempty"`
	Seasons  *[]ExportlistItem `json:"seasons,omitempty"`
	Episodes *[]ExportlistItem `json:"episodes,omitempty"`
	IDs      *[]int64          `json:"ids,omitempty"`
}

func (i ItemsList) String() string {
	return Stringify(i)
}

// Uniq make lists unique with oldest elements
func (i ItemsList) Uniq() *ItemsList {

	i.Movies = i.GetUniqueOldest(i.Movies)
	i.Shows = i.GetUniqueOldest(i.Shows)
	i.Seasons = i.GetUniqueOldest(i.Seasons)
	i.Episodes = i.GetUniqueOldest(i.Episodes)
	i.IDs = i.GetUniqIDs(i.IDs)
	return &i
}

// GetUniqueOldest returns a unique slice of Items, keeping the one with the oldest WatchedAt per ID.
func (ItemsList) GetUniqueOldest(items *[]ExportlistItem) *[]ExportlistItem {
	unique := make(map[int64]ExportlistItem)

	for _, item := range *items {
		id := *item.IDs.Trakt
		if item.WatchedAt == nil {
			continue // skip items with nil ID or WatchedAt
		}
		existing, found := unique[id]
		if !found || item.WatchedAt.After(existing.WatchedAt.Time) {
			unique[id] = item
		}
	}
	result := make([]ExportlistItem, consts.ZeroValue, len(unique))
	for _, item := range unique {
		result = append(result, item)
	}
	return &result
}

// GetUniqIDs returns a unique slice of ints.
func (ItemsList) GetUniqIDs(input *[]int64) *[]int64 {
	seen := make(map[int64]struct{}, len(*input))
	uniq := make([]int64, consts.ZeroValue, len(*input))

	for _, v := range *input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			uniq = append(uniq, v)
		}
	}

	return &uniq
}
