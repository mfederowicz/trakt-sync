// Package str used for structs
package str

// Air represents JSON response for media object
type Air struct {
	Day      *string `json:"day,omitempty"`
	Time     *string `json:"time,omitempty"`
	TimeZone *string `json:"timezone,omitempty"`
}

func (a Air) String() string {
	return Stringify(a)
}
