package str

import "strings"

// Define a type named "StrSlice" as a slice of strings
type StrSlice []string

// Get string value of slice
func (i *StrSlice) String() string {
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

// The second method is Set(value string) error
func (i *StrSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
