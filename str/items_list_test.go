package str

import (
	"reflect"
	"testing"
)

func watchedItem(id int64, day int) ExportlistItem {
	return ExportlistItem{IDs: &IDs{Trakt: Int64(id)}, WatchedAt: testStamp(day)}
}

func ratedItem(id int64, day int) ExportlistItem {
	return ExportlistItem{IDs: &IDs{Trakt: Int64(id)}, RatedAt: testStamp(day)}
}

// daysByID maps each Trakt ID to the day of the date picked by day.
func daysByID(items *[]ExportlistItem, day func(ExportlistItem) *Timestamp) map[int64]int {
	out := map[int64]int{}
	for _, item := range *items {
		out[*item.IDs.Trakt] = day(item).Day()
	}
	return out
}

func TestItemsListGetUniqueOldest(t *testing.T) {
	watched := func(i ExportlistItem) *Timestamp { return i.WatchedAt }
	rated := func(i ExportlistItem) *Timestamp { return i.RatedAt }

	t.Run("one entry per id with the oldest watched_at", func(t *testing.T) {
		items := &[]ExportlistItem{watchedItem(1, 4), watchedItem(2, 5), watchedItem(1, 2), watchedItem(1, 9)}
		got := ItemsList{}.GetUniqueOldest(items)
		if want := (map[int64]int{1: 2, 2: 5}); !reflect.DeepEqual(daysByID(got, watched), want) {
			t.Errorf("days by id are %v, want %v", daysByID(got, watched), want)
		}
	})

	t.Run("one entry per id with the oldest rated_at", func(t *testing.T) {
		items := &[]ExportlistItem{ratedItem(1, 8), ratedItem(1, 3), ratedItem(3, 1)}
		got := ItemsList{}.GetUniqueOldest(items)
		if want := (map[int64]int{1: 3, 3: 1}); !reflect.DeepEqual(daysByID(got, rated), want) {
			t.Errorf("days by id are %v, want %v", daysByID(got, rated), want)
		}
	})

	t.Run("items without watched_at and rated_at are kept", func(t *testing.T) {
		undated := ExportlistItem{IDs: &IDs{Trakt: Int64(1)}}
		got := ItemsList{}.GetUniqueOldest(&[]ExportlistItem{undated, undated, watchedItem(2, 5)})
		if len(*got) != 2 {
			t.Fatalf("result is %v, want two items", *got)
		}
		for _, item := range *got {
			if *item.IDs.Trakt == 1 && item.WatchedAt != nil {
				t.Errorf("item 1 is %v, want it without watched_at", item)
			}
		}
	})

	t.Run("an item with a date wins over one without", func(t *testing.T) {
		undated := ExportlistItem{IDs: &IDs{Trakt: Int64(1)}}
		for _, items := range []*[]ExportlistItem{
			{undated, watchedItem(1, 6)},
			{watchedItem(1, 6), undated},
		} {
			got := ItemsList{}.GetUniqueOldest(items)
			if want := (map[int64]int{1: 6}); len(*got) != 1 || (*got)[0].WatchedAt == nil || !reflect.DeepEqual(daysByID(got, watched), want) {
				t.Errorf("result is %v, want the item watched on day 6", *got)
			}
		}
	})

	t.Run("watched and rated items with the same id", func(t *testing.T) {
		for _, items := range []*[]ExportlistItem{
			{watchedItem(1, 6), ratedItem(1, 3)},
			{ratedItem(1, 3), watchedItem(1, 6)},
		} {
			got := ItemsList{}.GetUniqueOldest(items)
			if len(*got) != 1 || (*got)[0].RatedAt == nil || (*got)[0].RatedAt.Day() != 3 {
				t.Errorf("result is %v, want the item rated on day 3", *got)
			}
		}
	})

	t.Run("items without a trakt id are dropped", func(t *testing.T) {
		items := &[]ExportlistItem{{WatchedAt: testStamp(1)}, {IDs: &IDs{}, WatchedAt: testStamp(1)}, watchedItem(2, 5)}
		got := ItemsList{}.GetUniqueOldest(items)
		if want := (map[int64]int{2: 5}); !reflect.DeepEqual(daysByID(got, watched), want) {
			t.Errorf("days by id are %v, want %v", daysByID(got, watched), want)
		}
	})

	t.Run("nil list", func(t *testing.T) {
		if got := (ItemsList{}).GetUniqueOldest(nil); got != nil {
			t.Errorf("result is %v, want nil", got)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		got := ItemsList{}.GetUniqueOldest(&[]ExportlistItem{})
		if got == nil || len(*got) != 0 {
			t.Errorf("result is %v, want an empty list", got)
		}
	})
}

func TestItemsListGetUniqIDs(t *testing.T) {
	got := ItemsList{}.GetUniqIDs(&[]int64{3, 1, 3, 2, 1})
	if want := []int64{3, 1, 2}; !reflect.DeepEqual(*got, want) {
		t.Errorf("ids are %v, want %v", *got, want)
	}
}

func TestItemsListGetUniqIDsNil(t *testing.T) {
	if got := (ItemsList{}).GetUniqIDs(nil); got != nil {
		t.Errorf("ids are %v, want nil", got)
	}
}

func TestItemsListUniqNilLists(t *testing.T) {
	got := ItemsList{Shows: &[]ExportlistItem{watchedItem(2, 2), watchedItem(2, 1)}}.Uniq()

	if got.Movies != nil || got.Seasons != nil || got.Episodes != nil || got.IDs != nil {
		t.Errorf("nil lists are %v, want them to stay nil", got)
	}
	if len(*got.Shows) != 1 {
		t.Errorf("shows are %v, want one", *got.Shows)
	}
}

func TestItemsListUniq(t *testing.T) {
	watched := func(i ExportlistItem) *Timestamp { return i.WatchedAt }
	list := ItemsList{
		Movies:   &[]ExportlistItem{watchedItem(1, 2), watchedItem(1, 3)},
		Shows:    &[]ExportlistItem{watchedItem(2, 2)},
		Seasons:  &[]ExportlistItem{watchedItem(3, 2), watchedItem(3, 1)},
		Episodes: &[]ExportlistItem{},
		IDs:      &[]int64{7, 7, 8},
	}

	got := list.Uniq()

	if want := (map[int64]int{1: 2}); !reflect.DeepEqual(daysByID(got.Movies, watched), want) {
		t.Errorf("movies are %v, want %v", daysByID(got.Movies, watched), want)
	}
	if want := (map[int64]int{2: 2}); !reflect.DeepEqual(daysByID(got.Shows, watched), want) {
		t.Errorf("shows are %v, want %v", daysByID(got.Shows, watched), want)
	}
	if want := (map[int64]int{3: 1}); !reflect.DeepEqual(daysByID(got.Seasons, watched), want) {
		t.Errorf("seasons are %v, want %v", daysByID(got.Seasons, watched), want)
	}
	if len(*got.Episodes) != 0 {
		t.Errorf("episodes are %v, want none", *got.Episodes)
	}
	if want := []int64{7, 8}; !reflect.DeepEqual(*got.IDs, want) {
		t.Errorf("ids are %v, want %v", *got.IDs, want)
	}
	if len(*list.Movies) != 2 {
		t.Errorf("Uniq changed the original list: %v", *list.Movies)
	}
}
