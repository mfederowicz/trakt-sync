// Package str used for structs
package str

// ListReport represents JSON list report object
type ListReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (l ListReport) String() string {
	return Stringify(l)
}
