// Package str used for structs
package str

// MovieReport represents JSON movie report object
type MovieReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (m MovieReport) String() string {
	return Stringify(m)
}
