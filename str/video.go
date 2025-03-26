// Package str used for structs
package str

// Video represents JSON video object
type Video struct {
	Title       *string    `json:"title,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Site        *string    `json:"site,omitempty"`
	Type        *string    `json:"type,omitempty"`
	Size        *int       `json:"size,omitempty"`
	Official    *bool      `json:"official,omitempty"`
	PublishedAt *Timestamp `json:"published_at,omitempty"`
	Country     *string    `json:"country,omitempty"`
	Language    *string    `json:"language,omitempty"`
}

func (v Video) String() string {
	return Stringify(v)
}
