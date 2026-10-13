// Package deprecation keeps the history of module actions whose Trakt endpoint is going away or gone
package deprecation

import (
	"github.com/mfederowicz/trakt-sync/consts"
)

// Stage is how far the deprecation of an endpoint has gone
type Stage int

const (
	// StageDeprecated means the action still sends its request and a notice is shown
	StageDeprecated Stage = iota + consts.OneValue
	// StageLocked means no request is sent any more; the code is still in place
	StageLocked
	// StageRemoved means the code behind the action is gone; only the notice is left
	StageRemoved
)

const (
	actionSeparator = " -a "
	deprecatedSince = " is deprecated since "
	lockedSince     = " is locked since "
	removedSince    = " is removed since "
	reasonSeparator = ": "
	routeSeparator  = " "
	noRequestSent   = "; no request was sent"
)

// Notice is the deprecation history of one module action and its Trakt endpoint
type Notice struct {
	// Module is the CLI module name
	Module string
	// Action is the CLI action name
	Action string
	// Method is the HTTP method of the endpoint
	Method string
	// Route is the endpoint path template, without a leading slash
	Route string
	// Reason says what happened
	Reason string
	// Deprecated is the date (YYYY-MM-DD) the deprecation started
	Deprecated string
	// Locked is the date requests stopped being sent; empty while they still are
	Locked string
	// Removed is the date the code behind the action was removed; empty while it exists
	Removed string
}

// Stage returns the current stage of the notice
func (n Notice) Stage() Stage {
	if n.Removed != consts.EmptyString {
		return StageRemoved
	}
	if n.Locked != consts.EmptyString {
		return StageLocked
	}
	return StageDeprecated
}

// Blocked reports whether requests to the endpoint are no longer sent
func (n Notice) Blocked() bool {
	return n.Stage() >= StageLocked
}

// Message returns the notice for the CLI action, with its current stage
func (n Notice) Message() string {
	return n.Module + actionSeparator + n.Action + n.status()
}

// RouteMessage returns the notice for the endpoint, with its current stage
func (n Notice) RouteMessage() string {
	return n.Method + routeSeparator + n.Route + n.status()
}

// status returns the current stage with its date; earlier stages stay in the registry only
func (n Notice) status() string {
	text := deprecatedSince + n.Deprecated
	switch n.Stage() {
	case StageRemoved:
		text = removedSince + n.Removed
	case StageLocked:
		text = lockedSince + n.Locked
	}
	text += reasonSeparator + n.Reason
	if n.Blocked() {
		text += noRequestSent
	}
	return text
}
