// Package str used for structs
package str

// Connections represents JSON connections object
type Connections struct {
	Facebook  *bool `json:"facebook,omitempty"`
	Twitter   *bool `json:"twitter,omitempty"`
	Mastodon  *bool `json:"mastodon,omitempty"`
	Google    *bool `json:"google,omitempty"`
	Tumblr    *bool `json:"tumblr,omitempty"`
	Medium    *bool `json:"medium,omitempty"`
	Slack     *bool `json:"slack,omitempty"`
	Apple     *bool `json:"apple,omitempty"`
	Dropbox   *bool `json:"dropbox,omitempty"`
	Microsoft *bool `json:"microsoft,omitempty"`
}

func (c Connections) String() string {
	return Stringify(c)
}
