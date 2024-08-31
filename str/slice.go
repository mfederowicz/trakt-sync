// Package str used for structs
package str

import "strings"

// Slice as a slice of strings
type Slice []string

// Get string value of slice
func (i *Slice) String() string {
	if len(*i) == 0 {
		return ""
	}

	uniqueValues := make(map[string]struct{})
	var deduplicated []string
	for _, item := range *i {
		if _, seen := uniqueValues[item]; !seen {
			uniqueValues[item] = struct{}{}
			deduplicated = append(deduplicated, item)
		}
	}

	return strings.Join(deduplicated, ",")

}

// Set add element to slice 
func (i *Slice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
