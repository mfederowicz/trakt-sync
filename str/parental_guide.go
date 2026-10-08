package str

// ParentalGuide represents the parental guide filters: severity ranges min-max, from 0 (none) to 3 (severe)
type ParentalGuide struct {
	Alcohol        string
	Frightening    string
	IncludeUnrated bool
	Nudity         string
	Profanity      string
	Violence       string
}

func (p ParentalGuide) String() string {
	return Stringify(p)
}
