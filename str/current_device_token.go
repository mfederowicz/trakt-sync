package str

type CurrentDeviceToken struct {
	RefreshToken *string `json:"refresh_token"`
	ClientId     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
	RedirectUri  *string `json:"redirect_uri"`
	GrantType    *string `json:"grant_type"`
}

