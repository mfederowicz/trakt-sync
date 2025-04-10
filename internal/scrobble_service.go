// Package internal used for client and services
package internal

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// ScrobbleService  handles communication with the scrobble related
// methods of the Trakt API.
type ScrobbleService Service

// StartScrobble Start watching in a media server.
// API docs:https://trakt.docs.apiary.io/#reference/scrobble/start/start-watching-in-a-media-center 
func (s *ScrobbleService) StartScrobble(ctx context.Context, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error) {
	var url = "scrobble/start"
	printer.Println("start scrobble")
	req, err := s.client.NewRequest(http.MethodPost, url, scrobble)
	if err != nil {
		return nil, nil, err
	}

	sc := new(str.Scrobble)
	resp, err := s.client.Do(ctx, req, sc)
	if err != nil {
		return sc, resp, err
	}

	return sc, resp, nil
}

// PauseScrobble Pause watching in a media server.
// API docs: https://trakt.docs.apiary.io/#reference/scrobble/pause/pause-watching-in-a-media-center 
func (s *ScrobbleService) PauseScrobble(ctx context.Context, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error) {
	var url = "scrobble/pause"
	printer.Println("pause scrobble")
	req, err := s.client.NewRequest(http.MethodPost, url, scrobble)
	if err != nil {
		return nil, nil, err
	}

	sc := new(str.Scrobble)
	resp, err := s.client.Do(ctx, req, sc)
	if err != nil {
		return sc, resp, err
	}

	return sc, resp, nil
}


// StopScrobble Stop watching in a media server.
// API docs: https://trakt.docs.apiary.io/#reference/scrobble/stop/stop-or-finish-watching-in-a-media-center 
func (s *ScrobbleService) StopScrobble(ctx context.Context, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error) {
	var url = "scrobble/stop"
	printer.Println("stop scrobble")
	req, err := s.client.NewRequest(http.MethodPost, url, scrobble)
	if err != nil {
		return nil, nil, err
	}

	sc := new(str.Scrobble)
	resp, err := s.client.Do(ctx, req, sc)
	if err != nil {
		return sc, resp, err
	}

	return sc, resp, nil
}

