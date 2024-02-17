package str

type Season struct {
	Number *int `json:"number,omitempty"`
	Ids    *Ids `json:"ids,omitempty"`
}

func (s Season) String() string {
	return Stringify(s)
}
