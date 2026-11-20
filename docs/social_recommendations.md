#### Social recommendations (locked):
Trakt dropped the social recommendations endpoints from its API contract on 2026-10-13, and the live API answers 404
since 2026-11-20. Both actions are locked: they send no request and stop with a notice. [`recommendations`](./recommendations.md) is not affected.
```console
$ ./trakt-sync social_recommendations -a movies
social_recommendations -a movies is locked since 2026-11-20: Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent
```
The rest of this page describes what the actions did while Trakt served them.

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
