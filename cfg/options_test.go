package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidConfigTypeSlice(t *testing.T) {
	t.Helper()
	got := IsValidConfigTypeSlice([]string{"movie", "show"}, []string{"xxx"})
	if bool(got) {
		t.Fatalf("Expected %v, got %v", false, bool(got))
	}
}

func TestModuleConfigTypeComments(t *testing.T) {
	t.Helper()

	got := ModuleActionConfig["comments:trending"].Type

	assert.Equal(t, got, []string{"all", "movies", "shows", "seasons", "episodes", "lists"})

	got = ModuleActionConfig["comments:recent"].Type

	assert.Equal(t, got, []string{"all", "movies", "shows", "seasons", "episodes", "lists"})
}

func TestModuleConfigTypeUsers(t *testing.T) {
	t.Helper()

	got := ModuleActionConfig["users:watched"].Type

	assert.Equal(t, got, []string{"movies", "shows"})
}
