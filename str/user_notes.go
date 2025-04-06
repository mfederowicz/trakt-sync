// Package str used for structs
package str

// UserNotes represents JSON user notes object
type UserNotes struct {
	User  *UserProfile `json:"user,omitempty"`
	Notes *string      `json:"notes,omitempty"`
}

func (u UserNotes) String() string {
	return Stringify(u)
}
