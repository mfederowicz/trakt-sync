package str

import (
	"testing"
	"time"
)

func testStamp(day int) *Timestamp {
	return &Timestamp{Time: time.Date(2026, time.October, day, 12, 0, 0, 0, time.FixedZone("CEST", 2*60*60))}
}

func testMetadata() *Metadata {
	return &Metadata{
		MediaType:     String("bluray"),
		Resolution:    String("uhd_4k"),
		Audio:         String("dts"),
		AudioChannels: String("5.1"),
		ThreeD:        Bool(true),
	}
}

func TestExportlistItemJSONUptime(t *testing.T) {
	data := &ExportlistItem{
		WatchedAt:       testStamp(1),
		ListedAt:        testStamp(2),
		CollectedAt:     testStamp(3),
		LastCollectedAt: testStamp(4),
		UpdatedAt:       testStamp(5),
		LastUpdatedAt:   testStamp(6),
	}
	cases := []struct {
		field string
		got   func(i *ExportlistItemJSON) *Timestamp
		want  *Timestamp
	}{
		{field: "watched_at", got: func(i *ExportlistItemJSON) *Timestamp { return i.WatchedAt }, want: data.WatchedAt},
		{field: "listed_at", got: func(i *ExportlistItemJSON) *Timestamp { return i.ListedAt }, want: data.ListedAt},
		{field: "collected_at", got: func(i *ExportlistItemJSON) *Timestamp { return i.CollectedAt }, want: data.CollectedAt},
		{field: "last_collected_at", got: func(i *ExportlistItemJSON) *Timestamp { return i.LastCollectedAt }, want: data.LastCollectedAt},
		{field: "updated_at", got: func(i *ExportlistItemJSON) *Timestamp { return i.UpdatedAt }, want: data.UpdatedAt},
		{field: "last_updated_at", got: func(i *ExportlistItemJSON) *Timestamp { return i.LastUpdatedAt }, want: data.LastUpdatedAt},
	}
	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			item := &ExportlistItemJSON{}
			item.Uptime(&Options{Time: tc.field}, data)
			if tc.got(item) != tc.want {
				t.Errorf("%s is %v, want %v", tc.field, tc.got(item), tc.want)
			}
			if set := countTimes(item); set != 1 {
				t.Errorf("%d time fields set, want only %s", set, tc.field)
			}
		})
	}

	t.Run("unknown field", func(t *testing.T) {
		item := &ExportlistItemJSON{}
		item.Uptime(&Options{Time: "rated_at"}, data)
		if set := countTimes(item); set != 0 {
			t.Errorf("%d time fields set, want none", set)
		}
	})
}

func countTimes(i *ExportlistItemJSON) int {
	set := 0
	for _, ts := range []*Timestamp{i.WatchedAt, i.ListedAt, i.CollectedAt, i.LastCollectedAt, i.UpdatedAt, i.LastUpdatedAt} {
		if ts != nil {
			set++
		}
	}
	return set
}

func TestExportlistItemGetTime(t *testing.T) {
	// each case drops the field the previous one returned
	item := ExportlistItem{
		WatchedAt:       testStamp(1),
		ListedAt:        testStamp(2),
		UpdatedAt:       testStamp(3),
		LastUpdatedAt:   testStamp(4),
		CollectedAt:     testStamp(5),
		LastCollectedAt: testStamp(6),
	}
	steps := []struct {
		name string
		want *Timestamp
		drop func()
	}{
		{name: "watched_at", want: item.WatchedAt, drop: func() { item.WatchedAt = nil }},
		{name: "listed_at", want: item.ListedAt, drop: func() { item.ListedAt = nil }},
		{name: "updated_at", want: item.UpdatedAt, drop: func() { item.UpdatedAt = nil }},
		{name: "last_updated_at", want: item.LastUpdatedAt, drop: func() { item.LastUpdatedAt = nil }},
		{name: "collected_at", want: item.CollectedAt, drop: func() { item.CollectedAt = nil }},
		{name: "last_collected_at", want: item.LastCollectedAt, drop: func() { item.LastCollectedAt = nil }},
	}
	for _, step := range steps {
		if got := item.GetTime(); got != step.want {
			t.Errorf("GetTime() is %v, want %s %v", got, step.name, step.want)
		}
		step.drop()
	}
	if got := item.GetTime(); got != nil {
		t.Errorf("GetTime() without times is %v, want nil", got)
	}
}

