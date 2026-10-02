package str

// ResultNotFound represents JSON not found object
type ResultNotFound struct {
	Movies   *[]ExportlistItem `json:"movies,omitempty"`
	Shows    *[]Show           `json:"shows,omitempty"`
	Seasons  *[]Season         `json:"seasons,omitempty"`
	Episodes *[]Episodes       `json:"episodes,omitempty"`
	People   *[]Person         `json:"people,omitempty"`
	Users    *[]UserProfile    `json:"users,omitempty"`
}

func (c ResultNotFound) String() string {
	return Stringify(c)
}
