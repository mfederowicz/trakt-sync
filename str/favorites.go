// Package str used for structs
package str

// Favorites represents JSON favorites object
type Favorites struct {
	ItemCount *int       `json:"item_count,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

func (f Favorites) String() string {
	return Stringify(f)
}
