#### Media:
Movies and shows together in one list. Paging follows `per_page` and `pages_limit` from the config file; `-ex` sets extended info.
```console
$ ./trakt-sync media -a trending -> export_media_trending.json
```
```console
$ ./trakt-sync media -a popular -> export_media_popular.json
```
```console
$ ./trakt-sync media -a anticipated -> export_media_anticipated.json
```
All three actions take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`,
`-certifications` (comma separated where a filter takes several values), `-start_date`, `-end_date` and
`-watchnow favorites|any|any_all|free|free_all|subscriptions|subscriptions_all`.
```console
$ ./trakt-sync media -a trending -genres action,drama -years 2020-2026 -ratings 75-100
```
