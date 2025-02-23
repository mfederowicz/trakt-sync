<!-- TOC -->

- [trakt-sync](#trakt-sync)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage) 
    - [Command Line Flags](#command-line-flags)
    - [Command Line Commands](#command-line-commands)
    - [Sample export usage](#sample-export-usage) 
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

- `help` - Help on the trakt-sync command and subcommands.
- `history` - Returns movies and episodes that a user has watched, sorted by most recent.
- `watchlist` - Returns all items in a user's watchlist filtered by type.
- `collection` - Get all collected items in a user's collection.
- `users` - Returns all data for a user.
- `people` - Returns all data for selected person.
- `calendars` - By default, the calendar will return all shows or movies for the specified time period and can be global or user specific.
- `search` - Searches can use queries or ID lookups.
- `lists` - Returns data about lists: trending, popular, list, likes, like, items, comments.
### Sample export usage

#### Export all tvshows from the watching history:
```console
$ ./trakt-sync history -t shows -> export_history_shows_imdb.json
```

#### Export all episodes from the watching history:
```console
$ ./trakt-sync history -t episodes -f tmdb -> export_history_episodes_tmdb.json
```
```console
$ ./trakt-sync history -t episodes -f imdb -> export_history_episodes_imdb.json
```

#### Export all movies from watchlist:

```console
$ ./trakt-sync watchlist -t movies -f tmdb -> export_watchlist_movies_tmdb.json 
```
```console
$ ./trakt-sync watchlist -t movies -f imdb -> export_watchlist_movies_imdb.json
```
#### Export movies or shows from collection extended with metadata:
```console
$ ./trakt-sync collection -t movies --ex metadata
```
```console
$ ./trakt-sync collection -t shows --ex metadata
```
#### Export movies or shows or episodes from user lists:
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t episodes
```
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t shows
```
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t movies
```
#### Fetch lists for selected user:
```console
$ ./trakt-sync users -a lists -u username 
```
#### Fetch saved filters for selected user:
```console
$ ./trakt-sync users -a saved_filters -u username 
```
#### Fetch stats for selected user:
```console
$ ./trakt-sync users -a stats -u sean 
```
#### Fetch settings for current user:
```console
$ ./trakt-sync users -a settings 
```
#### Fetch watched movies for selected user:
```console
$ ./trakt-sync users -a watched -t movies -u sean 
```
#### Fetch watched shows for selected user:
```console
$ ./trakt-sync users -a watched -t shows -u sean 
```
#### Fetch watched shows for selected user without seasons:
```console
$ ./trakt-sync users -a watched -t shows -u sean --ex noseasons
```

#### Export people data:
```console
$ ./trakt-sync people -a updates -start_date 2024-10-13
```
```console
$ ./trakt-sync people -a updated_ids -start_date 2024-10-13
```
```console
$ ./trakt-sync people -a summary -i john-wayne
```
```console
$ ./trakt-sync people -a movies -i john-wayne
```
```console
$ ./trakt-sync people -a shows -i john-wayne
```
```console
$ ./trakt-sync people -a lists -i john-wayne
```
#### Export all or my calendars:

```console
$ ./trakt-sync calendars -a all-shows -> export_calendars_shows_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all-new-shows -> export_calendars_new_shows_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all-season-premieres -> export_calendars_season_premieres_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all-finales -> export_calendars_finales_20240707_7.json 
```
```console
$ ./trakt-sync calendars -a all-movies -> export_calendars_movies_20240707_7.json  
```
```console
$ ./trakt-sync calendars -a all-dvd -> export_calendars_dvd_20240707_7.json
```

#### Export search result by Text Query:

```console
$  ./trakt-sync search -a text-query -t movie -q freddy --field title
```
```console
$  ./trakt-sync search -a text-query -t movie -t show -q freddy --field tagline
```
```console
$  ./trakt-sync search -a text-query -t movie -t show -t list -q freddy --field name
```
```console
$  ./trakt-sync search -a text-query -t movie -t show -t list -q freddy --field title
```
```console
$  ./trakt-sync search -a text-query -t person -t list -q freddy --field name
```
```console
$  ./trakt-sync search -a text-query -t movie -t show -t list -q freddy --field title
```

#### Export search result by Id lookup:

```console
$ ./trakt-sync search -a id-lookup -i 12601 -t movie -t show
```
```console
$ ./trakt-sync search -a id-lookup --id_type tvdb -i 12601 -t movie -t show
```
```console
$ ./trakt-sync search -a id-lookup --id_type imdb -i 12601 -t movie
```
```console
$ ./trakt-sync search -a id-lookup --id_type imdb -i 12601 -t podcast
```
```console
$ ./trakt-sync search -a id-lookup --id_type imdb -i tt0266697
```
```console
$ ./trakt-sync search -a id-lookup --id_type tvdb -i 75725
```
```console
$ ./trakt-sync search -a id-lookup --id_type tvdb -i 75725 -t podcast
```
```console
$ ./trakt-sync search -a id-lookup -i 75725 
```
```console
$ ./trakt-sync search -a id-lookup -i 75725 -t episode
```
```console
$ ./trakt-sync search -a id-lookup --id_type tmdb -i 254265
```
#### Export lists data:
```console
$ ./trakt-sync lists -a trending
```
```console
$ ./trakt-sync lists -a popular
```
```console
$ ./trakt-sync lists -a list -trakt_id 2142753
```
```console
$ ./trakt-sync lists -a likes -trakt_id 2142753
```
```console
$ ./trakt-sync lists -a like -trakt_id 2142753
```
```console
$ ./trakt-sync lists -a like -trakt_id 2142753 -remove
```
```console
$ ./trakt-sync lists -a items -trakt_id 2142753
```
```console
$ ./trakt-sync lists -a items -trakt_id 2142753 -t movie,show
```

-- (temp not working - problems with api endpoint)
```console
$ ./trakt-sync lists -a comments -trakt_id 2142753 
```

## License

[MIT](./LICENSE)

