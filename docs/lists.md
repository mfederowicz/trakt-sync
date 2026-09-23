#### Lists:
```console
$ ./trakt-sync lists -a trending
```
```console
$ ./trakt-sync lists -a popular
```
`-t` on `trending`/`popular` is sent to the API as-is; the API contract does not list the allowed types.
```console
$ ./trakt-sync lists -a trending -t personal -> export_lists_trending_personal.json
```
```console
$ ./trakt-sync lists -a popular -t official -> export_lists_popular_official.json
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
```console
$ ./trakt-sync lists -a items -trakt_id 2142753 -t movie -sort_by added -sort_how desc
```
Without `-t`, `items` returns all item types (`movie,show,episode,season`).
-- (temp not working - problems with api endpoint)
```console
$ ./trakt-sync lists -a comments -trakt_id 2142753
```
```console
$ ./trakt-sync lists -a report -trakt_id 2142753 -r spam -message "only ads"
```
