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
