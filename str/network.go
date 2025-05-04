// Package str used for structs
package str

// Network represents JSON network object
type Network struct {
	Friends   *int `json:"friends,omitempty"`
	Followers *int `json:"followers,omitempty"`
	Following *int `json:"following,omitempty"`
}

func (n Network) String() string {
	return Stringify(n)
}
