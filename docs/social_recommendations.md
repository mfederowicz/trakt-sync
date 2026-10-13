#### Social recommendations (deprecated):
Deprecated since 2026-10-13: Trakt dropped the social recommendations endpoints from its API contract. Both actions still work (checked 2026-10-13) and print a deprecation notice; they may stop working. [`recommendations`](./recommendations.md) is not affected.

Movie and show recommendations based on the people you follow. `limit` follows `per_page` from the config file; `-ex` sets extended info (`full`, `images`, `colors`, `streaming_ids`).
```console
$ ./trakt-sync social_recommendations -a movies -> export_social_recommendations_movies.json
```
```console
$ ./trakt-sync social_recommendations -a shows -> export_social_recommendations_shows.json
```
Skip items you already watched, collected or watchlisted:
```console
$ ./trakt-sync social_recommendations -a movies -ignore_watched true -ignore_collected true -ignore_watchlisted true
```
Only recommendations from the last 30 days of activity:
```console
$ ./trakt-sync social_recommendations -a shows -watch_window 30
```
