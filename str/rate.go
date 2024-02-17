package str

// Rate represents the rate limit for the current client.
type Rate struct {
	// The time at which the current rate limit will reset.
	Reset Timestamp `json:"reset"`
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`
	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`
}

func (r Rate) String() string {
	return Stringify(r)
}
