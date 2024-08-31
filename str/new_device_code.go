// Package str used for structs
package str

// NewDeviceCode represents request new device code payload.
type NewDeviceCode struct {
	ClientID *string `json:"client_id"`
}

func (d NewDeviceCode) String() string {
	return Stringify(d)
}
