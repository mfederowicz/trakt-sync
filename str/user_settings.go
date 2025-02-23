// Package str used for structs
package str

// UserSettings represents JSON user stats object
type UserSettings struct {
	User        *UserProfile `json:"user,omitempty"`
	Account     *UserAccount `json:"account,omitempty"`
	Connections *Connections `json:"connections,omitempty"`
	SharingText *SharingText `json:"sharing_text,omitempty"`
	Limits      *Limits      `json:"limits,omitempty"`
}

func (u UserSettings) String() string {
	return Stringify(u)
}
