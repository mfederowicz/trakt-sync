// Package handlers used to handle module actions
package handlers

import "testing"

// TestMain switches off the pause between API calls once for the whole package, so tests can follow pages.
func TestMain(m *testing.M) {
	pageDelay = 0
	m.Run()
}
