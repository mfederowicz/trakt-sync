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
```console
$ ./trakt-sync people -a movies -i john-wayne
```
```console
$ ./trakt-sync people -a shows -i john-wayne
```
```console
$ ./trakt-sync people -a lists -i john-wayne
```
##### Report a person
```console
$ ./trakt-sync people -a report -i john-wayne -r metadata -message "birthday is wrong"
```
`-r` is one of: `duplicate`, `remove`, `data_refresh`, `metadata`, `adult`, `language`, `spam`, `tmdb`, `other`.
