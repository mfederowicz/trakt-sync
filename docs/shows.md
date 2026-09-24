#### Shows:
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
```
```console
$ ./trakt-sync shows -a collection_progress -i the-sopranos -hidden false
```
```console
$ ./trakt-sync shows -a collection_progress -i the-sopranos -specials false
```
```console
$ ./trakt-sync shows -a collection_progress -i the-sopranos -count_specials true
```
```console
$ ./trakt-sync shows -a collection_progress -i the-sopranos -hidden true -specials true -count_specials true
```
##### Get show watched progress
```console
$ ./trakt-sync shows -a watched_progress -i the-sopranos
```
```console
$ ./trakt-sync shows -a watched_progress -i the-sopranos -hidden false
```
```console
$ ./trakt-sync shows -a watched_progress -i the-sopranos -specials false
```
```console
$ ./trakt-sync shows -a watched_progress -i the-sopranos -count_specials true
```
```console
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
```
```console
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
##### Get show sentiments
```console
$ ./trakt-sync shows -a sentiments -i the-sopranos -> export_shows_sentiments_the-sopranos.json
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
##### Get where to watch
```console
$ ./trakt-sync shows -a watchnow -i the-sopranos -country us -> export_shows_watchnow_the-sopranos.json
```
```console
$ ./trakt-sync -ex streaming_ranks shows -a watchnow -i the-sopranos -country us -links tvos,direct,android,webos
```
`-links` adds provider links (any of `tvos`, `direct`, `android`, `webos`); `-ex streaming_ranks` adds the JustWatch rank.
```console
$ ./trakt-sync shows -a justwatch_links -i the-sopranos -country pl -> export_shows_justwatch_links_the-sopranos.json
```
`watchnow` and `justwatch_links` need `-country` (2 character code) and are marked Limited Access by Trakt;
an API app without access gets a limited access error.
##### Get all videos
```console
$ ./trakt-sync shows -a videos -i the-sopranos
```
##### Refresh show metadata
```console
$ ./trakt-sync shows -a refresh -i the-sopranos
```
##### Refresh show JustWatch links
```console
$ ./trakt-sync shows -a refresh_justwatch -i the-sopranos
```
`refresh_justwatch` is marked VIP only by Trakt; when Trakt answers with 426 the upgrade page opens in the browser.
The live API accepted it from a non-VIP account (checked 2026-09-24), unlike `movies -a refresh_justwatch`.
##### Report a show
```console
$ ./trakt-sync shows -a report -i the-sopranos -r metadata -message "overview is wrong"
```
`-r` is one of: `duplicate`, `remove`, `data_refresh`, `metadata`, `adult`, `runtime`, `language`, `spam`, `tmdb`, `other`.
