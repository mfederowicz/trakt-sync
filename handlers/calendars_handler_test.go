// Package handlers used to handle module actions
package handlers

import "testing"

func TestCalendarTarget(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{action: "my-media", want: "my"},
		{action: "all-media", want: "all"},
		{action: "my-streaming", want: "my"},
		{action: "all-streaming", want: "all"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			if got := calendarTarget(tt.action); got != tt.want {
				t.Errorf("calendarTarget(%q) = %q, want %q", tt.action, got, tt.want)
			}
		})
	}
}
