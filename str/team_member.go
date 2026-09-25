// Package str used for structs
package str

// TeamMember represents JSON Trakt team member object
type TeamMember struct {
	User *UserProfile `json:"user,omitempty"`
}

func (t TeamMember) String() string {
	return Stringify(t)
}
