// Package str used for structs
package str

// Crew represents JSON crew object
type Crew struct {
	Writing    *[]Job `json:"writing,omitempty"`
	Directing  *[]Job `json:"directing,omitempty"`
	Production *[]Job `json:"production,omitempty"`
}

func (c Crew) String() string {
	return Stringify(c)
}
