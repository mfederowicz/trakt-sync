// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CertificationsHandler interface to handle certifications
type CertificationsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
