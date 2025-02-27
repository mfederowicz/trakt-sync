// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CertificationsService  handles communication with the certifications related
// methods of the Trakt API.
type CertificationsService Service

// GetCertifications Get a list of all certifications, including names, slugs, and descriptions.
//
// API docs: https://trakt.docs.apiary.io/#reference/certifications/list/get-certifications
func (c *CertificationsService) GetCertifications(ctx context.Context, strType *string) (*str.Certifications, *str.Response, error) {
	var url = fmt.Sprintf("certifications/%s", *strType)
	printer.Println("fetch certifications url:" + url)
	
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.Certifications)
	resp, err := c.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch certifications err:", err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}


