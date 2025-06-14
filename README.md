<!-- TOC -->

- [trakt-sync](#trakt-sync)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage) 
    - [Command Line Flags](#command-line-flags)
    - [Command Line Commands](#command-line-commands)
    - [Sample export usage](#sample-export-usage) 
      - [calendars](#calendars)
      - [certifications](#certifications)
      - [checkin](#checkin)
      - [collection](#collection)
      - [comments](#comments)
      - [countries](#countries)
      - [genres](#genres)
      - [help](#help)
      - [history](#history)
      - [languages](#languages)
      - [lists](#lists)
      - [movies](#movies)
      - [networks](#networks)
      - [notes](#notes)
      - [people](#people)
      - [recommendations](#recommendations)
      - [scrobble](#scrobble)
      - [search](#search)
      - [shows](#shows)
      - [seasons](#seasons)
      - [episodes](#episodes)
      - [users](#users)
      - [watchlist](#watchlist)

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
- `networks` - Get a list of all TV networks
- `notes` - Manage notes created by user
- `people` - Returns all data for selected person.
- `recommendations` - Recommendations manage movie and shows recommendations for user
- `scrobble` - Scrobble for start/pause/stop movie,show,episode 
- `search` - Searches can use queries or ID lookups.
- `seasons` - Returns data about seasons: summary, season, episodes, translations, comments etc...
- `shows` - Returns data about movies: trending, popular, list, likes, like, items, comments etc...
- `users` - Returns all data for a user.
- `watchlist` - Returns all items in a user's watchlist filtered by type.
### Sample export usage

#### calendars:
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
#### certifications:
```console
$ ./trakt-sync certifications -> export_certifications_movies.json
```
```console
$ ./trakt-sync certifications -t movies -> export_certifications_movies.json
```
```console
$ ./trakt-sync certifications -t shows -> export_certifications_shows.json
```
#### checkin:
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
#### collection:
```console
$ ./trakt-sync collection -t movies --ex metadata
```
```console
$ ./trakt-sync collection -t shows --ex metadata
```
#### comments:
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
#### countries:
```console
$ ./trakt-sync countries -> export_countries_movies.json
```
```console
$ ./trakt-sync countries -t movies -> export_countries_movies.json
```
```console
$ ./trakt-sync countries -t shows -> export_countries_shows.json
```
#### genres:
```console
$ ./trakt-sync genres -> export_genres_movies.json
```
```console
$ ./trakt-sync genres -t movies -> export_genres_movies.json
```
```console
$ ./trakt-sync genres -t shows -> export_genres_shows.json
```
#### history:
```console
$ ./trakt-sync history -t shows -> export_history_shows_imdb.json
```
```console
$ ./trakt-sync history -t episodes -f tmdb -> export_history_episodes_tmdb.json
```
```console
$ ./trakt-sync history -t episodes -f imdb -> export_history_episodes_imdb.json
```
#### languages:
```console
$ ./trakt-sync languages -> export_languages_movies.json
```
```console
$ ./trakt-sync languages -t movies -> export_languages_movies.json
```
```console
$ ./trakt-sync languages -t shows -> export_languages_shows.json
```
#### lists:
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
#### movies:
```console
$ ./trakt-sync movies -a trending
```
```console
$ ./trakt-sync movies -a popular
```
##### Get the most favorited movies
```console
$ ./trakt-sync movies -a favorited -period daily
```
```console
$ ./trakt-sync movies -a favorited -period weekly
```
```console
$ ./trakt-sync movies -a favorited -period monthly
```
```console
$ ./trakt-sync movies -a favorited -period all
```
##### Get the most played movies
```console
$ ./trakt-sync movies -a played -period daily
```
```console
$ ./trakt-sync movies -a played -period weekly
```
```console
$ ./trakt-sync movies -a played -period monthly
```
```console
$ ./trakt-sync movies -a played -period all
```
##### Get the most watched movies
```console
$ ./trakt-sync movies -a watched -period daily
```
```console
$ ./trakt-sync movies -a watched -period weekly
```
```console
$ ./trakt-sync movies -a watched -period monthly
```
```console
$ ./trakt-sync movies -a watched -period all
```
##### Get the most collected movies
```console
$ ./trakt-sync movies -a collected -period daily
```
```console
$ ./trakt-sync movies -a collected -period weekly
```
```console
$ ./trakt-sync movies -a collected -period monthly
```
```console
$ ./trakt-sync movies -a collected -period all
```
##### Get the most anticipated movies
```console
$ ./trakt-sync movies -a anticipated
```
##### Get the weekend box office
```console
$ ./trakt-sync movies -a boxoffice
```
##### Get recenty updated movies 
```console
$ ./trakt-sync movies -a updates -start_date YYYY-MM-DD
```
##### Get recenty updated movie Trakt IDs 
```console
$ ./trakt-sync movies -a updated_ids -start_date YYYY-MM-DD
```
##### Get a movie
```console
$ ./trakt-sync movies -a summary -i the-sopranos
```
##### Get all movie aliases
```console
$ ./trakt-sync movies -a aliases -i the-sopranos
```
##### Get all movie releases
```console
$ ./trakt-sync movies -a releases -i the-sopranos -country us
```
##### Get all movie translations
```console
$ ./trakt-sync movies -a translations -i the-sopranos -language es
```
##### Get all movie comments
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s newest
```
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s oldest
```
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s likes
```
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s replies
```
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s highest
```
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s lowest
```
```console
$ ./trakt-sync movies -a comments -i the-sopranos -s plays
```
##### Get lists containing the movie
```console
$ ./trakt-sync movies -a lists -i the-sopranos -t all -s popular
```
```console
$ ./trakt-sync movies -a lists -i the-sopranos -t all -s likes
```
```console
$ ./trakt-sync movies -a lists -i the-sopranos -t all -s comments
```
```console
$ ./trakt-sync movies -a lists -i the-sopranos -t all -s items
```
```console
$ ./trakt-sync movies -a lists -i the-sopranos -t all -s added
```
```console
$ ./trakt-sync movies -a lists -i the-sopranos -t all -s updated
```
##### Get all people for movie
```console
$ ./trakt-sync movies -a people -i the-sopranos
```
##### Get movie ratings
```console
$ ./trakt-sync movies -a ratings -i the-sopranos
```
##### Get related movies
```console
$ ./trakt-sync movies -a related -i the-sopranos
```
##### Get movies stats
```console
$ ./trakt-sync movies -a stats -i the-sopranos
```
##### Get movies studios
```console
$ ./trakt-sync movies -a studios -i the-sopranos
```
##### Get users watching right now
```console
$ ./trakt-sync movies -a watching -i the-sopranos
```
##### Get all videos
```console
$ ./trakt-sync movies -a videos -i the-sopranos
```
##### Refresh movie metadata
```console
$ ./trakt-sync movies -a refresh -i the-sopranos
```
#### networks:
```console
$ ./trakt-sync networks -a list
```
#### notes:
##### adding notes for media types:
```console 
$ ./trakt-sync notes -a notes -t movie -i the-sopranos -notes "xyz"
```
```console
$ ./trakt-sync notes -a notes -t show -i breaking-bad -notes "greate show"
```
```console
$ ./trakt-sync notes -a notes -t season -i 250341 -notes "greate season"
```
```console
$ ./trakt-sync notes -a notes -t episode -i 250341 -notes "greate episode"
```
```console
$ ./trakt-sync notes -a notes -t person -i john-wayne -notes "greate person"
```
```console
$ ./trakt-sync notes -a notes -t history -i 1234567 -notes "history note"
```
##### adding notes depends on activities:
```console 
$ ./trakt-sync notes -a notes -t collection -item episode -i 73629 -notes "great episode"
```
```console 
$ ./trakt-sync notes -a notes -t collection -item movie -i despicable-me-4-2024 -notes "great animation"
```
```console 
$ ./trakt-sync notes -a notes -t rating -item movie -i despicable-me-4-2024 -notes "great animation"
```
```console 
$ ./trakt-sync notes -a notes -t rating -item episode -i 73629 -notes "overall 10/10"
```
```console 
$ ./trakt-sync notes -a notes -t rating -item movie -i the-gorge-2025 -notes "overall 7/10"
```
```console 
$ ./trakt-sync notes -a notes -t rating -item season -i 3961 -notes "overall 9/10"
```
```console 
$ ./trakt-sync notes -a notes -t rating -item show -i the-sopranos -notes "overall 9/10"
```
##### manage notes get/modify/delete:
```console
$ ./trakt-sync notes -a note -i 97857
```
```console
$ ./trakt-sync notes -a note -i 97857 -notes "super 10/10" -privacy public -spoiler
```
```console
$ ./trakt-sync notes -a note -i 97857 -delete
```
##### get items attachment to note:
```console
$ ./trakt-sync notes -a item -i 97854
```
#### people:
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
#### recommendations:
##### hide movie recommendations:
```console
$ ./trakt-sync recommendations -a movies -i black-bag-2025 -hide                                                                                                              
```
##### movies recommendations:
```console
$ ./trakt-sync recommendations -a movies                                                                                      
```
```console
$ ./trakt-sync recommendations -a movies -ignore_collected true -ignore_watchlisted true                                                                                      
```
##### hide show recommendations:
```console
$ ./trakt-sync recommendations -a shows -i wellington-paranormal -hide
```
##### movies recommendations:
```console
$ ./trakt-sync recommendations -a shows                                                                                      
```
```console
$ ./trakt-sync recommendations -a shows -ignore_collected false -ignore_watchlisted false                                                                                     
```
#### scrobble:
##### scrobble start/pause/stop movie:
```console
$ ./trakt-sync scrobble -a start -t movie -progress 3.45 -i guardians-of-the-galaxy-2014
```
```console
$ ./trakt-sync scrobble -a pause -t movie -progress 3.45 -i guardians-of-the-galaxy-2014
```
```console
$ ./trakt-sync scrobble -a stop -t movie -progress 3.45 -i guardians-of-the-galaxy-2014
```
##### scrobble start/pause/stop episode:
```console
$ ./trakt-sync scrobble -a start -t episode -i 73629 -progress 10.25
```
```console
$ ./trakt-sync scrobble -a pause -t episode -i 73629 -progress 10.25
```
```console
$ ./trakt-sync scrobble -a stop -t episode -i 73629 -progress 50.25
```
##### scrobble start/pause/stop show by episode code (format season x episode):
```console
$ ./trakt-sync scrobble -a start -t show_episode -i 136121 -episode_code 1x5 -progress 3.45
```
```console
$ ./trakt-sync scrobble -a pause -t show_episode -i 136121 -episode_code 1x5 -progress 3.45
```
```console
$ ./trakt-sync scrobble -a stop -t show_episode -i 136121 -episode_code 1x5 -progress 3.45
```
##### scrobble start/pause/stop show by episode abs number (useful for Anime and Donghua):
```console
$ ./trakt-sync scrobble -a start -t show_episode -i 37696 -episode_abs 164 -progress 50
```
```console
$ ./trakt-sync scrobble -a pause -t show_episode -i 37696 -episode_abs 164 -progress 50
```
```console
$ ./trakt-sync scrobble -a stop -t show_episode -i 37696 -episode_abs 164 -progress 60
```
#### search:
##### Export search result by Text Query:
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
##### Export search result by Id lookup:
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
#### shows:
```console
$ ./trakt-sync shows -a trending
```
```console
$ ./trakt-sync shows -a popular
```
##### Get the most favorited shows
```console
$ ./trakt-sync shows -a favorited -period daily
```
```console
$ ./trakt-sync shows -a favorited -period weekly
```
```console
$ ./trakt-sync shows -a favorited -period monthly
```
```console
$ ./trakt-sync shows -a favorited -period all
```
##### Get the most played shows
```console
$ ./trakt-sync shows -a played -period daily
```
```console
$ ./trakt-sync shows -a played -period weekly
```
```console
$ ./trakt-sync shows -a played -period monthly
```
```console
$ ./trakt-sync shows -a played -period all
```
##### Get the most watched shows
```console
$ ./trakt-sync shows -a watched -period daily
```
```console
$ ./trakt-sync shows -a watched -period weekly
```
```console
$ ./trakt-sync shows -a watched -period monthly
```
```console
$ ./trakt-sync shows -a watched -period all
```
##### Get the most collected shows
```console
$ ./trakt-sync shows -a collected -period daily
```
```console
$ ./trakt-sync shows -a collected -period weekly
```
```console
$ ./trakt-sync shows -a collected -period monthly
```
```console
$ ./trakt-sync shows -a collected -period all
```
##### Get the most anticipated shows
```console
$ ./trakt-sync shows -a anticipated
```
##### Get recenty updated shows 
```console
$ ./trakt-sync shows -a updates -start_date YYYY-MM-DD
```
##### Get recenty updated show Trakt IDs 
```console
$ ./trakt-sync shows -a updated_ids -start_date YYYY-MM-DD
```
##### Get a show
```console
$ ./trakt-sync shows -a summary -i the-sopranos
```
##### Get all show aliases
```console
$ ./trakt-sync shows -a aliases -i the-sopranos
```
##### Get all show certifications
```console
$ ./trakt-sync shows -a certifications -i the-sopranos
```
##### Get all show translations
```console
$ ./trakt-sync shows -a translations -i the-sopranos -language es
```
##### Get all show comments
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s newest
```
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s oldest
```
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s likes
```
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s replies
```
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s highest
```
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s lowest
```
```console
$ ./trakt-sync shows -a comments -i the-sopranos -s plays
```
##### Get lists containing the show
```console
$ ./trakt-sync shows -a lists -i the-sopranos -t all -s popular
```
```console
$ ./trakt-sync shows -a lists -i the-sopranos -t all -s likes
```
```console
$ ./trakt-sync shows -a lists -i the-sopranos -t all -s comments
```
```console
$ ./trakt-sync shows -a lists -i the-sopranos -t all -s items
```
```console
$ ./trakt-sync shows -a lists -i the-sopranos -t all -s added
```
```console
$ ./trakt-sync shows -a lists -i the-sopranos -t all -s updated
```
##### Get show collection progress 
```console
$ ./trakt-sync shows -a collection_progress -i the-sopranos
$ ./trakt-sync shows -a collection_progress -i the-sopranos -hidden false 
$ ./trakt-sync shows -a collection_progress -i the-sopranos -specials false 
$ ./trakt-sync shows -a collection_progress -i the-sopranos -count_specials true
$ ./trakt-sync shows -a collection_progress -i the-sopranos -hidden true -specials true -count_specials true
```
##### Get show watched progress 
```console
$ ./trakt-sync shows -a watched_progress -i the-sopranos
$ ./trakt-sync shows -a watched_progress -i the-sopranos -hidden false 
$ ./trakt-sync shows -a watched_progress -i the-sopranos -specials false 
$ ./trakt-sync shows -a watched_progress -i the-sopranos -count_specials true
$ ./trakt-sync shows -a watched_progress -i the-sopranos -hidden true -specials true -count_specials true
```
##### Reset show progress
```console
$ ./trakt-sync shows -a reset_show_progress -i the-sopranos
```
##### Undo Reset show progress
```console
$ ./trakt-sync shows -a reset_show_progress -i the-sopranos -undo
```
##### Get all people for a show
```console
$ ./trakt-sync shows -a people -i the-sopranos
$ ./trakt-sync shows -a people -i the-sopranos -ex guest_stars
```
##### Get show ratings
```console
$ ./trakt-sync shows -a ratings -i the-sopranos
```
##### Get related shows
```console
$ ./trakt-sync shows -a related -i the-sopranos
```
##### Get show studios
```console
$ ./trakt-sync shows -a studios -i the-sopranos
```
##### Get users watching right now
```console
$ ./trakt-sync shows -a watching -i the-sopranos
```
##### Get next episode
```console
$ ./trakt-sync shows -a next_episode -i the-sopranos
```
##### Get last episode
```console
$ ./trakt-sync shows -a last_episode -i the-sopranos
```
##### Get all videos
```console
$ ./trakt-sync shows -a videos -i the-sopranos
```
##### Refresh show metadata
```console
$ ./trakt-sync shows -a refresh -i the-sopranos
```
#### seasons:
##### Get all seasons for a show
```console
$ ./trakt-sync seasons -a summary -i the-sopranos -ex full
```
##### Get single seasons for a show
```console
$ ./trakt-sync seasons -a season -i the-sopranos -season 1 -ex full
```
##### Get all episodes for a single season
```console
$ ./trakt-sync seasons -a episodes -i the-sopranos -season 1 -translations es -ex full
```
##### Get all season translations - all languages
```console
$ ./trakt-sync seasons -a translations -i the-sopranos -season 1 
```
##### Get all season translations - selected language
```console
$ ./trakt-sync seasons -a translations -i the-sopranos -season 1 -language en
```
##### Get all season comments
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s newest
```
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s oldest
```
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s likes
```
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s replies
```
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s highest
```
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s lowest
```
```console
$ ./trakt-sync seasons -a comments -i the-sopranos -season 1 -s plays
```
##### Get lists containing this season
```console
$ ./trakt-sync seasons -a lists -i the-sopranos -season 1 -t all -s popular
```
```console
$ ./trakt-sync seasons -a lists -i the-sopranos -season 1 -t all -s likes
```
```console
$ ./trakt-sync seasons -a lists -i the-sopranos -season 1 -t all -s comments
```
```console
$ ./trakt-sync seasons -a lists -i the-sopranos -season 1 -t all -s items
```
```console
$ ./trakt-sync seasons -a lists -i the-sopranos -season 1 -t all -s added
```
```console
$ ./trakt-sync seasons -a lists -i the-sopranos -season 1 -t all -s updated
```
##### Get all people for season
```console
$ ./trakt-sync seasons -a people -i the-sopranos -season 1
```
##### Get season ratings
```console
$ ./trakt-sync seasons -a ratings -i the-sopranos -season 1
```
##### Get related seasons
```console
$ ./trakt-sync seasons -a related -i the-sopranos -season 1
```
##### Get seasons stats
```console
$ ./trakt-sync seasons -a stats -i the-sopranos -season 1
```
##### Get users watching right now
```console
$ ./trakt-sync seasons -a watching -i the-sopranos -season 1
```
##### Get all videos
```console
$ ./trakt-sync seasons -a videos -i the-sopranos -season 1
```

#### episodes:
##### Get a single episode for a show
```console
$ ./trakt-sync episodes -a summary -i the-sopranos -season 1 -episode 1 -ex full
```
##### Get all episode translations - all languages
```console
$ ./trakt-sync episodes -a translations -i the-sopranos -season 1 -episode 1 
```
##### Get all episode translations - selected language
```console
$ ./trakt-sync episodes -a translations -i the-sopranos -season 1 -episode 1 -language en
```
##### Get all episode comments
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s newest
```
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s oldest
```
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s likes
```
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s replies
```
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s highest
```
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s lowest
```
```console
$ ./trakt-sync episodes -a comments -i the-sopranos -season 1 -episode 1 -s plays
```
##### Get lists containing this episode
```console
$ ./trakt-sync episodes -a lists -i the-sopranos -season 1 -episode 1 -t all -s popular
```
```console
$ ./trakt-sync episodes -a lists -i the-sopranos -season 1 -episode 1 -t all -s likes
```
```console
$ ./trakt-sync episodes -a lists -i the-sopranos -season 1 -episode 1 -t all -s comments
```
```console
$ ./trakt-sync episodes -a lists -i the-sopranos -season 1 -episode 1 -t all -s items
```
```console
$ ./trakt-sync episodes -a lists -i the-sopranos -season 1 -episode 1 -t all -s added
```
```console
$ ./trakt-sync episodes -a lists -i the-sopranos -season 1 -episode 1 -t all -s updated
```
##### Get all people for episode
```console
$ ./trakt-sync episodes -a people -i the-sopranos -season 1 -episode 1
```
##### Get episode ratings
```console
$ ./trakt-sync episodes -a ratings -i the-sopranos -season 1 -episode 1
```
##### Get related episodes
```console
$ ./trakt-sync episodes -a related -i the-sopranos -season 1 -episode 1
```
##### Get episodes stats
```console
$ ./trakt-sync episodes -a stats -i the-sopranos -season 1 -episode 1
```
##### Get users watching right now
```console
$ ./trakt-sync episodes -a watching -i the-sopranos -season 1 -episode 1
```
##### Get all episodes videos
```console
$ ./trakt-sync episodes -a videos -i the-sopranos -season 1 -episode 1
```

#### users:

##### Export movies or shows or episodes from user lists:
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t episodes
```
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t shows
```
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t movies
```
##### Fetch lists for selected user:
```console
$ ./trakt-sync users -a lists -u username 
```
##### Fetch saved filters for selected user:
```console
$ ./trakt-sync users -a saved_filters -u username 
```
##### Fetch stats for selected user:
```console
$ ./trakt-sync users -a stats -u sean 
```
##### Fetch settings for current user:
```console
$ ./trakt-sync users -a settings 
```
##### Fetch watched movies for selected user:
```console
$ ./trakt-sync users -a watched -t movies -u sean 
```
##### Fetch watched shows for selected user:
```console
$ ./trakt-sync users -a watched -t shows -u sean 
```
##### Fetch watched shows for selected user without seasons:
```console
$ ./trakt-sync users -a watched -t shows -u sean --ex noseasons
```
#### watchlist

##### Export all movies from watchlist:
```console
$ ./trakt-sync watchlist -t movies -f tmdb -> export_watchlist_movies_tmdb.json 
```
```console
$ ./trakt-sync watchlist -t movies -f imdb -> export_watchlist_movies_imdb.json
```

## License

[MIT](./LICENSE)

