<!-- TOC -->

- [trakt-sync](#trakt-sync)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Command Line Flags](#command-line-flags)
    - [Command Line Commands](#command-line-commands)
  - [License](#license)

<!-- /TOC -->

## Installation
```bash
go install github.com/mfederowicz/trakt-sync@latest
```
## Configuration

After install, we need [API credentials](https://docs.trakt.tv/docs/create-an-app) (Client ID and Client Secret) and save them in config file (`$HOME/trakt-sync.toml`):
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

### Command Line Flags

`trakt-sync` accepts the following command line parameters:

- `-config [PATH]` - path to config file in TOML format, defaults to `$HOME/trakt-sync.toml` if present.
- `-version` - get trakt-sync version.

### Command Line Commands

`trakt-sync` accepts the following command line commands/modules:

- [`calendars`](./docs/calendars.md) - By default, the calendar will return all shows or movies for the specified time period and can be global or user specific.
- [`certifications`](./docs/certifications.md) - Certifications list
- [`checkin`](./docs/checkin.md) - Checkin movie,episode,show_episode,delete
- [`comments`](./docs/comments.md) - Comments comments,comment,replies,item,likes,like,trending,recent,updates.
- [`collection`](./docs/collection.md) - Get all collected items in a user's collection.
- [`countries`](./docs/countries.md) - Get a list of all countries, including names and codes.
- [`genres`](./docs/genres.md) - Get a list of all genres, including names and slugs.
- `help` - Help on the trakt-sync command and subcommands.
- [`history`](./docs/history.md) - Returns movies and episodes that a user has watched, sorted by most recent.
- [`languages`](./docs/languages.md) - Get a list of all laguages, including names and codes.
- [`lists`](./docs/lists.md) - Returns data about lists: trending, popular, list, likes, like, items, comments.
- [`movies`](./docs/movies.md) - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- [`networks`](./docs/networks.md) - Get a list of all TV networks
- [`notes`](./docs/notes.md) - Manage notes created by user
- [`people`](./docs/people.md) - Returns all data for selected person.
- [`recommendations`](./docs/recommendations.md) - Recommendations manage movie and shows recommendations for user
- [`scrobble`](./docs/scrobble.md) - Scrobble for start/pause/stop movie,show,episode
- [`search`](./docs/search.md) - Searches can use queries or ID lookups.
- [`seasons`](./docs/seasons.md) - Returns data about seasons: summary, season, episodes, translations, comments etc...
- [`shows`](./docs/shows.md) - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- [`sync`](./docs/sync.md) - Sync data useful for mediacenters: activities, playbacks, collections, ratings, watchlists, favorites.
- [`users`](./docs/users.md) - Returns all data for a user.
- [`watchlist`](./docs/watchlist.md) - Returns all items in a user's watchlist filtered by type.

## License

[MIT](./LICENSE)

