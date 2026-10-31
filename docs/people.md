#### People:
```console
$ ./trakt-sync people -a updates -start_date 2024-10-13
```
```console
$ ./trakt-sync people -a updated_ids -start_date 2024-10-13
```
```console
$ ./trakt-sync people -a summary -i john-wayne
```
The export holds the `biography` and a short `summary` (two or three sentences) as plain text. The API contract lists both under `-ex full`, but on 2026-10-14 the live API sent them without `-ex` too. `-ex full,markdown` returns both with markdown links to the movies and shows they mention:
```console
$ ./trakt-sync people -a summary -i john-wayne -ex full,markdown
```
```console
$ ./trakt-sync people -a movies -i john-wayne
```
```console
$ ./trakt-sync people -a shows -i john-wayne
```
```console
$ ./trakt-sync people -a lists -i john-wayne
```
```console
$ ./trakt-sync people -a lists -i john-wayne -t official -s likes
```
`-t` is one of: `all`, `personal` (default), `official`, `watchlist`, `favorites`.
`-s` is one of: `popular` (default), `likes`, `comments`, `items`, `added`, `updated`.

The Trakt API currently returns no lists for any person on this route (checked 2026-10-01, again 2026-10-31), so the command ends with `empty lists`.
Deprecated since 2026-10-20: the Trakt API contract marks the endpoint deprecated because lists no longer hold people. `people -a lists` still sends its request and prints a deprecation notice first; it may stop working.
##### Refresh person metadata (deprecated)
```console
$ ./trakt-sync people -a refresh -i john-wayne
```
Deprecated since 2026-10-13: Trakt dropped this endpoint from its API contract. `people -a refresh` still works (checked 2026-10-13) and prints a deprecation notice; it may stop working.
##### Report a person
```console
$ ./trakt-sync people -a report -i john-wayne -r metadata -message "birthday is wrong"
```
`-r` is one of: `duplicate`, `remove`, `data_refresh`, `metadata`, `adult`, `language`, `spam`, `tmdb`, `other`.
