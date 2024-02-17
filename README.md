Sample export usage

Export all movies from watchlist:

$ ./trakt-sync watchlist -t movies -f tmdb -> export_watchlist_movies_tmdb.json 
$ ./trakt-sync watchlist -t movies -f imdb -> export_watchlist_movies_imdb.json

Export all tvshows from the watching history:

$ ./trakt-sync history -t shows -> export_history_shows_imdb.json 

Export all episodes from the watching history:

$ ./trakt-sync history -t episodes -f tmdb -> export_history_episodes_tmdb.json
$ ./trakt-sync history -t episodes -f imdb -> export_history_episodes_imdb.json

Export all or my calendars:

$ ./trakt-sync calendars -a all-shows -> export_calendars_shows_20240707_7.json
$ ./trakt-sync calendars -a all-new-shows -> export_calendars_new_shows_20240707_7.json
$ ./trakt-sync calendars -a all-season-premieres -> export_calendars_season_premieres_20240707_7.json
$ ./trakt-sync calendars -a all-finales -> export_calendars_finales_20240707_7.json 
$ ./trakt-sync calendars -a all-movies -> export_calendars_movies_20240707_7.json  
$ ./trakt-sync calendars -a all-dvd -> export_calendars_dvd_20240707_7.json  

Export search result by Text Query:

$  ./trakt-sync search -a text-query -t movie -q freddy --field title
$  ./trakt-sync search -a text-query -t movie -t show -q freddy --field tagline
$  ./trakt-sync search -a text-query -t movie -t show -t list -q freddy --field name
$  ./trakt-sync search -a text-query -t movie -t show -t list -q freddy --field title
$  ./trakt-sync search -a text-query -t person -t list -q freddy --field name
$  ./trakt-sync search -a text-query -t movie -t show -t list -q freddy --field title

Export search result by Id lookup:

$ ./trakt-sync search -a id-lookup -i 12601 -t movie -t show
$ ./trakt-sync search -a id-lookup --id_type tvdb -i 12601 -t movie -t show
$ ./trakt-sync search -a id-lookup --id_type imdb -i 12601 -t movie
$ ./trakt-sync search -a id-lookup --id_type imdb -i 12601 -t podcast
$ ./trakt-sync search -a id-lookup --id_type imdb -i tt0266697
$ ./trakt-sync search -a id-lookup --id_type tvdb -i 75725
$ ./trakt-sync search -a id-lookup --id_type tvdb -i 75725 -t podcast
$ ./trakt-sync search -a id-lookup -i 75725 
$ ./trakt-sync search -a id-lookup -i 75725 -t episode
$ ./trakt-sync search -a id-lookup --id_type tmdb -i 254265
