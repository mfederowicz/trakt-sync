#### Smart lists:
Smart lists are saved filters that resolve to movies or shows. `-i` takes the smart list slug. Only public smart lists return data, unless you own the list.

##### Smart list definition (name, media type, filters):
```console
$ ./trakt-sync smart_lists -a summary -i top-sci-fi -> export_smart_lists_summary_top-sci-fi.json
```
##### Smart list items:
Paging follows `per_page` and `pages_limit` from the config file; `-ex` sets extended info (`full`, `images`, `colors`, `streaming_ids`).
```console
$ ./trakt-sync smart_lists -a items -i top-sci-fi -> export_smart_lists_items_top-sci-fi.json
```
Refine the items with media filters (`-watchnow favorites|any|any_all|free|free_all|subscriptions|subscriptions_all`,
`-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`) and skip what you already watched or watchlisted:
```console
$ ./trakt-sync smart_lists -a items -i top-sci-fi -genres action -years 2020-2026 -ratings 75-100 -ignore_watched true
```
```console
$ ./trakt-sync smart_lists -a items -i top-sci-fi -watchnow subscriptions -ignore_watchlisted true
```
