package str

// Code represents request new device code payload.
type NewDeviceCode struct {
	ClientId *string `json:"client_id"`
}

func (d NewDeviceCode) String() string {
	return Stringify(d)
}
