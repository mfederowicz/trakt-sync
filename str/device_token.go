// Package str used for structs
package str

// DeviceToken represents JSON response for /device/token
type DeviceToken struct {
	AccessToken  *string `json:"access_token"`
	TokenType    *string `json:"token_type"`
	ExpiresIn    *int    `json:"expires_in"`
	RefreshToken *string `json:"refresh_token"`
	Scope        *string `json:"scope"`
	CreatedAt    *int    `json:"created_at"`
}
