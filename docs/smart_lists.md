#### Smart lists:
Smart lists are saved filters that resolve to movies or shows. `-i` takes the smart list slug (`ids.slug` in the summary). Only public smart lists return data, unless you own the list; an unknown or private slug ends with "not found smart list for:<slug>".

##### Smart list definition (name, media type, filters):
```console
$ ./trakt-sync smart_lists -a summary -i <slug> -> export_smart_lists_summary_<slug>.json
```
##### Smart list items:
Paging follows `per_page` and `pages_limit` from the config file; `-ex` sets extended info (`full`, `images`, `colors`, `streaming_ids`).
```console
$ ./trakt-sync smart_lists -a items -i <slug> -> export_smart_lists_items_<slug>.json
```
Refine the items with media filters (`-watchnow favorites|any|any_all|free|free_all|subscriptions|subscriptions_all`,
`-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`) and skip what you already watched or watchlisted:
```console
$ ./trakt-sync smart_lists -a items -i <slug> -genres action -years 2020-2026 -ratings 75-100 -ignore_watched true
```
```console
$ ./trakt-sync smart_lists -a items -i <slug> -watchnow subscriptions -ignore_watchlisted true
```
Send the `-watchnow` filter for another region with `-watchnow_country <xx>` (two lowercase letters). Without it the API uses the list region, then the owner's region, then `us`:
```console
$ ./trakt-sync smart_lists -a items -i <slug> -watchnow free -watchnow_country pl
```
Filter by parental guide with `-parental_nudity`, `-parental_violence`, `-parental_profanity`, `-parental_alcohol` and `-parental_frightening`.
Each takes a severity range `min-max` from 0 (none) to 3 (severe). Titles without a parental guide are left out, unless you add `-parental_include_unrated`:
```console
$ ./trakt-sync smart_lists -a items -i <slug> -parental_violence 0-1 -parental_frightening 0-2 -parental_include_unrated
```
##### Smart list items by type and sort:
`-t all|movies|shows` narrows a list that holds both movies and shows, `-sort_by` and `-sort_how asc|desc` reorder it. Sort values: `rank` (the list's own order),
`random`, `title`, `released`, `runtime`, `percentage`, `votes`, `imdb_rating`, `imdb_votes`, `tmdb_rating`, `tmdb_votes`, `rt_tomatometer`, `rt_audience`, `metascore`,
and `added` on watchlist lists. Give at least one of the three flags after the module name; the others default to `-t all`, `-sort_by rank`, and `-sort_how asc` for `rank` and `title`, `desc` otherwise.
The filters above work here too, and the output file is the same.
```console
$ ./trakt-sync smart_lists -a items -i <slug> -t movies -sort_by imdb_rating
```
```console
$ ./trakt-sync smart_lists -a items -i <slug> -sort_by released -sort_how asc
```
