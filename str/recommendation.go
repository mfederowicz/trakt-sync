// Package str used for structs
package str

// Recommendation represents JSON recommendation object
type Recommendation struct {
	Title         *string      `json:"title,omitempty"`
	Year          *int         `json:"year,omitempty"`
	IDs           *IDs         `json:"ids,omitempty"`
	FavoritedBy   *[]UserNotes `json:"favorited_by,omitempty"`
	RecommendedBy *[]UserNotes `json:"recommended_by,omitempty"`
}

func (r Recommendation) String() string {
	return Stringify(r)
}
