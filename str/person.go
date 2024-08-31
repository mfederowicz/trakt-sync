// Package str used for structs
package str

// Person represents JSON person object
type Person struct {
	Name               *string    `json:"name,omitempty"`
	IDs                *IDs       `json:"ids,omitempty"`
	SocialIDs          *SocialIDs `json:"social_ids,omitempty"`
	Biography          *string    `json:"biography,omitempty"`
	Birthday           *string    `json:"birthday,omitempty"`
	Death              *string    `json:"death,omitempty"`
	Birthplace         *string    `json:"birthplace,omitempty"`
	Homepage           *string    `json:"homepage,omitempty"`
	Gender             *string    `json:"gender,omitempty"`
	KnownForDepartment *string    `json:"known_for_department,omitempty"`
	UpdatedAt          *Timestamp `json:"updated_at,omitempty"`
}

func (p Person) String() string {
	return Stringify(p)
}
