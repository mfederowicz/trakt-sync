// Package str used for structs
package str

// NewDeviceToken represents JSON when request for new token with code
type NewDeviceToken struct {
	Code         *string `json:"code"`
	ClientID     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
}
