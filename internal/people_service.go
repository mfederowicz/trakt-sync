// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// PeopleService  handles communication with the people related
// methods of the Trakt API.
type PeopleService Service

// GetListsContainingThisPerson Returns all lists that contain this person.
//
// API docs: https://trakt.docs.apiary.io/#reference/people/lists/get-lists-containing-this-person
func (p *PeopleService) GetListsContainingThisPerson(ctx context.Context, id *string, typeString *string, sort *string, opts *uri.ListOptions) ([]*str.PersonalList, *str.Response, error) {
	var url string

	url = fmt.Sprintf("people/%s/lists/%s/%s", *id, *typeString, *sort)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch lists url:" + url)
	req, err := p.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.PersonalList{}
	resp, err := p.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetMovieCredits Returns all movies where this person is in the cast or crew.
//
// API docs: https://trakt.docs.apiary.io/#reference/people/movies/get-movie-credits
func (p *PeopleService) GetMovieCredits(ctx context.Context, id *string, opts *uri.ListOptions) (*str.PersonMovies, *str.Response, error) {
	var url = fmt.Sprintf("people/%s/movies", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch movie credits url:" + url)
	req, err := p.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.PersonMovies)
	resp, err := p.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch person err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetShowCredits Returns all shows where this person is in the cast or crew, including the episode_count for which they appear.
//
// API docs: https://trakt.docs.apiary.io/#reference/people/shows/get-show-credits
func (p *PeopleService) GetShowCredits(ctx context.Context, id *string, opts *uri.ListOptions) (*str.PersonShows, *str.Response, error) {
	var url = fmt.Sprintf("people/%s/shows", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch shows credits url:" + url)
	req, err := p.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.PersonShows)
	resp, err := p.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch shows credits err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSinglePerson Returns a single person's details.
//
// API docs: https://trakt.docs.apiary.io/#reference/people/summary/get-a-single-person
func (p *PeopleService) GetSinglePerson(ctx context.Context, id *string, opts *uri.ListOptions) (*str.Person, *str.Response, error) {
	var url string

	url = fmt.Sprintf("people/%s", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch person url:" + url)
	req, err := p.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Person)
	resp, err := p.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch person err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetRecentlyUpdatedPeople Returns all people updated since the specified UTC date and time.
//
// API docs: https://trakt.docs.apiary.io/#reference/people/updates/get-recently-updated-people
func (p *PeopleService) GetRecentlyUpdatedPeople(ctx context.Context, startDate *string, opts *uri.ListOptions) ([]*str.PersonItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("people/updates/%s", *startDate)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updates url:" + url)
	req, err := p.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.PersonItem{}
	resp, err := p.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetRecentlyUpdatedPeopleTraktIDs Returns all people Trakt IDs updated since the specified UTC date and time.
//
// API docs: https://trakt.docs.apiary.io/#reference/people/updated-ids
func (p *PeopleService) GetRecentlyUpdatedPeopleTraktIDs(ctx context.Context, startDate *string, opts *uri.ListOptions) ([]*int, *str.Response, error) {
	var url string

	url = fmt.Sprintf("people/updates/id/%s", *startDate)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updates url:" + url)
	req, err := p.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*int{}
	resp, err := p.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
