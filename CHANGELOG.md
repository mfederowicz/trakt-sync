# Changelog

All notable changes to this project are documented here, in the
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format.

Versioning follows [SemVer](https://semver.org/):

- **Major** (`X.0.0`) - breaking changes to the CLI surface: a removed or renamed module,
  action or flag, or a changed output file name or JSON shape that scripts may depend on.
- **Minor** (`1.X.0`) - a module or action lands (new endpoints/commands), or other
  user-facing behavior changes in a compatible way.
- **Patch** (`1.15.X`) - bug fixes, CI/tooling changes, docs-only changes, refactors with no
  behavior change.

A version is tagged once a module (or a meaningful fix) is done - there's no fixed release
schedule.

## [Unreleased]

### Added

- `sync -a playback -t all`, and `-ex` for `sync -a playback`, which was not sent before.
- `sync -a get_minimal_collection -t movies|shows|episodes [-available_on plex]`: the collection as a compact map of
  Trakt IDs to collected dates (shows nested by season and episode), for syncing local state.
- `sync -a get_up_next [-include_stats] [-lifetime_stats]`: shows you are watching, with their next episode and
  progress.
- `sync -a get_watched_progress [-hide_completed | -hide_not_completed] [-only_rewatching] [-lifetime_stats]`:
  watched progress of your shows. Both actions take `-sort_by` / `-sort_how` (sent only when given) and `-ex`.
- `sync -a get_up_next_nitro [-intent all|continue|start|completed] [-watchnow <filter>]`: up next with media filters
  (`-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`, `-start_date`,
  `-end_date`). The global `-genres`, `-years`, `-countries` and `-runtimes` flags were accepted before but not used.

### Changed

- `CHANGELOG.md` now covers every release back to 1.0.0.
- GitHub release notes now show the version's `CHANGELOG.md` section (Added / Changed / Fixed) instead of a list of
  commits.

### Fixed

- A crash inside any command was caught but then reported nothing, as if the command had succeeded; it now
  prints `panic error:<reason>` (or `fatal error`).
- `checkin -a movie -trakt_id <id>` and `checkin -a show_episode -trakt_id <id>` looked the item up without its ID
  (`GET /movies/`, `GET /shows/`) since 1.8.0, so the checkin could not work; they now use `-trakt_id`.
- `checkin -a show_episode`: an active checkin (409) without an expiry time in the response now reports the
  existing checkin instead of crashing.
- `users -a history -start_at <date> -end_at <date>` failed with `flag provided but not defined: -start_at` and
  silently used the default window; `users` now accepts both flags.
- `sync` and `users` ignored `-o` and always wrote to their generated file name (e.g.
  `export_sync_history_movies.json`); `-o` now sets the output file.
- `sync -a playback` without `-t` returned only movies (the default type), although the docs describe it as all
  playback. It now returns movies and episodes; `-t movies|episodes` still narrows it.

## [1.17.0] - 2026-09-25

### Added

- `movies -a watchnow -i <id> -country <code>`: watch now sources (streaming, rent, purchase) for a movie in one country;
  `-links tvos,direct,android,webos` adds provider links and `-ex streaming_ranks` adds the JustWatch rank.
- `movies -a justwatch_links -i <id> -country <code>`: the JustWatch link for a movie in one country.
- `shows -a report -i <id> -r <reason> [-message "..."]`: report a show for moderator review.
- `shows -a sentiments -i <id>`: good and bad sentiments from a show's comments and reactions.
- `shows -a watchnow -i <id> -country <code>`: watch now sources for a show in one country, with the same
  `-links` and `-ex streaming_ranks` options as `movies -a watchnow`.
- `shows -a justwatch_links -i <id> -country <code>`: the JustWatch link for a show in one country.
- `seasons -a report -i <show> -season <n> -r <reason> [-message "..."]`: report a season for moderator review
  (`-season 0` is specials).
- `episodes -a report -i <show> -season <n> -episode <n> -r <reason> [-message "..."]`: report an episode for
  moderator review.
- `episodes -a watchnow -i <show> -season <n> -episode <n> -country <code>`: watch now sources for an episode, with the
  same `-links` and `-ex streaming_ranks` options as `shows -a watchnow`.
- `seasons -a justwatch_links -i <show> -season <n> -country <code>`: the JustWatch link for a season.
- `shows -a refresh_justwatch -i <id>`: queue a refresh of the show's JustWatch links (VIP only).
- `people -a report -i <id> -r <reason> [-message "..."]`: report a person for moderator review.
- `seasons -a report`, `episodes -a report` and `episodes -a watchnow` accept the season's or episode's own Trakt ID
  in `-i` when `-season` (and for episodes `-episode`) is left out; with them, `-i` is still the show.
- `search -a text_query --field`: accepts `original_title` for `-t movie` / `-t show` and `show_title` for
  `-t episode`, the remaining search fields from the Trakt API docs.
- `search -a exact_query -t movie|show -q <query>`: exact title matches for movies or shows.
- `search -a trending -t movies|shows|people [-q <query>]`: globally trending searches, the items people most often
  picked from search results; `-q` filters on the typed search text.
- `search -a add_recent|remove_recent -t movies|shows|people|lists -q <query> -i <trakt id>`: add or remove a pick
  in the global search trends that `search -a trending` returns.

### Changed

- `search -a text_query` and `search -a id_lookup` replace `text-query` and `id-lookup`, to match the underscore action
  names of the other modules. The old names still work and print a deprecation note.
- `search` with an unknown or missing `-a` prints the same "Available actions" list as the other modules.

### Removed

- `search -t podcast` and `-t podcast_episode`, and the `podcast` / `podcast_episode` fields in search results: the
  Trakt API contract has no podcast search types or podcast results.

### Fixed

- `search -a text_query --field <name>`: the field filter was sent as `field` instead of `fields`, so the API
  ignored it and searched its default fields. It is now sent as `fields`.
- `comments` on a season or episode always failed with `set traktId`, even with `-trakt_id`/`-i` given: the check
  read a field the `comments` module never sets. It now checks the id that is actually passed.
- `checkin -a episode -trakt_id <id>` crashed with a nil pointer: the episode was looked up by an empty id instead of
  `-trakt_id`. It now sends the `-trakt_id` directly.
- `notes -t season|episode`, `notes -t rating|collection -item season|episode`, `comments` on seasons and episodes, and
  `scrobble -t episode` no longer look the item up with `GET seasons/{id}` / `GET episodes/{id}`, which are not in the
  Trakt API docs. The request carries only the numeric Trakt ID from `-i`; a non-numeric id now fails before anything
  is posted (notes used to post without the item).
- `movies -a sentiments` and `shows -a sentiments`: an unknown id wrote an empty `{}` file, because the API answers
  it with an empty object instead of 404; it now fails with `no sentiments for:<id>` and writes nothing.
- `movies -a hot` and `movies -a streaming`: the live Trakt API returns 404 for these routes; the error now says the
  route is documented but not served, instead of a bare 404.

## [1.16.0] - 2026-09-24

### Added

- `calendars -a {my,all}-media`: movies and shows releasing in the selected period.
- `calendars -a {my,all}-streaming`: movies with a streaming release in the selected period.
- `calendars -a hot-releases|hot-premieres|hot-new-shows|hot-finales`: hot releases, premieres,
  new shows and finales in the selected period.
- New `media` module: `media -a trending|popular|anticipated` exports movies and shows together in one list.
- `comments -a reactions` and `comments -a reactions_summary`: export the reactions on a comment, or their totals by type.
- `comments -a reaction -reaction <type>`: add a reaction to a comment; add `-remove` to take it back.
- `comments -a report -r <reason> [-message "..."]`: report a comment for moderator review.
- `lists -a trending|popular -t <type>`: lists of one type. The value is sent to the API as-is; the API
  contract does not list the allowed types.
- `lists -a items -sort_by <field> -sort_how asc|desc`: sort list items.
- `lists -a report -trakt_id <id> -r <reason> [-message "..."]`: report a list for moderator review.
- `movies -a hot`: hot movies based on current list activity.
- `movies -a streaming -period daily|weekly|monthly`: the most streamed movies (default `weekly`).
- `movies -a report -i <id> -r <reason> [-message "..."]`: report a movie for moderator review.
- `movies -a refresh_justwatch -i <id>`: queue a refresh of the movie's JustWatch links (VIP only).
- `movies -a sentiments -i <id>`: good and bad sentiments from a movie's comments and reactions.

### Changed

- `users -a add_hidden_items|remove_hidden_items`: `-section` is now checked before the request, like
  `hidden_items` (`calendar` is the default).
- `users -a hidden_items`: `-section progress_watched_reset` is no longer accepted; it is not a section in
  the Trakt API contract.
- `sync -a get_collection`: `-t` is now checked before the request (`movies`, `shows`, `episodes`, `media`,
  `seasons`); an unknown type fails with an error instead of calling the API. `-t media` returns movies,
  shows and episodes together.
- All modules: a failed action now prints its error as `<module>/<action>: <cause>`, with a space after the colon.

### Fixed

- VIP account limits (`420`) for `users -a add_list|add_list_items`, `sync -a add_to_watchlist` and
  `lists -a items`: a non-VIP user now gets the Trakt VIP page opened (from `X-Upgrade-URL`), a VIP user gets
  `account limit exceeded (limit: N)`. Before, the command showed a generic error.
- VIP-only actions (`426`): when the response has no `X-Upgrade-URL` header, `https://trakt.tv/vip` is opened
  instead of an empty URL.
- `movies -a favorited|played|watched|collected -period <p>`: `-period` was ignored and `weekly` was always
  used (the `shows` default overwrote it).
- `movies|shows|people -a updates|updated_ids -start_date <date>`: `-start_date` was ignored and the last
  60 days were always used. Without `-start_date` the last 60 days are still used.
- `calendars -start_date <date>`: `-start_date` was ignored; the request used a date 60 days in the past in
  the wrong format. The default is today again, for `-days` days (7).
- `type` from the config file was ignored: commands that use the global `-t` flag (for example `sync`,
  `watchlist`, `collection`, `history`, `certifications`, `comments`) always used `-t`'s default `movies`
  when `-t` was not given. The config file value is now used, and `-t` still overrides it.
- A global flag given before the module name could hide a config file value for an unrelated one-letter
  flag with the same first letter (for example `-translations` made `type` fall back to the `-t` default).
- `users -a history -item_id <id>` without `-t`: the request URL had no type (`history//<id>`), which is
  not an API route; the command now asks for `-t` (for example `-t movies`).
- `lists -a trending`: the output file name was empty, so the export was not written; it is now
  `export_lists_trending.json`.
- `lists -a items` without `-t`: it used the global default type `movies`, which is not a list item type;
  it now requests all item types (`movie,show,episode,season`).
- Nil pointer panics when a request fails before Trakt answers (no network, DNS, timeout) are fixed in
  the shared client, the API services and the command handlers; the command now prints the error.
- Posting a comment or a note: an error other than a validation error (for example 401) now shows the
  HTTP error; before, the command crashed with a nil pointer panic.
- `lists -a list` and `users -a list`: a 500 from Trakt now gives an error instead of a panic.
- `users -a report|list_report`: a 400 or 409 response without a `message` no longer panics.
- `checkin -a movie|episode`: if you are already checked in, the command now fails with the 409 error;
  before, it exited without printing anything.
- `users -a follow|block`: an unknown user now gives "user not found"; before, the command crashed
  with a nil pointer panic.
- `certifications -t movies|shows`: the "write data to:" message is printed to stdout on its own line;
  before, it went to stderr with no newline and ran into the next line of output.

## [1.15.3] - 2026-09-23

### Added

- `API_COVERAGE.md`: every Trakt API route from the official contract, whether trakt-sync
  implements it, and the Go method that calls it. Linked from the README.
- Contributor and AI agent guidelines: `AGENTS.md`, `CLAUDE.md` and `.agents/rules/`.
- `CHANGELOG.md` (this file). Every PR adds its entry under `[Unreleased]`.

### Changed

- README: API credentials are now created in the
  [Trakt developer portal](https://developer.trakt.tv/apps/new) (a verified GitHub account is
  required); new "API documentation" section; `episodes` added to the command list.
- Releases are now built and published by GitHub Actions (GoReleaser) when a `v*` tag is
  pushed.

### Fixed

- `users -a list_like` now likes the list. It sent `DELETE`, which removed an existing like
  instead.
- `users -a follower_requests` now returns the follow requests waiting for your approval. It
  returned your own pending requests to other users instead.
- `users -a follower_requests` and `users -a following_requests` now send `-ex` (extended info)
  to the API; it was silently dropped.
- `-version` now shows the real version and commit for release binaries and `make build`
  builds, instead of `dev` / `none` or a Go pseudo-version.

## [1.15.2] - 2026-06-16

### Fixed

- `sync -a get_watched` and `users -a watched` returned only the first page of results; they now fetch every page
  (up to `pages_limit`).

## [1.15.1] - 2026-06-16

### Changed

- Module docs moved from `README.md` into `docs/<module>.md`; the README links to them.

### Fixed

- `collection` returned only the first page of results; it now fetches every page (up to `pages_limit`).

## [1.15.0] - 2026-05-26

### Added

- `users`: many new actions for the user endpoints: `profile`, `follower_requests`,
  `following_requests`, `follow_request`, `hidden_items`, `add_hidden_items`, `remove_hidden_items`, `likes`,
  `collection`, `comments`, `notes`, `history`, `ratings`, `watchlist`, `watchlist_comments`, `favorites`,
  `favorites_comments`, `watching`, `report`.
- `users`: manage personal lists with `list`, `add_list`, `update_list`, `delete_list`, `reorder_lists`,
  `collaborations`, `list_likes`, `list_like`, `list_items`, `add_list_items`, `remove_list_items`,
  `reorder_list_items`, `update_list_item`, `list_comments` and `list_report`.
- `users`: social actions `follow`, `unfollow`, `followers`, `following`, `friends`, `blocked_users`, `block` and
  `unblock`.

### Changed

- `README.md` and `sync` action descriptions updated.

## [1.14.0] - 2026-05-11

### Added

- New `sync` module: `last_activities`, `playback`, `remove_playback`.
- `sync -a get_collection|add_to_collection|remove_from_collection`.
- `sync -a get_watched|get_history|add_to_history|remove_from_history`.
- `sync -a get_ratings|add_to_ratings|remove_from_ratings`.
- `sync -a get_watchlist` (with `-sort_by` / `-sort_how`), `update_watchlist`, `add_to_watchlist`,
  `remove_from_watchlist`, `reorder_watchlist`, `update_watchlist_item`.
- `sync -a get_favorites`, `update_favorites`, `add_to_favorites`, `remove_from_favorites`, `reorder_favorites`,
  `update_favorite_item`.

## [1.13.0] - 2025-06-01

### Added

- New `episodes` module: `summary`, `translations`, `comments`, `lists`, `people`, `ratings`, `stats`, `watching`,
  `videos`.

## [1.12.0] - 2025-05-30

### Added

- New `seasons` module: `summary`, `season`, `episodes`, `translations`, `comments`, `lists`, `people`, `ratings`,
  `stats`, `watching`, `videos`.

## [1.11.0] - 2025-05-27

### Changed

- Dates are handled in the timezone from your Trakt account settings. The settings are downloaded after login and
  stored in the JSON file set by `settings_path` in the config file.
- `-start_date` is rounded down to the full hour, and dates default to the RFC 3339 format.

### Fixed

- An invalid or missing settings or token file no longer stops the program; the file is created or refreshed.

## [1.10.0] - 2025-05-21

### Added

- New `shows` module: `trending`, `popular`, `favorited`, `played`, `watched`, `collected`, `anticipated`,
  `updates`, `updated_ids`, `summary`, `aliases`, `certifications`, `translations`, `comments`, `lists`,
  `collection_progress`, `watched_progress`, `reset_show_progress`, `people`, `ratings`, `related`, `stats`,
  `studios`, `watching`, `next_episode`, `last_episode`, `videos`, `refresh`.

## [1.9.1] - 2025-05-04

### Changed

- Internal: linter and `gofmt` cleanup, and a lint check in GitHub Actions.

## [1.9.0] - 2025-05-04

### Added

- New `networks` module: `networks -a list`.
- New `notes` module: add notes to movies, shows, seasons, episodes, people, history, collection and rating items
  (`notes -a notes -t <type> -i <id> -notes "..."`), and get, update or delete one (`notes -a note -i <id>`), or
  get its item (`notes -a item -i <id>`).
- New `recommendations` module: `recommendations -a movies|shows`, with `-ignore_collected`,
  `-ignore_watchlisted`, or `-i <id> -hide` to hide a recommendation.
- New `scrobble` module: `scrobble -a start|pause|stop -t movie|episode|show_episode -i <id> -progress <n>`; for
  `show_episode` pick the episode with `-episode_code 1x5` or `-episode_abs 164`.

### Changed

- `README.md` reorganized with a section per module.
- Tests run in GitHub Actions.

## [1.8.0] - 2025-03-26

### Added

- New `movies` module: `trending`, `popular`, `favorited`, `played`, `watched`, `collected` (with `-period`),
  `anticipated`, `boxoffice`, `updates`, `updated_ids`, `summary`, `aliases`, `releases`, `translations`,
  `comments`, `lists`, `people`, `ratings`, `related`, `stats`, `studios`, `watching`, `videos`, `refresh`.

### Changed

- `lists` and `comments` accept a Trakt slug as well as a numeric Trakt ID.

## [1.7.0] - 2025-03-16

### Added

- New `countries`, `genres` and `languages` modules: `-t movies|shows` exports the list for movies or shows.

## [1.6.1] - 2025-03-15

### Added

- `pages_limit` in the config file limits how many pages paginated exports fetch.

## [1.6.0] - 2025-03-11

### Added

- New `comments` module: `comment` (get, update or `-delete`), `comments` (post on a movie, show, season, episode
  or list), `replies`, `item`, `likes`, `like` (`-remove` to unlike), `trending`, `recent` and `updates`.
- New `certifications` module: `-t movies|shows` exports the certifications.

## [1.5.0] - 2025-02-26

### Added

- New `checkin` module: `checkin -a movie|episode -trakt_id <id> -msg "..."`,
  `checkin -a show_episode -trakt_id <id> -episode_code 1x5` (or `-episode_abs 6`), and `checkin -a delete`.
- `users -a settings`: settings of the current user.
- `people -a refresh`: queue a refresh of a person's data.

## [1.4.2] - 2025-02-21

### Changed

- `help` lists the commands in alphabetical order.
- `lists` docs updated.

## [1.4.1] - 2025-02-20

### Changed

- Same code as 1.4.0; this is the GitHub release for it.

## [1.4.0] - 2025-02-20

### Added

- New `lists` module: `trending`, `popular`, `list`, `likes`, `like` (`-remove` to unlike), `items` (filter with
  `-t movie,show`) and `comments`, selected with `-trakt_id <id>`.

## [1.3.1] - 2024-12-07

### Changed

- Internal: error handling cleanup.

## [1.3.0] - 2024-11-12

### Added

- `users -a watched -t movies|shows -u <user>`: watched movies or shows of a user; `-ex noseasons` leaves out the
  seasons.

## [1.2.0] - 2024-11-07

### Added

- `users -a stats -u <user>`: stats of a user.

## [1.1.0] - 2024-10-29

### Added

- `users -a saved_filters -u <user>`: saved filters of a user (VIP); without VIP it opens the upgrade page in the
  browser.

### Changed

- **Breaking:** the `lists` command moved to `users -a lists` (`users -a lists -u <user> [-i <list id> -t <type>]`).

## [1.0.10] - 2024-10-14

### Fixed

- `people -a updates` and `people -a updated_ids` ignored `-start_date` and always used the current date.
- `lists -u <user>` ignored the `-u` flag.

## [1.0.9] - 2024-10-14

### Changed

- Internal: linter cleanup and refactoring.

## [1.0.8] - 2024-10-13

### Changed

- Internal: linter cleanup and refactoring.

## [1.0.7] - 2024-10-06

### Changed

- Internal: linter cleanup and refactoring.

## [1.0.6] - 2024-10-04

### Changed

- Internal: linter cleanup and refactoring of the `people` command.

## [1.0.5] - 2024-10-03

### Changed

- Internal: linter cleanup, console output moved to one printer package, unhandled errors checked.

## [1.0.4] - 2024-09-30

### Fixed

- `-version` printed `dev` for binaries installed with `go install` instead of the module version.

## [1.0.3] - 2024-09-30

### Changed

- Internal: linter cleanup.

## [1.0.2] - 2024-09-29

### Changed

- Internal: linter cleanup.

## [1.0.1] - 2024-09-12

### Added

- Release binaries built and published with GoReleaser.

## [1.0.0] - 2024-09-12

### Added

- First release, with the `calendars`, `collection`, `help`, `history`, `lists`, `people`, `search` and `watchlist`
  commands exporting Trakt data to JSON.

[Unreleased]: https://github.com/mfederowicz/trakt-sync/compare/v1.17.0...HEAD
[1.17.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.16.0...v1.17.0
[1.16.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.3...v1.16.0
[1.15.3]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.2...v1.15.3
[1.15.2]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.1...v1.15.2
[1.15.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.0...v1.15.1
[1.15.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.10.0...v1.11.0
[1.10.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.9.1...v1.10.0
[1.9.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.9.0...v1.9.1
[1.9.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.6.1...v1.7.0
[1.6.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.6.0...v1.6.1
[1.6.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.4.2...v1.5.0
[1.4.2]: https://github.com/mfederowicz/trakt-sync/compare/v1.4.1...v1.4.2
[1.4.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.3.1...v1.4.0
[1.3.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.10...v1.1.0
[1.0.10]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.9...v1.0.10
[1.0.9]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.8...v1.0.9
[1.0.8]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.7...v1.0.8
[1.0.7]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.6...v1.0.7
[1.0.6]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.5...v1.0.6
[1.0.5]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/mfederowicz/trakt-sync/releases/tag/v1.0.0
