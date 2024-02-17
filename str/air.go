package str

type Air struct {
	Day      *string `json:"day,omitempty"`
	Time     *string `json:"time,omitempty"`
	TimeZone *string `json:"timezone,omitempty"`
}

func (a Air) String() string {
	return Stringify(a)
}

