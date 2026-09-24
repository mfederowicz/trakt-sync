// Package str used for structs
package str

// RecentSearch represents JSON body for adding or removing a recent search
type RecentSearch struct {
	Query string `json:"query"`
	ID    int64  `json:"id"`
	Type  string `json:"type"`
}

func (r RecentSearch) String() string {
	return Stringify(r)
}
