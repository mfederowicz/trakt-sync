// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/str"
)

var (
	_smartListsAction            = SmartListsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_smartListsID                = SmartListsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.SmartListIDUsage)
	_smartListsWatchNow          = SmartListsCmd.Flag.String("watchnow", consts.EmptyString, consts.WatchNowUsage)
	_smartListsSubgenres         = SmartListsCmd.Flag.String("subgenres", consts.EmptyString, consts.SubgenresUsage)
	_smartListsRatings           = SmartListsCmd.Flag.String("ratings", consts.EmptyString, consts.RatingsFilterUsage)
	_smartListsCertifications    = SmartListsCmd.Flag.String("certifications", consts.EmptyString, consts.CertificationsUsage)
	_smartListsIgnoreWatched     = SmartListsCmd.Flag.String("ignore_watched", string(cfg.DefaultConfig().IgnoreWatched), consts.IgnoreWatchedUsage)
	_smartListsIgnoreWatchlisted = SmartListsCmd.Flag.String("ignore_watchlisted", string(cfg.DefaultConfig().IgnoreWatchlisted), consts.IgnoreWatchlistedUsage)
	_smartListsWatchNowCountry   = SmartListsCmd.Flag.String("watchnow_country", consts.EmptyString, consts.SmartListCountryUsage)
	_smartListsNudity            = SmartListsCmd.Flag.String("parental_nudity", consts.EmptyString, consts.ParentalRangeUsage)
	_smartListsViolence          = SmartListsCmd.Flag.String("parental_violence", consts.EmptyString, consts.ParentalRangeUsage)
	_smartListsProfanity         = SmartListsCmd.Flag.String("parental_profanity", consts.EmptyString, consts.ParentalRangeUsage)
	_smartListsAlcohol           = SmartListsCmd.Flag.String("parental_alcohol", consts.EmptyString, consts.ParentalRangeUsage)
	_smartListsFrightening       = SmartListsCmd.Flag.String("parental_frightening", consts.EmptyString, consts.ParentalRangeUsage)
	_smartListsUnrated           = SmartListsCmd.Flag.Bool("parental_include_unrated", false, consts.ParentalUnratedUsage)
)

// SmartListsCmd returns smart list definitions and the items they resolve to.
var SmartListsCmd = &Command{
	Name:    "smart_lists",
	Usage:   "",
	Summary: "Returns a smart list definition (summary) or the items it resolves to (items).",
	Help:    `smart_lists command`,
}

func smartListsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	smartListsItemsSort(options, cmd.flagIsSet("t"), cmd.flagIsSet("sort_by"), cmd.flagIsSet("sort_how"))
	var handler handlers.SmartListsHandler
	allHandlers := map[string]handlers.Handler{
		consts.Summary: handlers.SmartListsSummaryHandler{},
		consts.Items:   handlers.SmartListsItemsHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Summary, consts.Items})
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

// smartListsItemsSort keeps -t, -sort_by and -sort_how only when given on the command line;
// the config file and built-in defaults would otherwise always pick the typed and sorted items route.
func smartListsItemsSort(options *str.Options, typeSet bool, sortBySet bool, sortHowSet bool) {
	if !typeSet {
		options.Type = consts.EmptyString
	}
	if !sortBySet {
		options.SortBy = consts.EmptyString
	}
	if !sortHowSet {
		options.SortHow = consts.EmptyString
	}
}

func init() {
	SmartListsCmd.Run = smartListsFunc
}
