// Package str used for structs
package str

// CurrentDeviceToken represents JSON current device token object
type CurrentDeviceToken struct {
	RefreshToken *string `json:"refresh_token"`
	ClientID     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
	RedirectURI  *string `json:"redirect_uri"`
	GrantType    *string `json:"grant_type"`
}
