// Package str used for structs
package str

// Sharing represents JSON sharing object
type Sharing struct {
	Twitter  *bool `json:"twitter,omitempty"`
	Mastodon *bool `json:"mastodon,omitempty"`
	Tumblr   *bool `json:"tumblr,omitempty"`
}

func (s Sharing) String() string {
	return Stringify(s)
}
