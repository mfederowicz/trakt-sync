// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// SearchHandler interface to handle search module action
type SearchHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

func noSearchTypeOrInvalidConfigTypeSlice(options *str.Options, slice []string) bool {
	return (options.Action == consts.TextQuery && len(options.SearchType) == consts.ZeroValue) || !cfg.IsValidConfigTypeSlice(slice, options.SearchType)
}

func validSearchIDTypes(options *str.Options, slice []string) bool {
	return len(options.SearchIDType) > consts.ZeroValue && !cfg.IsValidConfigType(slice, options.SearchIDType)
}

func checkSearchFieldFlag(options *str.Options) error {
	if len(options.SearchType) > consts.ZeroValue {
		for _, stype := range options.SearchType {
			if !cfg.IsValidConfigTypeSlice(cfg.SearchFieldConfig[stype], options.SearchField) {
				return fmt.Errorf("invalid --field flag values: %v for selected type: %v, avalable values:%v",
					options.SearchField, stype, cfg.SearchFieldConfig[stype])
			}
		}
	}

	return nil
}

func checkSearchRequiredFields(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}

	// Check search_type flag slice
	if noSearchTypeOrInvalidConfigTypeSlice(options, moduleConfig.SearchType) {
		return fmt.Errorf("invalid -t flag values: %v, avaliable values: %v", options.SearchType, moduleConfig.SearchType)
	}

	// Check search_field flag slice
	sfError := checkSearchFieldFlag(options)
	if sfError != nil {
		return fmt.Errorf("field error:%s", sfError)
	}

	// Check id_type values
	if validSearchIDTypes(options, moduleConfig.SearchIDType) {
		return fmt.Errorf("invalid --id_type flag value: %v avalable values:%v", options.SearchIDType, moduleConfig.SearchIDType)
	}

	return nil
}
