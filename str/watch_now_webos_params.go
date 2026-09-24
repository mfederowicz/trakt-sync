// Package str used for structs
package str

// WatchNowWebosParams represents JSON webOS link params object
type WatchNowWebosParams struct {
	ContentTarget *string `json:"contentTarget,omitempty"`
}

func (w WatchNowWebosParams) String() string {
	return Stringify(w)
}
