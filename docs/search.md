#### Search:
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
