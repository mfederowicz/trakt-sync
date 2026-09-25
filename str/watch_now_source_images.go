// Package str used for structs
package str

// WatchNowSourceImages represents JSON watch now source images object
type WatchNowSourceImages struct {
	Logo    *string `json:"logo,omitempty"`
	Channel *string `json:"channel,omitempty"`
}

func (w WatchNowSourceImages) String() string {
	return Stringify(w)
}
