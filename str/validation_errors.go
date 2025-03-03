// Package str used for structs
package str

// ValidationErrors represents errors object
type ValidationErrors struct {
	Errors *Errors `json:"errors,omitempty"` // errors object

}

func (v ValidationErrors) String() string {
	return Stringify(v)
}
