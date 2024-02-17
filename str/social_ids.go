package str

type SocialIds struct {
	Twitter   *string `json:"twitter,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Wikipedia *string `json:"wikipedia,omitempty"`
}

func (s SocialIds) String() string {
	return Stringify(s)
}
