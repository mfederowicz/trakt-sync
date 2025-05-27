// Package str used for structs
package str

// DeviceToken represents JSON response for /device/token
type DeviceToken struct {
	AccessToken  *string `json:"access_token"`
	TokenType    *string `json:"token_type"`
	ExpiresIn    *int64  `json:"expires_in"`
	RefreshToken *string `json:"refresh_token"`
	Scope        *string `json:"scope"`
	CreatedAt    *int64  `json:"created_at"`
}

// ToToken convert device_token to token
func (d *DeviceToken) ToToken() *Token {
	return &Token{
		AccessToken:  *d.AccessToken,
		TokenType:    *d.TokenType,
		RefreshToken: *d.RefreshToken,
		Scope:        *d.Scope,
		ExpiresIn:    *d.ExpiresIn,
		CreatedAt:    *d.CreatedAt,
	}
}
