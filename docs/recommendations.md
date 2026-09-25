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
