// Package str used for structs
package str

// SocialActivity represents JSON social activity object: a movie or an episode someone in the user's social graph watched
type SocialActivity struct {
	ID         *int64       `json:"id,omitempty"`
	ActivityAt *Timestamp   `json:"activity_at,omitempty"`
	Action     *string      `json:"action,omitempty"`
	User       *UserProfile `json:"user,omitempty"`
	UserRating *int         `json:"user_rating,omitempty"`
	Type       *string      `json:"type,omitempty"`
	Movie      *Movie       `json:"movie,omitempty"`
	Show       *Show        `json:"show,omitempty"`
	Episode    *Episode     `json:"episode,omitempty"`
}

func (s SocialActivity) String() string {
	return Stringify(s)
}
