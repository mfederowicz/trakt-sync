// Package str used for structs
package str

// Avatar represents JSON avatar object
type Avatar struct {
	Full *string `json:"full,omitempty"`
}

// Images represents JSON images object
type Images struct {
	Avatar *Avatar `json:"avatar,omitempty"`
}
