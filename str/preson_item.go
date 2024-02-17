package str

type PersonItem struct {
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
	Person    *Person        `json:"person,omitempty"`
}

func (p PersonItem) String() string {
	return Stringify(p)
}
