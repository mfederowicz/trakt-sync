package str

import "encoding/json"

// UserNotes represents JSON user notes object
type UserNotes struct {
	User  *UserProfile `json:"user,omitempty"`
	Notes *string      `json:"notes,omitempty"`
}

// userNotesFields is UserNotes without its methods, so unmarshalling it does not recurse.
type userNotesFields UserNotes

// UnmarshalJSON decodes both shapes of an entry: the nested one (`user` object plus `notes`) and the flat one
// of the contract (profile fields plus `notes`), whose profile fields go to User.
func (u *UserNotes) UnmarshalJSON(data []byte) error {
	fields := userNotesFields{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	*u = UserNotes(fields)
	if u.User != nil {
		return nil
	}
	profile := UserProfile{}
	if err := json.Unmarshal(data, &profile); err != nil {
		return err
	}
	if profile != (UserProfile{}) {
		u.User = &profile
	}
	return nil
}

func (u UserNotes) String() string {
	return Stringify(u)
}
