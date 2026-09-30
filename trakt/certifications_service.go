// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
)

// CertificationsService  handles communication with the certifications related
// methods of the Trakt API.
type CertificationsService Service

// GetCertifications Get a list of all certifications, including names, slugs, and descriptions.
//
// API docs: https://docs.trakt.tv/reference/getcertificationslist
func (c *CertificationsService) GetCertifications(ctx context.Context, strType *string) (*str.Certifications, *str.Response, error) {
	var url = fmt.Sprintf("certifications/%s", *strType)
	c.client.debug("fetch certifications url:" + url)

	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.Certifications)
	resp, err := c.client.Do(ctx, req, &result)

	if err != nil {
		c.client.debug("fetch certifications err:", err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}
