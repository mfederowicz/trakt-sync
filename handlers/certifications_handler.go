// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CertificationsHandler interface to handle certifications
type CertificationsHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
