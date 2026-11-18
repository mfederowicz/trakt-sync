package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A season or episode missing under a show comes back in not_found.seasons or not_found.episodes as that show,
// holding only the missing season or episode (trakt/trakt-api 496cfc4); the show and what is missing must survive decoding.
const showShapedNotFound = `{
	"seasons": [{"title": "Breaking Bad", "year": 2008, "ids": {"trakt": 1388}, "seasons": [{"number": 99}]}],
	"episodes": [{"title": "Breaking Bad", "year": 2008, "ids": {"trakt": 1388}, "seasons": [{"number": 1, "episodes": [{"number": 99, "rating": 8}]}]}]
}`

func wantShowShapedSeasons() *[]Show {
	return &[]Show{{Title: String("Breaking Bad"), Year: Int(2008), IDs: &IDs{Trakt: Int64(1388)}, Seasons: &[]Season{{Number: Int(99)}}}}
}

func wantShowShapedEpisodes() *[]Show {
	rating := float32(8)
	return &[]Show{{Title: String("Breaking Bad"), Year: Int(2008), IDs: &IDs{Trakt: Int64(1388)},
		Seasons: &[]Season{{Number: Int(1), Episodes: &[]Episode{{Number: Int(99), Rating: &rating}}}}}}
}

func TestResultNotFoundShowShaped(t *testing.T) {
	got := ResultNotFound{}
	assert.NoError(t, json.Unmarshal([]byte(showShapedNotFound), &got))
	assert.Equal(t, wantShowShapedSeasons(), got.Seasons)
	assert.Equal(t, wantShowShapedEpisodes(), got.Episodes)
}

func TestCollectionResultNotFoundShowShaped(t *testing.T) {
	got := CollectionResultNotFound{}
	assert.NoError(t, json.Unmarshal([]byte(showShapedNotFound), &got))
	assert.Equal(t, wantShowShapedSeasons(), got.Seasons)
	assert.Equal(t, wantShowShapedEpisodes(), got.Episodes)
}

// A season or episode sent by its own ids is echoed with those ids only.
func TestResultNotFoundByIDs(t *testing.T) {
	got := ResultNotFound{}
	assert.NoError(t, json.Unmarshal([]byte(`{"seasons":[{"ids":{"trakt":5}}],"episodes":[{"ids":{"trakt":6}}]}`), &got))
	assert.Equal(t, &[]Show{{IDs: &IDs{Trakt: Int64(5)}}}, got.Seasons)
	assert.Equal(t, &[]Show{{IDs: &IDs{Trakt: Int64(6)}}}, got.Episodes)
}
