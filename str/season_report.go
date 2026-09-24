// Package str used for structs
package str

// SeasonReport represents JSON season report object
type SeasonReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (s SeasonReport) String() string {
	return Stringify(s)
}
