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
`-languages` (language codes, comma separated) filters them too.
`-parental_nudity`, `-parental_violence`, `-parental_profanity`, `-parental_alcohol` and `-parental_frightening` filter by parental guide: each takes a severity range `min-max` from 0 (none) to 3 (severe), ie: `-parental_violence 0-1`. Titles without a parental guide are left out, unless you add `-parental_include_unrated`.
`-theme <name>` applies a seasonal theme such as `halloween` or `christmas`: Trakt turns it into filters for the list, a filter you set yourself wins over it, and an unknown theme is ignored.
```console
$ ./trakt-sync media -a trending -theme halloween
```
```console
$ ./trakt-sync media -a trending -genres action,drama -years 2020-2026 -ratings 75-100
```
##### Anticipated recommendations
Needs your login (OAuth). Returns unreleased movies and shows ranked for your taste, at most 100 in total, each with a `score` and its `released` date (`YYYY-MM-DD`, left out when Trakt does not know it yet). Takes no filters or `-ex`; paging follows `per_page` and `pages_limit`.
```console
$ ./trakt-sync media -a anticipated_recommendations -> export_media_anticipated_recommendations.json
```
