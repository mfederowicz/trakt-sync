// Package handlers used to handle module actions
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersSmartListsHandler struct for handler
type UsersSmartListsHandler struct{}

// Handle to handle users: smart_lists action
func (UsersSmartListsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all smart lists of: " + options.UserName)
	result, _, err := client.Users.GetSmartLists(client.BuildCtxFromOptions(options), &options.UserName)
	if err != nil {
		return fmt.Errorf("fetch smart lists error: %w", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeSmartList(options, result)
}

// readSmartListWrite reads a smart list create / update body from the -items file or stdin; unknown keys are rejected.
func readSmartListWrite(common *CommonLogic, options *str.Options) (*str.SmartListWrite, error) {
	data, err := common.ReadInputBytes(*options)
	if err != nil {
		return nil, err
	}

	list := new(str.SmartListWrite)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(list); err != nil {
		return nil, fmt.Errorf("invalid smart list JSON: %w", err)
	}

	return list, nil
}

// validSmartListWrite checks the contract values of a smart list body; create also needs name, source and media_type.
func validSmartListWrite(list *str.SmartListWrite, create bool) error {
	if create && (list.Name == nil || len(*list.Name) == consts.ZeroValue || list.Source == nil || list.MediaType == nil) {
		return errors.New("smart list needs name, source and media_type")
	}
	if !create && list.Name == nil && list.Source == nil && list.MediaType == nil && list.Filters == nil && list.Privacy == nil {
		return errors.New("smart list update needs at least one of name, source, media_type, filters, privacy")
	}
	if list.Name != nil && len(*list.Name) == consts.ZeroValue {
		return errors.New("smart list name must not be empty")
	}

	checks := []struct {
		name  string
		value *string
		valid []string
	}{
		{name: "source", value: list.Source, valid: cfg.SmartListSources},
		{name: "media_type", value: list.MediaType, valid: cfg.SmartListMediaTypes},
		{name: "privacy", value: list.Privacy, valid: cfg.SmartListPrivacy},
	}
	// a field that is sent must hold a contract value; cfg.IsValidConfigType would accept ""
	for _, c := range checks {
		if c.value != nil && !slices.Contains(c.valid, *c.value) {
			return fmt.Errorf("%s '%s' is not valid, avaliable values: %v", c.name, *c.value, c.valid)
		}
	}

	return nil
}
