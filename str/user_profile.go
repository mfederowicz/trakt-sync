package str

// UserProfile represents JSON user profile object
type UserProfile struct {
	Email           *string    `json:"email,omitempty"`
	Username        *string    `json:"username,omitempty"`
	Private         *bool      `json:"private,omitempty"`
	Deleted         *bool      `json:"deleted,omitempty"`
	Name            *string    `json:"name,omitempty"`
	Vip             *bool      `json:"vip,omitempty"`
	VipEp           *bool      `json:"vip_ep,omitempty"`
	Director        *bool      `json:"director,omitempty"`
	IDs             *IDs       `json:"ids,omitempty"`
	JoinedAt        *Timestamp `json:"joined_at,omitempty"`
	HiddenAt        *Timestamp `json:"hidden_at,omitempty"`
	Location        *string    `json:"location,omitempty"`
	About           *string    `json:"about,omitempty"`
	Gender          *string    `json:"gender,omitempty"`
	Age             *int       `json:"age,omitempty"`
	Images          *Images    `json:"images,omitempty"`
	VipOg           *bool      `json:"vip_og,omitempty"`
	VipYears        *int       `json:"vip_years,omitempty"`
	VipCoverImage   *string    `json:"vip_cover_image,omitempty"`
	VipVeteranSince *Timestamp `json:"vip_veteran_since,omitempty"`
	VipVeteranYears *int       `json:"vip_veteran_years,omitempty"`
	VipVeteranTier  *int       `json:"vip_veteran_tier,omitempty"`
	VipVeteranTitle *string    `json:"vip_veteran_title,omitempty"`
	VipGraceEndsAt  *Timestamp `json:"vip_grace_ends_at,omitempty"`
}

func (u UserProfile) String() string {
	return Stringify(u)
}
