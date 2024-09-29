package cfg

import (
	"testing"
)

func TestIsValidConfigTypeSlice(t *testing.T) {
	t.Helper()

	got := IsValidConfigTypeSlice([]string{"movie","show"}, []string{"xxx"})
	if bool(got) != false {
		t.Fatalf("Expected %v, got %v", false, bool(got))
	}
}
