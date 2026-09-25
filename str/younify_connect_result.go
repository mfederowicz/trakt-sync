// Package str used for structs
package str

// YounifyConnectResult represents JSON streaming connection response with the web auth URL
type YounifyConnectResult struct {
	URL *string `json:"url,omitempty"`
}

func (y YounifyConnectResult) String() string {
	return Stringify(y)
}
