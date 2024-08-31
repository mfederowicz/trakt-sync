// Package str used for structs
package str

// SocialIDs represents JSON object with social media handlers
type SocialIDs struct {
	Twitter   *string `json:"twitter,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Wikipedia *string `json:"wikipedia,omitempty"`
}

func (s SocialIDs) String() string {
	return Stringify(s)
}
