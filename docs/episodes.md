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
