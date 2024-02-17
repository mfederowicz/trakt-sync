package str

type Crew struct {
	Writing    *[]Job `json:"writing,omitempty"`
	Directing  *[]Job `json:"directing,omitempty"`
	Production *[]Job `json:"production,omitempty"`
}

func (c Crew) String() string {
	return Stringify(c)
}
