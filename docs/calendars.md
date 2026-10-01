#### Calendars:
The old hyphenated action names (`my-shows`, `all-new-shows`, `hot-releases`, ...) still work, but are deprecated: use the underscore names (`my_shows`, `all_new_shows`, `hot_releases`). Every `all_*` action also has a `my_*` variant for your own calendar.
```console
$ ./trakt-sync calendars -a all_shows -> export_calendars_shows_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_new_shows -> export_calendars_new_shows_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_season_premieres -> export_calendars_season_premieres_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_finales -> export_calendars_finales_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_movies -> export_calendars_movies_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_dvd -> export_calendars_dvd_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_media -> export_calendars_media_20240707_7.json
```
```console
$ ./trakt-sync calendars -a all_streaming -> export_calendars_streaming_20240707_7.json
```
```console
$ ./trakt-sync calendars -a hot_releases -> export_calendars_hot_releases_20240707_7.json
```
```console
$ ./trakt-sync calendars -a hot_premieres -> export_calendars_hot_premieres_20240707_7.json
```
```console
$ ./trakt-sync calendars -a hot_new_shows -> export_calendars_hot_new_shows_20240707_7.json
```
```console
$ ./trakt-sync calendars -a hot_finales -> export_calendars_hot_finales_20240707_7.json
```
