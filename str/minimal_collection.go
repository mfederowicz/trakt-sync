// Package str used for structs
package str

// MinimalCollection represents JSON minimal movie or episode collection: Trakt ID -> collected_at
type MinimalCollection map[string]*Timestamp

func (m MinimalCollection) String() string {
	return Stringify(m)
}
