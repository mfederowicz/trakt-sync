// Package str used for structs
package str

import (
	"errors"
)

// Errors represents JSON response for errors object
type Errors struct {
	Comment *[]string `json:"comment,omitempty"`
}

// GetComments returns a joined string of errors for the Comment field or nil if empty
func (e *Errors) GetComments() error {
	return joinErrors(e.Comment)
}

// consts
const (
	Empty = 0
)

// Helper function to join errors or return nil if empty
// Helper function to convert []string to []error and use errors.Join()
func joinErrors(errs *[]string) error {
	if errs == nil || len(*errs) == Empty {
		return nil
	}

	// Convert []string to []error
	errList := make([]error, len(*errs))
	for i, errMsg := range *errs {
		errList[i] = errors.New(errMsg)
	}

	// Use errors.Join() to combine them
	return errors.Join(errList...)
}

func (e Errors) String() string {
	return Stringify(e)
}
