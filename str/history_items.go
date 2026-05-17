// Package str used for structs
package str

// HistoryItems represents JSON list object
type HistoryItems struct {
	Movies   *[]Movie       `json:"movies,omitempty"`
	Shows    *[]Show        `json:"shows,omitempty"`
	Seasons  *[]Season      `json:"seasons,omitempty"`
	Episodes *[]Episode     `json:"episodes,omitempty"`
	Users    *[]UserProfile `json:"users,omitempty"`
}

func (h HistoryItems) String() string {
	return Stringify(h)
}
