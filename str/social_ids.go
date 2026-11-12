package str

// SocialIDs represents JSON object with social media handlers
type SocialIDs struct {
	Twitter   *string `json:"twitter,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Wikipedia *string `json:"wikipedia,omitempty"`
	// TikTok, YouTube, Wikidata and WikipediaLanguage are sent for people only.
	TikTok            *string `json:"tiktok,omitempty"`
	YouTube           *string `json:"youtube,omitempty"`
	Wikidata          *string `json:"wikidata,omitempty"`
	WikipediaLanguage *string `json:"wikipedia_language,omitempty"`
}

func (s SocialIDs) String() string {
	return Stringify(s)
}
