// Package str used for structs
package str

// Certifications represents JSON certifications object
type Certifications struct {
	Us []*Certification `json:"us,omitempty"`
}

func (c Certifications) String() string {
	return Stringify(c)
}
