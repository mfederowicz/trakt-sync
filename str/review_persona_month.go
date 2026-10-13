package str

// ReviewPersonaMonth represents JSON year in review persona of one month
type ReviewPersonaMonth struct {
	Month   *int    `json:"month,omitempty"`
	Persona *string `json:"persona,omitempty"`
}

func (r ReviewPersonaMonth) String() string {
	return Stringify(r)
}
