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
	Users    *[]ExportlistItem `json:"users,omitempty"`
	Lists    *[]PersonalList   `json:"lists,omitempty"`
	People   *[]ExportlistItem `json:"people,omitempty"`
	IDs      *[]int64          `json:"ids,omitempty"`
	List     *PersonalList     `json:"list,omitempty"`
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

// GetUniqueOldest returns a unique slice of Items, keeping the one with the oldest WatchedAt (or RatedAt) per ID.
// An item without both dates is kept unless another item with the same ID has one.
func (i ItemsList) GetUniqueOldest(items *[]ExportlistItem) *[]ExportlistItem {
	if items == nil {
		return nil
	}
	unique := map[int64]ExportlistItem{}
	for _, item := range *items {
		if item.IDs == nil || item.IDs.Trakt == nil {
			continue // skip items with nil ID
		}

		id := *item.IDs.Trakt
		existing, found := unique[id]
		if !found || i.isOlder(i.itemDate(item), i.itemDate(existing)) {
			unique[id] = item
		}
	}
	result := make([]ExportlistItem, consts.ZeroValue, len(unique))
	for _, item := range unique {
		result = append(result, item)
	}
	return &result
}

// itemDate returns the date an item is compared by: WatchedAt, or RatedAt when it was not watched.
func (ItemsList) itemDate(item ExportlistItem) *Timestamp {
	if item.WatchedAt != nil {
		return item.WatchedAt
	}
	return item.RatedAt
}

// isOlder tells if date should replace existing: any date replaces a missing one.
func (ItemsList) isOlder(date, existing *Timestamp) bool {
	if date == nil {
		return false
	}
	return existing == nil || date.Before(existing.Time)
}

// GetUniqIDs returns a unique slice of ints.
func (ItemsList) GetUniqIDs(input *[]int64) *[]int64 {
	if input == nil {
		return nil
	}
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
