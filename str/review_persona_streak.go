package str

// ReviewPersonaStreak represents JSON year in review persona longest streak object
type ReviewPersonaStreak struct {
	Longest   *int       `json:"longest,omitempty"`
	StartedAt *Timestamp `json:"started_at,omitempty"`
}

func (r ReviewPersonaStreak) String() string {
	return Stringify(r)
}
