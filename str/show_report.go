// Package str used for structs
package str

// ShowReport represents JSON show report object
type ShowReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (s ShowReport) String() string {
	return Stringify(s)
}
