package str

type NewDeviceToken struct {
	Code         *string `json:"code"`
	ClientId     *string `json:"client_id"`
	ClientSecret *string `json:"client_secret"`
}
