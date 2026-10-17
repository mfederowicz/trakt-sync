package str

// AvailableOn represents JSON source of a collected item, sent with extended available_on; an empty list is left out of exports
type AvailableOn struct {
	Name *string `json:"name,omitempty"`
}

func (a AvailableOn) String() string {
	return Stringify(a)
}
