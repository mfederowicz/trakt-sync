// Package str used for structs
package str

// UserReport represents JSON user report object
type UserReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (u UserReport) String() string {
	return Stringify(u)
}
