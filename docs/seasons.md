#### Seasons:
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
