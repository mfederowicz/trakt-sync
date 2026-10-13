package str

// ReviewPersona represents JSON year in review persona object; scores holds a 0-100 score for every persona
type ReviewPersona struct {
	Persona            *string                   `json:"persona,omitempty"`
	RunnerUp           *string                   `json:"runner_up,omitempty"`
	Confidence         *string                   `json:"confidence,omitempty"`
	Rarity             *int                      `json:"rarity,omitempty"`
	CardNumber         *int                      `json:"card_number,omitempty"`
	Traits             []string                  `json:"traits,omitempty"`
	Highlights         []*ReviewPersonaHighlight `json:"highlights,omitempty"`
	RunnerUpHighlights []*ReviewPersonaHighlight `json:"runner_up_highlights,omitempty"`
	Scores             map[string]int            `json:"scores,omitempty"`
	Streak             *ReviewPersonaStreak      `json:"streak,omitempty"`
	Monthly            []*ReviewPersonaMonth     `json:"monthly,omitempty"`
}

func (r ReviewPersona) String() string {
	return Stringify(r)
}
