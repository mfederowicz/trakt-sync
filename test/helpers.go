// Package test used for process tests
package test

// MapsStringBoolEqual  check if two maps are equal
func MapsStringBoolEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, ok := b[key]; !ok || bValue != value {
			return false
		}
	}
	return true
}
