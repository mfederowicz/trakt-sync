// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestArgsToMap checks the flag names read from the arguments after the global flags:
// -flag=value gives the same name as -flag value, and an empty value does not panic.
func TestArgsToMap(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want map[string]bool
	}{
		{name: "flag and value", args: []string{"movies", "-a", "trending"}, want: map[string]bool{"movies": true, "a": true}},
		{name: "flag=value", args: []string{"movies", "-a=trending"}, want: map[string]bool{"movies": true, "a": true}},
		{name: "double dash flag=value", args: []string{"search", "--field=title"}, want: map[string]bool{"search": true, "field": true}},
		{name: "value with equals sign", args: []string{"search", "-q=a=b"}, want: map[string]bool{"search": true, "q": true}},
		{name: "separate value with equals sign", args: []string{"search", "-q", "a=b"}, want: map[string]bool{"search": true, "q": true}},
		{name: "argument after flag=value is not its value", args: []string{"movies", "-a=trending", "foo"}, want: map[string]bool{"movies": true, "a": true, "foo": true}},
		{name: "flag=value then flag", args: []string{"shows", "-status=ended", "-a", "popular"}, want: map[string]bool{"shows": true, "status": true, "a": true}},
		{name: "empty value", args: []string{"movies", "-i", ""}, want: map[string]bool{"movies": true, "i": true}},
		{name: "bool flag last", args: []string{"sync", "-hide_completed"}, want: map[string]bool{"sync": true, "hide_completed": true}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, argsToMap(tt.args)); diff != "" {
				t.Errorf("argsToMap(%v) mismatch (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}

// TestValidArgs checks that a known flag is accepted in both forms and an unknown name is not.
func TestValidArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "flag and value", args: []string{"shows", "-a", "popular", "-status", "ended"}, want: true},
		{name: "flag=value", args: []string{"shows", "-a=popular", "-status=ended"}, want: true},
		{name: "empty value", args: []string{"movies", "-a", "summary", "-i", ""}, want: true},
		{name: "unknown flag=value", args: []string{"shows", "-no_such_flag=1"}, want: false},
		{name: "unknown argument after flag=value", args: []string{"shows", "-a=popular", "no_such_argument"}, want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := validArgs(tt.args); got != tt.want {
				t.Errorf("validArgs(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
