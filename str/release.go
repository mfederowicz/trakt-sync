// Package str used for structs
package str

// Release represents JSON release object
type Release struct {
	Country       *string `json:"country,omitempty"`
	Certification *string `json:"certification,omitempty"`
	ReleaseDate   *Timestamp `json:"release_date,omitempty"`
	ReleaseType   *string `json:"release_type,omitempty"`
	Note          *string `json:"note,omitempty"`
}

func (r Release) String() string {
	return Stringify(r)
}
