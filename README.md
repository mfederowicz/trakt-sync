<!-- TOC -->

- [trakt-sync](#trakt-sync)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Command Line Flags](#command-line-flags)
    - [Command Line Commands](#command-line-commands)
    - [Sample export usage](#sample-export-usage)
      - [Calendars](#Calendars)
      - [Certifications](#Certifications)
      - [Checkin](#Checkin)
      - [Collection](#Collection)
      - [Comments](#Comments)
      - [Countries](#Countries)
      - [Genres](#Genres)
      - [Help](#Help)
      - [History](#History)
      - [Languages](#Languages)
      - [Lists](#Lists)
      - [Movies](#Movies)
      - [Networks](#Networks)
      - [Notes](#Notes)
      - [People](#People)
      - [Recommendations](#Recommendations)
      - [Scrobble](#Scrobble)
      - [Search](#Search)
      - [Shows](#Shows)
      - [Seasons](#Seasons)
      - [Episodes](#Episodes)
      - [Sync](#Sync)
      - [Users](#Users)
      - [Watchlist](#Watchlist)
  - [License](#license)

<!-- /TOC -->

## Installation
```bash
go install github.com/mfederowicz/trakt-sync@latest
```
## Configuration

After install, we should create [API app](https://trakt.tv/oauth/applications/new) and save credentials in config file (`$HOME/trakt-sync.toml`):
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
- `certifications` - Certifications list
- `checkin` - Checkin movie,episode,show_episode,delete
- `comments` - Comments comments,comment,replies,item,likes,like,trending,recent,updates.
- `collection` - Get all collected items in a user's collection.
- `countries` - Get a list of all countries, including names and codes.
- `genres` - Get a list of all genres, including names and slugs.
- `help` - Help on the trakt-sync command and subcommands.
- `history` - Returns movies and episodes that a user has watched, sorted by most recent.
- `languages` - Get a list of all laguages, including names and codes.
- `lists` - Returns data about lists: trending, popular, list, likes, like, items, comments.
- `movies` - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- `networks` - Get a list of all TV networks
- `notes` - Manage notes created by user
- `people` - Returns all data for selected person.
- `recommendations` - Recommendations manage movie and shows recommendations for user
- `scrobble` - Scrobble for start/pause/stop movie,show,episode
- `search` - Searches can use queries or ID lookups.
- `seasons` - Returns data about seasons: summary, season, episodes, translations, comments etc...
- `shows` - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- `sync` - Sync data useful for mediacenters: activities, playbacks, collections, ratings, watchlists, favorites.
- `users` - Returns all data for a user.
- `watchlist` - Returns all items in a user's watchlist filtered by type.

## License

[MIT](./LICENSE)

