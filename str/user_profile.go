package str

type UserProfile struct {
	Userame       *string    `json:"username,omitempty"`
	Private       *bool      `json:"private,omitempty"`
	Name          *string    `json:"name,omitempty"`
	Vip           *bool      `json:"vip,omitempty"`
	VipEp         *bool      `json:"vip_ep,omitempty"`
	Ids           *Ids       `json:"ids,omitempty"`
	JoinedAt      *Timestamp `json:"joined_at,omitempty"`
	Location      *string    `json:"location,omitempty"`
	About         *string    `json:"about,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	Age           *int       `json:"age,omitempty"`
	Images        *Images    `json:"images,omitempty"`
	VipOg         *bool      `json:"vip_og,omitempty"`
	VipYears      *int       `json:"vip_years,omitempty"`
	VipCoverImage *string    `json:"vip_cover_image,omitempty"`
}

func (u UserProfile) String() string {
	return Stringify(u)
}
