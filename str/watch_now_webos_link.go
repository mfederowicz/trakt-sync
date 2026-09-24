// Package str used for structs
package str

// WatchNowWebosLink represents JSON webOS link object
type WatchNowWebosLink struct {
	ID     *string              `json:"id,omitempty"`
	Params *WatchNowWebosParams `json:"params,omitempty"`
}

func (w WatchNowWebosLink) String() string {
	return Stringify(w)
}
