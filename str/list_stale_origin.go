package str

// ListStaleOrigin represents JSON tool that kept a personal list updated and stopped syncing it
type ListStaleOrigin struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

func (l ListStaleOrigin) String() string {
	return Stringify(l)
}
