package str

// ReviewPersonaHighlight represents JSON year in review persona highlight object
type ReviewPersonaHighlight struct {
	Kind  *string  `json:"kind,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (r ReviewPersonaHighlight) String() string {
	return Stringify(r)
}
