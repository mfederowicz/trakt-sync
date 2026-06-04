#### Movies:
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

