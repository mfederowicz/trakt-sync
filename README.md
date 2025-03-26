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

- `calendars` - By default, the calendar will return all shows or movies for the specified time period and can be global or user specific.
- `certifications` - Certifications list
- `checkin` - Checkin movie,episode,show_episode,delete
- `collection` - Get all collected items in a user's collection.
- `comments` - Comments comments,comment,replies,item,likes,like,trending,recent,updates.
- `countries` - Get a list of all countries, including names and codes.
- `genres` - Get a list of all genres, including names and slugs.
- `help` - Help on the trakt-sync command and subcommands.
- `history` - Returns movies and episodes that a user has watched, sorted by most recent.
- `languages` - Get a list of all laguages, including names and codes.
- `lists` - Returns data about lists: trending, popular, list, likes, like, items, comments.
- `movies` - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- `people` - Returns all data for selected person.
- `search` - Searches can use queries or ID lookups.
- `users` - Returns all data for a user.
- `watchlist` - Returns all items in a user's watchlist filtered by type.
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
#### Manage comments:

```console
$ ./trakt-sync comments -a comment -comment_id 779883 -comment "minions,minions,minions movie ever ok" 
```
```console
$ ./trakt-sync comments -a comment -comment_id 779883 -delete
```
```console
$ ./trakt-sync comments -a comments -t episode -trakt_id 172245 -comment "super episode, interesting plot ok" 
```
```console
$ ./trakt-sync comments -a replies -comment_id 779896 -reply "reply msg min 5 words" -spoiler 
```
```console
$ ./trakt-sync comments -a replies -comment_id 71340
```
```console
$ ./trakt-sync comments -a item -comment_id 664237 -ex full
```
```console
$ ./trakt-sync comments -a likes -comment_id 773108 -remove
```
```console
$ ./trakt-sync comments -a like -comment_id 773108
```
```console
$ ./trakt-sync comments -a like -comment_id 773108 -remove
```
```console
$ ./trakt-sync comments -a trending -comment_type reviews
```
```console
$ ./trakt-sync comments -a recent -include_replies false
```
```console
$ ./trakt-sync comments -a recent -include_replies true
```
```console
$ ./trakt-sync comments -a updates -include_replies false
```

#### Export countries list for movies or shows:
```console
$ ./trakt-sync countries -> export_countries_movies.json
```
```console
$ ./trakt-sync countries -t movies -> export_countries_movies.json
```
```console
$ ./trakt-sync countries -t shows -> export_countries_shows.json
```

#### Export languages list for movies or shows:
```console
$ ./trakt-sync languages -> export_languages_movies.json
```
```console
$ ./trakt-sync languages -t movies -> export_languages_movies.json
```
```console
$ ./trakt-sync languages -t shows -> export_languages_shows.json
```

#### Export genres list for movies or shows:
```console
$ ./trakt-sync genres -> export_genres_movies.json
```
```console
$ ./trakt-sync genres -t movies -> export_genres_movies.json
```
```console
$ ./trakt-sync genres -t shows -> export_genres_shows.json
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
#### Export certifications for movies or shows:
```console
$ ./trakt-sync certifications -> export_certifications_movies.json
```
```console
$ ./trakt-sync certifications -t movies -> export_certifications_movies.json
```
```console
$ ./trakt-sync certifications -t shows -> export_certifications_shows.json
```

#### Manage checkin:

```console
$ ./trakt-sync checkin -a movie -trakt_id 28 -msg "super movie"
```
```console
$ ./trakt-sync checkin -a episode -trakt_id 3190486 -msg "super episode"
```
```console
$ ./trakt-sync checkin -a show_episode -trakt_id 37696 -episode_abs 6 -msg "super episode"
```
```console
$ ./trakt-sync checkin -a show_episode -trakt_id 136121 -episode_code 1x5 -msg "super episode"
```
```console
$ ./trakt-sync checkin -a delete 
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
#### Export movies data:
```console
$ ./trakt-sync movies -a trending
```
```console
$ ./trakt-sync movies -a popular
```
```console
$ ./trakt-sync movies -a favorited -period daily
$ ./trakt-sync movies -a favorited -period weekly
$ ./trakt-sync movies -a favorited -period monthly
$ ./trakt-sync movies -a favorited -period all
```
```console
$ ./trakt-sync movies -a played -period daily
$ ./trakt-sync movies -a played -period weekly
$ ./trakt-sync movies -a played -period monthly
$ ./trakt-sync movies -a played -period all
```
```console
$ ./trakt-sync movies -a watched -period daily
$ ./trakt-sync movies -a watched -period weekly
$ ./trakt-sync movies -a watched -period monthly
$ ./trakt-sync movies -a watched -period all
```
```console
$ ./trakt-sync movies -a collected -period daily
$ ./trakt-sync movies -a collected -period weekly
$ ./trakt-sync movies -a collected -period monthly
$ ./trakt-sync movies -a collected -period all
```
```console
$ ./trakt-sync movies -a anticipated
```
```console
$ ./trakt-sync movies -a boxoffice
```
```console
$ ./trakt-sync movies -a updates -start_date YYYY-MM-DD
```
```console
$ ./trakt-sync movies -a updated_ids -start_date YYYY-MM-DD
```
```console
$ ./trakt-sync movies -a summary -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a aliases -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a releases -i tron-legacy-2010 -country us
```
```console
$ ./trakt-sync movies -a translations -i tron-legacy-2010 -language es
```
```console
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s newest
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s oldest
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s likes
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s replies
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s highest
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s lowest
$ ./trakt-sync movies -a comments -i tron-legacy-2010 -s plays
```
```console
$ ./trakt-sync movies -a lists -i tron-legacy-2010 -t all -s popular
$ ./trakt-sync movies -a lists -i tron-legacy-2010 -t all -s likes
$ ./trakt-sync movies -a lists -i tron-legacy-2010 -t all -s comments
$ ./trakt-sync movies -a lists -i tron-legacy-2010 -t all -s items
$ ./trakt-sync movies -a lists -i tron-legacy-2010 -t all -s added
$ ./trakt-sync movies -a lists -i tron-legacy-2010 -t all -s updated
```
```console
$ ./trakt-sync movies -a people -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a ratings -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a related -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a stats -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a studios -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a watching -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a videos -i tron-legacy-2010
```
```console
$ ./trakt-sync movies -a refresh -i tron-legacy-2010
```

## License

[MIT](./LICENSE)

