// Package str used for structs
package str

// PersonReport represents JSON person report object
type PersonReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (p PersonReport) String() string {
	return Stringify(p)
}
