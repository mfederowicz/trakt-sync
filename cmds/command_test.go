// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

func TestValidModuleActionTypeNotFoundModule(t *testing.T) {
	t.Helper()
	var options = &str.Options{Module: "xyz"}
	var testCmd = &Command{
		Name: "test",
	}

	got := testCmd.ValidModuleActionType(options)
	assert.Contains(t, got.Error(), "not found config for module 'xyz'")
}

func TestValidModuleActionTypeNotFoundType(t *testing.T) {
	t.Helper()
	var options = &str.Options{Module: "users", Action: "watched", Type: "abc"}
	var testCmd = &Command{
		Name: "test",
	}

	got := testCmd.ValidModuleActionType(options)
	assert.Contains(t, got.Error(), "type 'abc' is not valid for module 'users'")
}

func TestValidModuleActionTypeOk(t *testing.T) {
	t.Helper()
	var options = &str.Options{Module: "users", Action: "watched", Type: "shows"}
	var testCmd = &Command{
		Name: "test",
	}

	got := testCmd.ValidModuleActionType(options)
	assert.Equal(t, got, nil)
}
