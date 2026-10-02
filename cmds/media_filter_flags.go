// Package cmds used for commands modules
package cmds

import (
	"flag"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

// mediaFilterFlags are the media filter flags a module registers in its own flag set;
// -genres, -years, -countries and -runtimes are global flags.
type mediaFilterFlags struct {
	watchNow       *string
	subgenres      *string
	ratings        *string
	certifications *string
	startDate      *string
	endDate        *string
}

// newMediaFilterFlags registers the media filter flags in fs.
func newMediaFilterFlags(fs *flag.FlagSet) mediaFilterFlags {
	return mediaFilterFlags{
		watchNow:       fs.String("watchnow", consts.EmptyString, consts.WatchNowUsage),
		subgenres:      fs.String("subgenres", consts.EmptyString, consts.SubgenresUsage),
		ratings:        fs.String("ratings", consts.EmptyString, consts.RatingsFilterUsage),
		certifications: fs.String("certifications", consts.EmptyString, consts.CertificationsUsage),
		startDate:      fs.String("start_date", consts.EmptyString, consts.FilterStartDateUsage),
		endDate:        fs.String("end_date", consts.EmptyString, consts.EndDateUsage),
	}
}

// newMediaFilterFlagsWithoutDates registers the media filter flags in fs, except -start_date and -end_date,
// for a module where those names are taken or make no sense.
func newMediaFilterFlagsWithoutDates(fs *flag.FlagSet) mediaFilterFlags {
	return mediaFilterFlags{
		watchNow:       fs.String("watchnow", consts.EmptyString, consts.WatchNowUsage),
		subgenres:      fs.String("subgenres", consts.EmptyString, consts.SubgenresUsage),
		ratings:        fs.String("ratings", consts.EmptyString, consts.RatingsFilterUsage),
		certifications: fs.String("certifications", consts.EmptyString, consts.CertificationsUsage),
	}
}

// apply copies the media filter flags and the global filter flags to options.
func (m mediaFilterFlags) apply(options str.Options) str.Options {
	options.WatchNow = *m.watchNow
	options.Subgenres = *m.subgenres
	options.Ratings = *m.ratings
	options.Certifications = *m.certifications
	if m.startDate != nil {
		options.MediaStartDate = *m.startDate
	}
	if m.endDate != nil {
		options.MediaEndDate = *m.endDate
	}
	options.Genres = *_genres
	options.Years = *_years
	options.Countries = *_countries
	options.Runtimes = *_runtimes
	return options
}
