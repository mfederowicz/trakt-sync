[![Test](https://github.com/mfederowicz/trakt-sync/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/mfederowicz/trakt-sync/actions/workflows/test.yaml?query=branch%3Amain)
[![Go Reference](https://pkg.go.dev/badge/github.com/mfederowicz/trakt-sync/trakt.svg)](https://pkg.go.dev/github.com/mfederowicz/trakt-sync/trakt)
[![Coverage](https://img.shields.io/badge/coverage-80%25-green)](https://github.com/mfederowicz/trakt-sync/actions/workflows/test.yaml?query=branch%3Amain)
[![Version](https://img.shields.io/badge/version-v1.21.0-blue)](https://github.com/mfederowicz/trakt-sync/releases/latest)

<!-- TOC -->

- [trakt-sync](#trakt-sync)
  - [Intended use](#intended-use)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Command Line Flags](#command-line-flags)
    - [Command Line Commands](#command-line-commands)
  - [API documentation](#api-documentation)
  - [License](#license)

<!-- /TOC -->

## Intended use

`trakt-sync` is a tool for managing and exporting **your own** Trakt data (history, collection,
watchlist, ratings, lists, ...). It must not be used to bulk-collect public data such as other
users' lists, comments or ratings, or to feed that data into other services. Paging options
(`per_page`, `pages_limit`) exist to fetch your own data efficiently, not to crawl the API.
Every use has to follow the [Trakt API Use Policy](https://developer.trakt.tv/?section=guides&guide=api-use-policy).

## Installation
```bash
go install github.com/mfederowicz/trakt-sync@latest
```
## Configuration

After install, we need API credentials (Client ID and Client Secret). Create a new API app in the [Trakt developer portal](https://developer.trakt.tv/apps/new) (a verified GitHub account is required) and save them in config file (`$HOME/trakt-sync.toml`):
```console
client_id = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
client_secret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxx"
token_path = "~/.config/trakt-sync/token.json"
errorCode = 0
warningCode = 0
per_page = 500
pages_limit = 10
```

## Usage

`trakt-sync` supports a `-config` flag whose value should correspond to a TOML file.
If not provided, `trakt-sync` will try to use a global config file (assumed to be located at `$HOME/trakt-sync.toml`).
Otherwise, if no configuration TOML file is found then `trakt-sync` uses a built-in parameters depends on selected module.
A value given on the command line (for example `-t shows`) wins over the same option in the config file (`type = "shows"`), which wins over the built-in default.
One exception: `type` in the config file is ignored by `episodes`, `lists`, `movies`, `people`, `scrobble`, `seasons`, `shows` and `users`. There `-t` has a module-specific meaning (for example a list type), so give it on the command line.

### Command Line Flags

`trakt-sync` accepts the following command line parameters:

- `-config [PATH]` - path to config file in TOML format, defaults to `$HOME/trakt-sync.toml` if present.
- `-version` - get trakt-sync version.

### Command Line Commands

`trakt-sync` accepts the following command line commands/modules:

- [`calendars`](./docs/calendars.md) - By default, the calendar will return all shows or movies for the specified time period and can be global or user specific.
- [`certifications`](./docs/certifications.md) - Certifications list
- [`checkin`](./docs/checkin.md) - Checkin movie,episode,show_episode,delete
- [`comments`](./docs/comments.md) - Comments comments,comment,replies,item,likes,like,trending,recent,updates,reactions,reactions_summary,reaction,report.
- [`collection`](./docs/collection.md) - Get all collected items in a user's collection.
- [`countries`](./docs/countries.md) - Get a list of all countries, including names and codes.
- [`episodes`](./docs/episodes.md) - Returns data about episodes: summary, season, episodes, translations, comments etc...
- [`genres`](./docs/genres.md) - Get a list of all genres, including names and slugs.
- `help` - Help on the trakt-sync command and subcommands.
- [`history`](./docs/history.md) - Returns movies and episodes that a user has watched, sorted by most recent.
- [`languages`](./docs/languages.md) - Get a list of all laguages, including names and codes.
- [`lists`](./docs/lists.md) - Returns data about lists: trending, popular, list, likes, like, items, comments, report.
- [`media`](./docs/media.md) - Returns movies and shows together: trending, popular, anticipated.
- [`movies`](./docs/movies.md) - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- [`networks`](./docs/networks.md) - Get a list of all TV networks
- [`notes`](./docs/notes.md) - Manage notes created by user
- [`people`](./docs/people.md) - Returns all data for selected person.
- [`recommendations`](./docs/recommendations.md) - Recommendations manage movie and shows recommendations for user
- [`scrobble`](./docs/scrobble.md) - Scrobble for start/pause/stop movie,show,episode
- [`search`](./docs/search.md) - Searches can use queries or ID lookups.
- [`seasons`](./docs/seasons.md) - Returns data about seasons: summary, season, episodes, translations, comments etc...
- [`shows`](./docs/shows.md) - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- [`smart_lists`](./docs/smart_lists.md) - Returns a smart list definition or the items it resolves to.
- [`social_recommendations`](./docs/social_recommendations.md) - Movie and show recommendations based on the people you follow.
- [`sync`](./docs/sync.md) - Sync data useful for mediacenters: activities, playbacks, collections, ratings, watchlists, favorites.
- [`team`](./docs/team.md) - Returns Trakt team members.
- [`users`](./docs/users.md) - Returns all data for a user.
- [`watchlist`](./docs/watchlist.md) - Returns all items in a user's watchlist filtered by type.
- [`watchnow`](./docs/watchnow.md) - Returns watch now sources (streaming providers), all or by country.
- [`younify`](./docs/younify.md) - Streaming service connections: connections, connect, refresh, disconnect.

## Library usage

The API client behind the CLI is also a Go package, `github.com/mfederowicz/trakt-sync/trakt`.
It is **experimental**: its API may still change in minor releases, listed in the changelog with a
**Library:** prefix. See the versioning rules at the top of [CHANGELOG.md](./CHANGELOG.md).

```bash
go get github.com/mfederowicz/trakt-sync/trakt
```
```go
client := trakt.NewClient(nil).
	WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
	WithUserAgent("my-app/1.0")

movies, _, err := client.Movies.GetTrendingMovies(context.Background(), &uri.ListOptions{Limit: 10})
```
Package docs: [pkg.go.dev/github.com/mfederowicz/trakt-sync/trakt](https://pkg.go.dev/github.com/mfederowicz/trakt-sync/trakt).
Runnable programs (device login, trending, paginated history, errors) are in [`example/`](./example).
The [Trakt API Use Policy](https://developer.trakt.tv/?section=guides&guide=api-use-policy) applies to every app built on it.

## API documentation

- [Trakt developer portal](https://developer.trakt.tv) - official guides (authentication, pagination, rate limiting, required headers) and management of your API apps.
- [trakt/trakt-api](https://github.com/trakt/trakt-api) - API contract used as the source of truth for this project; API announcements are published there.
- [API coverage](./API_COVERAGE.md) - which Trakt API routes trakt-sync implements, and what is still missing.

## License

[MIT](./LICENSE)

