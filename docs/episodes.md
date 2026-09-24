#### Episodes:
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
##### Get where to watch an episode
```console
$ ./trakt-sync episodes -a watchnow -i the-sopranos -season 1 -episode 2 -country us -> export_episodes_watchnow_the-sopranos.json
```
```console
$ ./trakt-sync -ex streaming_ranks episodes -a watchnow -i the-sopranos -season 1 -episode 2 -country us -links tvos,direct
```
Or by the episode's own Trakt ID (`-i` without `-season` and `-episode`):
```console
$ ./trakt-sync episodes -a watchnow -i 73482 -country us -> export_episodes_watchnow_73482.json
```
`-country` is required; `-links` and `-ex streaming_ranks` work as in `shows -a watchnow`. Marked Limited Access by Trakt.
##### Report an episode
```console
$ ./trakt-sync episodes -a report -i the-sopranos -season 1 -episode 2 -r runtime -message "runtime is 50 min"
```
Or report an episode by its own Trakt ID: pass it in `-i` without `-season` and `-episode`:
```console
$ ./trakt-sync episodes -a report -i 73482 -r runtime
```
With `-season` and `-episode`, `-i` is the show; without both, `-i` is the episode's Trakt ID. Only one of them is an error. `-r` is one of: `duplicate`, `remove`, `data_refresh`, `metadata`, `adult`, `runtime`, `language`, `spam`, `tmdb`, `other`.
