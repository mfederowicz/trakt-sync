// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
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

func TestValidSectionHiddenItems(t *testing.T) {
	tests := []struct {
		action  string
		section string
		wantErr bool
	}{
		{action: "hidden_items", section: cfg.DefaultConfig().UsersSection},
		{action: "hidden_items", section: "dropped"},
		{action: "hidden_items", section: "progress_watched"},
		{action: "hidden_items", section: "progress_watched_reset", wantErr: true},
		{action: "remove_hidden_items", section: cfg.DefaultConfig().UsersSection},
		{action: "remove_hidden_items", section: "progress_watched"},
		{action: "remove_hidden_items", section: "abc", wantErr: true},
		{action: "add_hidden_items", section: "dropped"},
		{action: "add_hidden_items", section: "abc", wantErr: true},
	}

	testCmd := &Command{Name: "test"}
	for _, tt := range tests {
		t.Run(tt.action+"/"+tt.section, func(t *testing.T) {
			got := testCmd.ValidSection(&str.Options{Module: "users", Action: tt.action, Section: tt.section})
			if tt.wantErr {
				assert.ErrorContains(t, got, "section '"+tt.section+"' is not valid")
				return
			}
			assert.NoError(t, got)
		})
	}
}

func TestUsersWatchedType(t *testing.T) {
	tests := []struct {
		name     string
		flagType string
		wantType string
		wantErr  bool
	}{
		{name: "no -t uses movies", flagType: cfg.DefaultConfig().UsersType, wantType: "movies"},
		{name: "shows", flagType: "shows", wantType: "shows"},
		{name: "episodes is not a watched type", flagType: "episodes", wantType: "episodes", wantErr: true},
	}

	testCmd := &Command{Name: "users"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := UpdateOptionsWithCommandUsersFlags(testCmd, &str.Options{Module: "users", Action: "watched", UserName: "sean", Type: tt.flagType})
			assert.Equal(t, tt.wantType, options.Type)

			err := testCmd.ValidModuleActionType(options)
			if tt.wantErr {
				assert.ErrorContains(t, err, "type 'episodes' is not valid")
				return
			}
			assert.NoError(t, err)
		})
	}
}

// TestAvflagsHasEveryCommand checks that ValidFlags accepts every module name;
// a module missing from Avflags fails with "invalid flags" before its action runs.
func TestAvflagsHasEveryCommand(t *testing.T) {
	for _, cmd := range Commands {
		if !Avflags[cmd.Name] {
			t.Errorf("command %q is missing from Avflags", cmd.Name)
		}
	}
}
