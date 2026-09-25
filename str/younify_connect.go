// Package str used for structs
package str

// YounifyConnect represents JSON streaming connection request object
type YounifyConnect struct {
	ServiceID *string `json:"service_id,omitempty"`
	ReturnURL *string `json:"return_url,omitempty"`
}

func (y YounifyConnect) String() string {
	return Stringify(y)
}
