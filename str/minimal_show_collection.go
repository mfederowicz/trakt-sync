// Package str used for structs
package str

// MinimalShowCollection represents JSON minimal show collection: show Trakt ID -> season -> episode -> collected_at
type MinimalShowCollection map[string]map[string]MinimalCollection

func (m MinimalShowCollection) String() string {
	return Stringify(m)
}
