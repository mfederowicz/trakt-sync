// Package str used for structs
package str

// SettingsUpdateUser represents JSON profile values of a settings update
type SettingsUpdateUser struct {
	Name     *string `json:"name,omitempty"`
	About    *string `json:"about,omitempty"`
	Location *string `json:"location,omitempty"`
	Private  *bool   `json:"private,omitempty"`
	Dob      *string `json:"dob,omitempty"`
}

func (s SettingsUpdateUser) String() string {
	return Stringify(s)
}