func TestExportlistItemUpdateCollectedData(t *testing.T) {
	source := &ExportlistItem{
		WatchedAt:   testStamp(1),
		CollectedAt: testStamp(2),
		HiddenAt:    testStamp(3),
		Metadata:    testMetadata(),
		Seasons: &[]Season{
			{Number: Int(1), Episodes: &[]Episode{
				{Number: Int(1), CollectedAt: testStamp(4), Metadata: testMetadata()},
				{Number: Int(2), CollectedAt: testStamp(5), Metadata: &Metadata{}},
			}},
			{Number: Int(2)},
		},
	}

	item := &ExportlistItem{}
	item.UpdateCollectedData(source)

	for name, pair := range map[string][2]*Timestamp{
		"watched_at":   {item.WatchedAt, source.WatchedAt},
		"collected_at": {item.CollectedAt, source.CollectedAt},
		"hidden_at":    {item.HiddenAt, source.HiddenAt},
	} {
		got, want := pair[0], pair[1]
		if got == nil || !got.Equal(want.Time) || got.Location() != time.UTC {
			t.Errorf("%s is %v, want %v in UTC", name, got, want)
		}
	}
	if *item.MediaType != "bluray" || *item.Resolution != "uhd_4k" || *item.Audio != "dts" || *item.AudioChannels != "5.1" || !*item.ThreeD {
		t.Errorf("metadata was not copied: %v", item)
	}

	if item.Seasons == nil || len(*item.Seasons) != 2 {
		t.Fatalf("seasons are %v, want 2", item.Seasons)
	}
	first, second := (*item.Seasons)[0], (*item.Seasons)[1]
	if *first.Number != 1 || *second.Number != 2 {
		t.Errorf("season numbers are %d and %d, want 1 and 2", *first.Number, *second.Number)
	}
	if second.Episodes != nil {
		t.Errorf("season without episodes got %v", second.Episodes)
	}
	if first.Episodes == nil || len(*first.Episodes) != 2 {
		t.Fatalf("episodes are %v, want 2", first.Episodes)
	}
	episode := (*first.Episodes)[0]
	if *episode.Number != 1 || *episode.MediaType != "bluray" || !*episode.ThreeD {
		t.Errorf("episode is %v, want number 1 with its metadata", episode)
	}
	if !episode.CollectedAt.Equal(testStamp(4).Time) || episode.CollectedAt.Location() != time.UTC {
		t.Errorf("episode collected_at is %v, want %v in UTC", episode.CollectedAt, testStamp(4))
	}
	if bare := (*first.Episodes)[1]; *bare.Number != 2 || bare.MediaType != nil {
		t.Errorf("episode without metadata is %v, want only number 2", bare)
	}
}

func TestExportlistItemUpdateCollectedDataLastCollected(t *testing.T) {
	source := &ExportlistItem{LastCollectedAt: testStamp(7)}

	item := &ExportlistItem{}
	item.UpdateCollectedData(source)

	if item.CollectedAt == nil || !item.CollectedAt.Equal(source.LastCollectedAt.Time) {
		t.Errorf("collected_at is %v, want last_collected_at %v", item.CollectedAt, source.LastCollectedAt)
	}
	if item.MediaType != nil || item.Seasons != nil {
		t.Errorf("item is %v, want no metadata and no seasons", item)
	}
}

func TestMediaUpdateCollectedData(t *testing.T) {
	source := &ExportlistItem{CollectedAt: testStamp(2), Metadata: testMetadata()}
	want := source.CollectedAt.Time

	movie := &Movie{}
	movie.UpdateCollectedData(source)
	show := &Show{}
	show.UpdateCollectedData(source)
	episode := &Episode{}
	episode.UpdateCollectedData(source)
	person := &Person{}
	person.UpdateCollectedData(source)
	season := &Season{}
	season.UpdateCollectedData(source)

	for name, got := range map[string]*Timestamp{
		"movie":   movie.CollectedAt,
		"show":    show.CollectedAt,
		"episode": episode.CollectedAt,
		"person":  person.CollectedAt,
	} {
		if got == nil || !got.Equal(want) || got.Location() != time.UTC {
			t.Errorf("%s collected_at is %v, want %v in UTC", name, got, want)
		}
	}
	for name, got := range map[string]*string{
		"movie":   movie.MediaType,
		"show":    show.MediaType,
		"episode": episode.MediaType,
		"season":  season.MediaType,
	} {
		if got == nil || *got != "bluray" {
			t.Errorf("%s media_type is %v, want bluray", name, got)
		}
	}
	if season.Resolution == nil || *season.Resolution != "uhd_4k" || season.ThreeD == nil || !*season.ThreeD {
		t.Errorf("season metadata is %v, want resolution and 3d copied", season)
	}

	plain := &ExportlistItem{CollectedAt: testStamp(2)}
	bare := &Movie{}
	bare.UpdateCollectedData(plain)
	if bare.MediaType != nil || bare.CollectedAt == nil {
		t.Errorf("movie without metadata is %v, want only collected_at", bare)
	}
}
