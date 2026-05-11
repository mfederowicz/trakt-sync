// Package str used for structs
package str

import (
	"strconv"
	"strings"
)

// SliceInt as a slice of ints
type SliceInt []int

// Get string value of slice
func (i *SliceInt) String() string {
	const (
		emptySliceLen = 0
	)

	if len(*i) == emptySliceLen {
		return ""
	}

	uniqueValues := map[int]struct{}{}
	var deduplicated []string
	for _, item := range *i {
		if _, seen := uniqueValues[item]; !seen {
			uniqueValues[item] = struct{}{}
			deduplicated = append(deduplicated, strconv.Itoa(item))
		}
	}

	return strings.Join(deduplicated, ",")
}

// Set add element to slice
func (i *SliceInt) Set(value int) error {
	*i = append(*i, value)
	return nil
}
