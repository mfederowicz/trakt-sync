// Package handlers used to handle module actions
package handlers

import "testing"

func TestCalendarTarget(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{action: "my_media", want: "my"},
		{action: "all_media", want: "all"},
		{action: "my_streaming", want: "my"},
		{action: "all_streaming", want: "all"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			if got := calendarTarget(tt.action); got != tt.want {
				t.Errorf("calendarTarget(%q) = %q, want %q", tt.action, got, tt.want)
			}
		})
	}
}
