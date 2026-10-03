#### Recommendations:
##### Hide movie recommendations:
```console
$ ./trakt-sync recommendations -a movies -i black-bag-2025 -hide
```
##### Movies recommendations:
```console
$ ./trakt-sync recommendations -a movies
```
```console
$ ./trakt-sync recommendations -a movies -ignore_collected true -ignore_watchlisted true
```
Skip movies you already watched, based on the last 30 days of activity:
```console
$ ./trakt-sync recommendations -a movies -ignore_watched true -watch_window 30
```
##### Hide show recommendations:
```console
$ ./trakt-sync recommendations -a shows -i wellington-paranormal -hide
```
##### Shows recommendations:
```console
$ ./trakt-sync recommendations -a shows
```
```console
$ ./trakt-sync recommendations -a shows -ignore_collected false -ignore_watchlisted false
```
```console
$ ./trakt-sync recommendations -a shows -ignore_watched true -watch_window 30
```
##### Filter the recommendations:
`-a movies` and `-a shows` take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`,
`-certifications` (comma separated where a filter takes several values), `-start_date`, `-end_date` and
`-watchnow favorites|any|any_all|free|free_all|subscriptions|subscriptions_all`.
`-languages` (language codes, comma separated) filters them too.
```console
$ ./trakt-sync recommendations -a movies -genres horror -years 2020-2026
```
```console
$ ./trakt-sync recommendations -a shows -watchnow subscriptions -ratings 75-100
```
