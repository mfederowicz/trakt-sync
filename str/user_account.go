package str

// UserAccount represents JSON user account object
type UserAccount struct {
	Timezone   *string `json:"timezone,omitempty"`
	DateFormat *string `json:"date_format,omitempty"`
	Time24HR   *bool   `json:"time_24hr,omitempty"`
	CoverImage *string `json:"cover_image,omitempty"`
	ShareCode  *string `json:"share_code,omitempty"`
	DisplayAds *bool   `json:"display_ads,omitempty"`
}

func (u UserAccount) String() string {
	return Stringify(u)
}
