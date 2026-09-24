// Package str used for structs
package str

// ShowProgress represents JSON up next or watched progress item: a show with its progress
type ShowProgress struct {
	Show     *Show            `json:"show,omitempty"`
	Progress *WatchedProgress `json:"progress,omitempty"`
}

func (s ShowProgress) String() string {
	return Stringify(s)
}
