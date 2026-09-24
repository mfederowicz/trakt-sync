#### Search:
The old action names `text-query` and `id-lookup` still work, but are deprecated: use `text_query` and `id_lookup`.

##### Export search result by Text Query:
Every `--field` value must be valid for every `-t` type given.

```console
$  ./trakt-sync search -a text_query -t movie -q freddy --field title
```
```console
$  ./trakt-sync search -a text_query -t movie -q freddy --field tagline
```
```console
$  ./trakt-sync search -a text_query -t movie -t show -q freddy --field overview
```
```console
$  ./trakt-sync search -a text_query -t movie -t show -t episode -q freddy --field title
```
```console
$  ./trakt-sync search -a text_query -t person -t list -q freddy --field name
```
```console
$  ./trakt-sync search -a text_query -t movie -t show -q freddy --field original_title
```
```console
$  ./trakt-sync search -a text_query -t episode -q freddy --field show_title
```
##### Export search result by Id lookup:
```console
$ ./trakt-sync search -a id_lookup -i 12601 -t movie -t show
```
```console
$ ./trakt-sync search -a id_lookup --id_type tvdb -i 12601 -t movie -t show
```
```console
$ ./trakt-sync search -a id_lookup --id_type imdb -i 12601 -t movie
```
```console
$ ./trakt-sync search -a id_lookup --id_type imdb -i 12601 -t podcast
```
```console
$ ./trakt-sync search -a id_lookup --id_type imdb -i tt0266697
```
```console
$ ./trakt-sync search -a id_lookup --id_type tvdb -i 75725
```
```console
$ ./trakt-sync search -a id_lookup --id_type tvdb -i 75725 -t podcast
```
```console
$ ./trakt-sync search -a id_lookup -i 75725
```
```console
$ ./trakt-sync search -a id_lookup -i 75725 -t episode
```
```console
$ ./trakt-sync search -a id_lookup --id_type tmdb -i 254265
```
##### Export exact text query results:
`-t` takes one value: `movie` or `show`; `-q` is required.
```console
$ ./trakt-sync search -a exact_query -t movie -q "tron legacy"
```
```console
$ ./trakt-sync -ex full search -a exact_query -t show -q dark
```
##### Export trending search results:
Globally trending searches: the items people most often picked from Trakt search results, with `count` as the number
of picks. It is not a list of trending movies, shows or people. `-t` takes one value: `movies`, `shows` or `people`.
`-q` filters on the text people typed into search, not on the item's name, so `-t shows -q reacher` finds Reacher,
but `-t people -q ritchson` is empty unless people searched for that name.
```console
$ ./trakt-sync search -a trending -t movies
```
```console
$ ./trakt-sync search -a trending -t people -q keanu
```
##### Add or remove a recent search:
Records (or removes) that the item with Trakt ID `-i` was picked from search results for the text `-q`. This changes
the global search trends that `search -a trending` returns for every user, not a personal search history. `-t` takes
one value: `movies`, `shows`, `people` or `lists`; `-i` must be the numeric Trakt ID. Nothing is written to a file.
```console
$ ./trakt-sync search -a add_recent -t shows -q reacher -i 139606
```
```console
$ ./trakt-sync search -a remove_recent -t shows -q reacher -i 139606
```
