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

The Go library (`trakt`, `str`, `uri`) follows the same version numbers:

- **While it is experimental** (now) - breaking changes to its Go API (a removed or renamed
  exported identifier, a changed signature, or a behavior callers rely on) may land in a minor
  release. They are listed under `### Changed` or `### Removed` with a **Library:** prefix.
  New exported identifiers are minor, fixes are patch, as for the CLI.
- **Once it is declared stable** (the maintainer's call; the README and the `trakt` package docs
  will say so) - a breaking change to its Go API needs a major version.
- Any major version, CLI or library, also changes the module path to
  `github.com/mfederowicz/trakt-sync/v2` (Go modules ignore `v2+` tags without it).

A version is tagged once a module (or a meaningful fix) is done - there's no fixed release
schedule.

## [Unreleased]

### Added

- `recommendations -a movies|shows` and `social_recommendations -a movies|shows` read the `favorited_by` and `recommended_by` entries in both shapes: the nested one (`user` object plus `notes`) and the flat one of the Trakt API contract (profile fields plus `notes`). A flat entry kept only its `notes` before. The export writes `user` plus `notes` in both cases. Both shapes were tested against a mock server only: on 2026-10-04 the live API returned empty lists.
- Library: `str.UserNotes` decodes a flat entry (profile fields plus `notes`) into `User` and `Notes`, next to the nested one.
- Exports keep three newer groups of fields the API returns: `email` in `users -a settings`, `height` of a person with `-ex full` (for example `people -a summary`), and the VIP veteran fields of a user profile with `-ex vip` (`vip_veteran_since`, `vip_veteran_years`, `vip_veteran_tier`, `vip_veteran_title`, `vip_grace_ends_at`). These fields were dropped before. On 2026-10-04 the live API sent `height` (in centimetres) and the VIP veteran fields; `email` is in the Trakt API contract, but the live API did not send it yet.
- Library: `str.UserProfile` has the new fields `Email`, `VipVeteranSince`, `VipVeteranYears`, `VipVeteranTier`, `VipVeteranTitle` and `VipGraceEndsAt`; `str.Person` has `Height`.

### Changed

### Fixed

- The config file keys `ignore_collected`, `ignore_watched` and `ignore_watchlisted` can be written without quotes (`ignore_collected = true`). Before, an unquoted value stopped every run with `cannot parse the config file: toml: ... incompatible types: TOML value has type bool; destination has type string`. The quoted form (`"true"`) works as before.
- `users -a collection`: the docs and the 1.24.0 notes said it takes the media filters (`-genres`, `-years`, `-watchnow` and so on). The flags are sent, as the Trakt API contract lists them for this route, but on 2026-10-04 the live API ignored `-genres`, `-years` and `-languages` and returned the whole collection. `docs/users.md` now says so.

## [1.24.0] - 2026-10-04

### Added

- `movies -a trending|popular|anticipated|watched|played|collected|favorited|hot|streaming` take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`, `-start_date`, `-end_date` and `-watchnow`. The flags were accepted or missing before, but no filter was sent, so the lists were never filtered.
- Library: new `uri.MediaFilters` and the `uri.ListOptions.Filters` field carry the media filters (`watchnow`, `genres`, `subgenres`, `years`, `ratings`, `start_date`, `end_date`, `runtimes`, `countries`, `certifications`) of the list routes.
- `shows -a trending|popular|anticipated|watched|played|collected|favorited` take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`, `-start_date`, `-end_date` and `-watchnow`, plus `-status` (for example `-status ended`). No filter was sent before, so the lists were never filtered. The live API does not always apply `status`: on 2026-10-03 it returned the unfiltered list.
- `media -a trending|popular|anticipated` and `recommendations -a movies|shows` take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`, `-start_date`, `-end_date` and `-watchnow`.
- Every `calendars` action takes the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications` and `-watchnow`. The dates of a calendar stay `-start_date` and `-days`.
- `lists -a trending|popular|items` take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`, `-start_date`, `-end_date` and `-watchnow`.
- `users -a watchlist|history|collection|list_items` take the Trakt media filters: `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`, `-start_date`, `-end_date` and `-watchnow`. In `history`, `-start_at` / `-end_at` still set the watched-at window; `-start_date` / `-end_date` are the separate media date filters.
- `users -a activities` sends the rest of the Trakt media filters: `-subgenres`, `-ratings`, `-certifications`, `-start_date`, `-end_date` and `-watchnow`, next to `-genres`, `-years`, `-runtimes` and `-countries`. The six flags were accepted before, but not sent. An unknown `-watchnow` value stops the run before any request.
- Library: `uri.SocialActivityOptions` has the new fields `Certifications`, `EndDate`, `Ratings`, `StartDate`, `Subgenres` and `WatchNow`.
- The config file keys `ignore_collected`, `ignore_watched`, `ignore_watchlisted` (quoted `"true"` or `"false"`) and `watch_window` now set the defaults of the flags of the same name in `recommendations`, `social_recommendations`, `smart_lists -a items` and `users -a activities`. The keys were read but never applied; only the flags worked. A flag still wins over the config file.
- `users -a activities` takes `-ignore_watched`, `-ignore_collected` and `-ignore_watchlisted` (`true` or `false`), like `recommendations`. As of 2026-09-25 Trakt answers this route with 401 for API apps, so the flags were tested against a mock server only.
- Library: `uri.SocialActivityOptions` has the new fields `IgnoreCollected`, `IgnoreWatched` and `IgnoreWatchlisted`.
- `movies -a trending|popular|anticipated|watched|played|collected|favorited|hot|streaming` take three more filters: `-imdb_ratings` (for example `8.0-10.0`), `-rt_meters` and `-rt_user_meters` (for example `90-100`).
- `shows -a trending|popular|anticipated|watched|played|collected|favorited` take the same three filters: `-imdb_ratings`, `-rt_meters` and `-rt_user_meters`.
- `-languages` (for example `-languages en,pl`) filters the `movies` and `shows` list actions. Since 1.22.0 the flag only printed a note that the Trakt API no longer supports it; the live API does apply the filter, so the note is gone and the value is sent.
- Exports of movies and shows keep more of what the API returns: `images` (with `-ex images`), `colors` (with `-ex colors`), and with `-ex full` also `subgenres`, `original_title`, `social_ids`, `after_credits` / `during_credits` (movies) and `last_aired` / `total_runtime` (shows). These fields were dropped before. This covers every export that holds a movie or a show, for example `movies`, `shows`, `media`, `recommendations`, `calendars`, `users -a watchlist|history`.
- Library: `str.Movie`, `str.Show`, `str.Recommendation` and `str.Media` have the new fields `Images`, `Colors` and `SocialIDs` (new types `str.MediaImages`, `str.MediaColors`); `str.Movie`, `str.Show` and `str.Recommendation` also `Subgenres` and `OriginalTitle`, plus `AfterCredits` / `DuringCredits` (movie) and `LastAired` / `TotalRuntime` (show).
- Library: `str.Recommendation` has the extended info fields of a movie or a show (`Language`, `Genres`, `Overview`, `Runtime`, `Rating`, `Released`, `FirstAired`, `Airs`, `Network` and so on).
- Library: `uri.MediaFilters` has the new fields `ImdbRatings`, `Languages`, `RtMeters` and `RtUserMeters`. `uri.ListOptions.Languages` is no longer marked deprecated.
- `-languages` also filters `media -a trending|popular|anticipated`, `recommendations -a movies|shows`, every `calendars` action, `lists -a items` and `users -a watchlist|history|list_items`. It is sent by `lists -a trending|popular` and `users -a collection` too, but on 2026-10-03 the live API ignored it there.

### Changed

- A word that is not a flag or a flag's value now stops the run with a message that names it, for example `movies: unexpected argument "foo"` (it was `invalid flags`). A stray word that happened to be a module or flag name, as in `movies shows -a trending`, was accepted before and the flags after it were silently ignored; it is rejected now. `help <module>` works as before.
- Every `calendars` action and `recommendations -a movies|shows` now end with `empty result` and exit status 1 when the API returns no items, like the other list actions. Before, they wrote a file holding an empty list `[]` and ended with exit status 0. The output file is not written (an existing one is left as it is).
- `users -a watchlist|history|favorites|ratings` do the same: an empty list, or a filter that matches nothing, ends with `empty result` and exit status 1 instead of a file holding `[]`.

### Fixed

- A flag value that starts with a dash (for example a negative number) was read as a flag name and could stop the run with `invalid flags`.
- Flags written as `-flag=value` after the module name (for example `movies -a=trending` or `shows -status=ended`) stopped the run with `invalid flags`; only `-flag value` worked. Both forms are accepted now. An empty value (`-i ""`) no longer crashes the flag check.
- `users -allow_comments` and `users -display_numbers` (used by `users -a update_list`) stopped the run with `invalid flags`. They are accepted now. The deprecated `-query`, `-studio_ids` and `notes -notes_id` also ended with `invalid flags` instead of their deprecation note.
- Library: `uri.StatusOptions` listed the show status `running series`; the API value is `returning series`.
- `recommendations -a movies|shows -ex full` and `social_recommendations -a movies|shows -ex full` exported only `title`, `year`, `ids` and the `favorited_by` / `recommended_by` lists; the extended fields the API returned were dropped. The export now keeps `language`, `languages`, `genres`, `overview`, `runtime`, `rating`, `certification`, `country`, `status` and the other `full` fields the `movies` and `shows` exports have.
- `users -a watchlist` with a filter that matches less than the whole watchlist (for example `-languages pl` or `-genres drama`) kept asking for empty pages, up to `pages_limit` or the page count of the unfiltered watchlist. It now stops at the first empty page.
- `lists -a items` and `users -a list_items` with a filter that matches nothing (for example `-languages pl` on a list without Polish titles) kept asking for the next page, up to `pages_limit` or the page count of the unfiltered list, because the page count the API sends ignores the filters. They now stop at the first empty page.

## [1.23.0] - 2026-10-02

### Changed

- `sync -a add_to_history` with `-t movies`, `-t seasons`, `-t episodes` or `-t all` now adds every play of the input back. An item watched several times came back with one play before (the first one in the input file), so its other plays were lost from the Trakt history. `-t shows` already kept every episode play. With `-t all`, an episode entry that also names its show was sent twice, with the show and as an episode; it is now sent once, with the show. To keep one play per item, leave only that entry in the input file.

### Fixed

- `sync -a add_to_history` with `-t shows` or `-t all` sent a show that has no episode in the input file (a whole show) with an empty `seasons` list and without its `watched_at`. Such a show is now sent without `seasons`, the form Trakt documents for adding all episodes of a show, and with its `watched_at` when the entry has one.
- Library: `str.Show` has a new `WatchedAt` field, used when a whole show is added to the history.
- Library: `str.HiddenItem` has a new `User` field for hidden users. The live API does not return them yet: `GET /users/hidden/comments` answers with an empty list even when a user is hidden (see the Findings in `API_COVERAGE.md`), so `users -a hidden_items -section comments` still exports `[]`.
- `users -a add_hidden_items|remove_hidden_items -section comments` wrote a result file without the users: the `users` counter under `added` / `deleted` and the `users` and `people` lists under `not_found` were dropped, so the file did not show whether a user was hidden. They are now written.
- Library: `str.ResultCounters` has a new `Users` field and `str.ResultNotFound` new `People` and `Users` fields; they were dropped from add and remove results before.
- A device login that failed (no code, code denied, expired, already used or unknown, or not approved in time) went on to run the module without a token. The run now stops with `Error: device login failed: ...` and exit status 1. A code that is not approved in time is reported as `Error: device login failed: time out, the device code was not approved` instead of `Time out!`.
- The device login kept polling until the code's lifetime ran out after Trakt had already given a final answer: code denied, expired, already used or unknown. It now stops at once with `Error: device login failed: device code denied, your device is not connected` (or `device code expired`, `device code already used`, `invalid device code`). A device login that could not get a code at all is now reported as `Error: ...`; it ended without any message before.
- The device login (first run, or an expired token that cannot be refreshed) crashed with `panic: runtime error: invalid memory address or nil pointer dereference` when a polling request got no response, for example on a network error. The failed attempt is now reported as `Error: ...` and polling goes on. Polling also ends with a time out when the code's lifetime is not a multiple of the polling interval; it never ended before.
- `sync -a get_collection -t seasons` stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` on every run. It now exports the collected seasons, each with its own IDs, and follows the pages of the collection.
- Library: `SyncService.GetCollectedSeasons` returned a nil `*str.Response` on success and a list in which every item pointed at the last collected season. It now returns the response of the collection request and one item per season.
- `users -a add_hidden_items|remove_hidden_items -section comments` sent an empty `users` list, so no user was hidden or unhidden although the command ended without an error. The `user` items of the input file are now sent (without `-t`, or with `-t user`).
- `history`, `watchlist` and `collection` with `-t episodes` stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` when an exported episode had no title, which Trakt allows. Such an episode is now exported with the title `no episode title`, as an episode with an empty title already was.
- `shows -a last_episode` and `shows -a next_episode` stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` when the episode had no title, which Trakt allows (typical for an upcoming episode). The episode is now exported and the message shows `no episode title`.
- `comments -a comments` without `-t`, or with a `-t` other than `movie`, `show`, `season`, `episode` or `list` (for example the plural `movies`), stopped with `panic error:runtime error: invalid memory address or nil pointer dereference`. It now prints the possible types and stops with `unknown type "movies"`.
- `checkin -a movie|episode|show_episode` and `comments -a comments` stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` when the user settings returned by Trakt had no `connections` object (the API allows it to be missing). The checkin or comment is now sent without sharing overrides, so Trakt applies the account's own sharing settings.
- `scrobble -a start|pause|stop -t movie`, `checkin -a movie`, `notes -a notes -t movie|show|person|rating|collection` and `comments -a comments -t movie` ignored a failed lookup of the movie, show or person (for example an unknown ID) and sent the request to Trakt without the item. They now stop with the lookup error, such as `fetch movie error:...`, and send nothing. `comments -a comments` also stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` when the user settings could not be read; it now reports `user connections error:...`.
- `checkin -a show_episode` and `scrobble -a start|pause|stop -t show_episode` without `-episode_code` and `-episode_abs` ended with exit status 0 although nothing was sent to Trakt. They now stop with `set episode ie: -episode_code 1x5 or -episode_abs 6` and exit status 1.
- `sync -a reorder_watchlist|reorder_favorites` and `users -a reorder_list_items|reorder_lists` stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` when an input item had no `id` (for `reorder_lists`: a list without a Trakt ID). They now stop before any request with an error that names the item, such as `item at index 1: has no id`. A `null` item in the input of the collection, watchlist, favorites, list items and hidden items actions is reported the same way.
- `sync -a add_to_history|remove_from_history|add_to_ratings|remove_from_ratings` silently dropped every input item that had neither `watched_at` nor `rated_at`, so the request went out with empty lists and nothing changed on Trakt, although the command ended without an error. Such items are now sent: the dates are optional in the Trakt API (a history item without `watched_at` is marked as watched now, and removing needs only the IDs).
- Library: `str.ItemsList.GetUniqueOldest` dropped items without `WatchedAt` and `RatedAt`. It now keeps them, unless another item with the same Trakt ID has a date.
- `sync -a add_to_history|remove_from_history|add_to_ratings|remove_from_ratings` stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` when an item in the input had no Trakt ID (for example a movie with only an `imdb` ID), when a show's episode had no `season` number, or when the list held a `null` item. It now stops before any request with an error that names the item, such as `item at index 1: movie has no trakt id`. A show item without an `episode` that follows one with an episode no longer panics either.
- Library: `str.ItemsList.GetUniqueOldest` kept the entry with the latest `watched_at` / `rated_at` per Trakt ID, although its name and docs say the oldest. It now keeps the oldest. It also no longer panics on a nil list, on an item without a Trakt ID (such items are skipped), or when a watched and a rated item share an ID, and `Uniq` / `GetUniqIDs` leave nil lists nil instead of panicking. Nothing changes for CLI users: `sync` builds these lists with one entry per ID already.

## [1.22.0] - 2026-10-01

### Added

- Library: examples on pkg.go.dev for the `trakt` package (client setup, OAuth token, pagination, typed errors, timezones) and for `uri.AddQuery`. `go test` compiles them, so they stay in step with the API.

### Changed

- The `write data to:` / `write result to:` line went to stderr without a line break, so it ran into the next output. It now goes to stdout on its own line, like the newer modules, and a JSON encoding error stops the command instead of writing an empty file. This covers every module.
- `calendars` actions use underscores like the other modules: `my_shows`, `all_new_shows`, `hot_releases` and so on. The old hyphenated names (`my-shows`, `all-new-shows`, `hot-releases`, ...) still work, but print a deprecation note. Export file names are unchanged.
- Library: a `trakt` client without `WithUserAgent` now sends `User-Agent: trakt-sync-go` (`trakt.DefaultUserAgent`), so apps built on the library no longer show up as the trakt-sync CLI. The CLI still sends `trakt-sync/<version>`. The package docs now tell apps to set their own User-Agent with `WithUserAgent`.
- CI also builds and tests the `trakt`, `str` and `uri` packages and `example/` with the Go version in `go.mod` (1.21), so the library keeps working for that minimum.
- `people -a lists` has new `-t` (`all`, `personal`, `official`, `watchlist`, `favorites`) and `-s` (`popular`, `likes`, `comments`, `items`, `added`, `updated`) flags and asks for `personal` lists sorted by `popular` by default. Before, it sent the action name as the list type and `rank` as the sort, which are not valid values. The Trakt API currently returns no lists for any person on this route, so the command still ends with `empty lists`.
- CI fails when the test coverage of a package drops below its floor in `.github/coverage-floors.txt`; `make cover-check` runs the same check locally.

### Deprecated

- Library: the `uri.ListOptions` fields `EpisodeTypes`, `Languages`, `Metascores`, `NetworkIDs`, `StudioIDs` and `TmdbRatings` are marked deprecated, because their query parameters are no longer part of the Trakt API. They still exist and will be removed in the next major version.
- Flags that no longer do anything now print a note when used, and will be removed in the next major version: `-studio_ids` and `-languages` (no longer supported by the Trakt API), `-query` (use `-q` with `search`), `notes -notes_id` (use `-i`) and `scrobble -delete` (never had an effect). Commands still run as before.

### Fixed

- `sync -a add_to_history|remove_from_history|add_to_ratings|remove_from_ratings` with a `-t` other than `all`, `movies`, `shows`, `seasons` or `episodes` (for example `-t movie`) stopped with `panic error:runtime error: invalid memory address or nil pointer dereference` while reading the items. It now stops with an error that lists the valid types.
- An unknown `module` in the config file made every command stop with `type 'movies' is not valid for module 'history'`, right after the note `Forcing module to history`. The fallback to `history` now works, so commands run.
- `verbose = true` in the config file had no effect: verbose output only appeared with `-v`. The config file value is now used, and `-v=false` before the module name turns it off for one run.
- `scrobble -progress` was ignored, so scrobbles were sent without a progress value; a `scrobble -a stop` could therefore not mark an item as watched. The given progress is now sent.
- `comments -spoiler` was ignored: comments and replies were always posted without the spoiler flag. `-spoiler` now marks them as spoilers.
- If the new token (after login or a token refresh) or the refreshed user settings could not be encoded, `token.json` / `user_settings.json` was overwritten with an empty file. The encoding error is now reported and the stored file is kept.
- An unknown or missing `-a` (and an unknown `-t` for `certifications`, `countries`, `genres`, `languages`, and for `notes`/`scrobble` actions, or `-item` for `notes`) printed the list of valid values but ended with exit status 0, so scripts could not notice a typo. It now prints the list followed by an error such as `calendars: unknown action "typo"` and exits with status 1. `notes` with an unknown action used to report a misleading privacy error; it now shows the actions list like the other modules.
- README: the config file's `type` was described as applying to every module. It is ignored by `episodes`, `lists`, `movies`, `people`, `scrobble`, `seasons`, `shows` and `users`, which take `-t` only from the command line; the README now says so.

## [1.21.0] - 2026-10-01

### Added

- Library docs for the `trakt` package: a package overview for pkg.go.dev, a "Library usage" section in the README, and runnable programs in `example/` (device login, trending movies, paginated history, error handling). Method docs link to the new API reference at docs.trakt.tv, and the versioning rules at the top of this file now cover the library.

### Changed

- The Trakt API client is now an importable Go package, `github.com/mfederowicz/trakt-sync/trakt` (it was `internal/`, which other modules cannot import). It is experimental: its API may still change in minor releases, as the library versioning rules at the top of this file describe. Nothing changes for CLI users.
- Library: `trakt` API cleanup before its first release. Methods get consistent names: `Oauth.PollForAccessToken`, `Checkin.CheckIn`, `Comments.AddComment`, `Users.GetSettings`, `Users.LikeList`, `Users.ReportList`, `Users.GetWatching` and `Users.GetListItemsByType`. Every service method now returns the `*str.Response` as well (the `Sync` add/remove/update/reorder methods, `Shows.GetShowCollectionProgress`/`GetShowWatchedProgress` and `Users.AddHiddenItems`/`RemoveHiddenItems` did not), so callers can read rate limit and pagination headers. Path parameters are plain values instead of pointers (`id string`, `days int`), so a nil can no longer panic; an optional segment is left out when it is empty (`""` or `0`). In `Users`, an empty user means the authenticated user (`me`) in every method except `Follow`, `Unfollow`, `Block`, `Unblock` and `Report`, and type or section segments are always sent. `str` and `uri` get package docs for pkg.go.dev, and `str.Format`/`str.Formatc` are removed along with the `github.com/wissance/stringFormatter` dependency. Nothing changes for CLI users.
- Debug lines such as `fetch ... url:` and `create new checkin`, and the bare request URLs that `users` actions such as `likes`, `collection`, `followers` and `history` printed, now print only with `-v`, which also shows each API request as `METHOD <url>` (with `client_secret` redacted). This covers `calendars`, `checkin`, `comments`, `lists`, `movies`, `notes`, `people`, `recommendations`, `scrobble`, `search`, `shows`, `sync`, `users`, `countries`, `certifications`, `genres`, `languages` and `networks`. Errors are still printed as before.

### Fixed

- Library: `uri.AddQuery` panicked when given nil options (for example a nil `*uri.ListOptions` passed to a `trakt` service method). Nil options now add no query.
- After an expired access token was refreshed, user settings were fetched with the old token, so that settings refresh failed. They are now fetched with the new token.
- Some API errors were only printed as `General error occurred:` or dropped: 412, 502, 503, 504 and other statuses without their own error type. The command then carried on as if the request had worked. These errors now stop the command with a readable message (`GET <url>: 503 <message>`) and exit status 1.
- A 429 (rate limit) response without a `Retry-After` header, or a 426 (VIP required) response without an `X-Upgrade-URL` header, ended the command with a `panic error` (nil pointer). It now reports the API error instead.
- `lists -a items` and `sync -a add_to_watchlist` now open the VIP upgrade page on a 426, as the other VIP actions already did.
- `users -a plex_*`: Plex errors on 502/503/504 now show Plex's `error_code`, message and guidance, not just the status.

## [1.20.0] - 2026-09-30

### Changed

- `token.json`, `user_settings.json` and JSON exports (`-o`) are now written readable only by the owner
  (0600); an existing 0644 token or settings file is tightened the next time it is written.
- Every request now sends `User-Agent: trakt-sync/<version>`, `trakt-api-version: 2` and
  `Content-Type: application/json`, and no longer sends an empty `Authorization` header. `-v` prints the
  User-Agent that is sent.
- README: new "Intended use" section - trakt-sync is for your own Trakt data only.

### Fixed

- `-v` printed the full access token and client ID; both are now masked to their last 4 characters.
- The list of available actions printed for an unknown `-a` now shows only real actions: `movies` lists
  `updates` and `related` instead of `updated` and `releated`; `shows` lists `related` and
  `reset_show_progress` and no longer offers `boxoffice`, `releases` or `releated`; `users` no longer
  offers `follow_request`.

## [1.19.1] - 2026-09-26

### Changed

- CI: workflows run on `ubuntu-26.04` and use Node 24-based `actions/checkout@v7` / `actions/setup-go@v7`,
  which clears the Node.js 20 deprecation and `ubuntu-latest` migration warnings.

## [1.19.0] - 2026-09-25

### Added

- New `younify` module for streaming service connections: `younify -a connections` (services and your connection
  status), `younify -a connect -service_id <id> [-return_url <url>]` (prints the web auth URL to link a service),
  `younify -a refresh -service_id <id> [-all_data]` (queue an incremental or full re-sync) and
  `younify -a disconnect -service_id <id>`. Trakt does not open younify to API apps yet: `connections` answers
  401 (the developer portal too), which the CLI reports as `younify is not open to API apps yet`.
- `users -a smart_lists|smart_list|add_smart_list|update_smart_list|delete_smart_list`: list, get, create, update and
  delete your smart lists. Create and update read the smart list JSON from `-items <file>` or stdin, like `add_list`.
- `users -a comment_reactions` (comments you reacted to), `users -a activities -t friends|followers|following` (what
  your social graph watched), `users -a month_in_review -year <y> -month <m>` and `users -a year_in_review -year <y>`.
  Trakt answers the last three with 401 for API apps for now; the CLI reports it as `this route is not open to API
  apps yet`.
- `users -a watchlist|favorites -sort <rank|added|title|released|runtime|popularity|percentage|votes>` uses the
  contract's `/{type}/{sort}` routes (with `-t all`, `movies` or `shows`); without `-sort` nothing changes.
- `users -a update_settings` (profile and browsing settings from `-items <file>` or stdin), `users -a add_saved_filters`
  and `users -a delete_saved_filter -i <id>` (both VIP only).
- `users -a data_syncs [-t younify|plex|import]`, `users -a data_sync|data_sync_paused|data_sync_skipped -i <id>` and
  `users -a undo_data_sync -i <id>`: the syncs your connected apps ran, their paused and skipped items (with the keys
  specific to the source kept), and undoing a sync. Trakt answers `data_syncs` with 401 for API apps for now; the CLI
  reports it as `this route is not open to API apps yet`.
- Plex settings in `users`: `plex_settings`, `update_plex_settings` (JSON from `-items <file>` or stdin),
  `plex_connect [-return_url <url>]`, `plex_disconnect`, `plex_servers`, `plex_server -i <server id>` and
  `plex_sync [-i <server id>] [-all_data]`. Plex's own errors (e.g. `bad_auth`) are shown with Trakt's guidance.
  Trakt answers the Plex routes with 401 for API apps for now; the CLI reports it as `this route is not open to API
  apps yet`.

### Changed

- `docs/users.md` documents the contract's singular list item types for `users -a lists -i <id> -t`: `movie`, `show`,
  `movie,show`, `movie,show,season,episode` (the action passes `-t` through unchanged).

### Fixed

- `users -a add_smart_list|update_smart_list`: an empty `source`, `media_type` or `privacy` in the JSON (e.g.
  `"source": ""`), or an empty `name` on update, passed the input check and was sent to Trakt; it is now rejected before the
  request.
- `trakt-sync help`: the command list is sized to the longest command name; before, the summary of
  `social_recommendations` was shifted out of its column.

## [1.18.0] - 2026-09-25

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
- New `team` module: `team -a members [-ex full|images]` exports the Trakt team members to
  `export_team_members.json`. User profiles in JSON exports now also include `deleted` and `director` when the API
  sends them.
- New `social_recommendations` module: `social_recommendations -a movies|shows`
  exports recommendations based on the people you follow, with `-ignore_watched`, `-ignore_collected`,
  `-ignore_watchlisted`, `-watch_window <days>` and `-ex`.
- New `watchnow` module: `watchnow -a sources [-country us]` exports the watch now sources
  (streaming providers) supported by Trakt, for all countries or one. Trakt marks it limited access; without access
  it fails with a "limited access" message.
- `recommendations -a movies|shows` takes `-ignore_watched true|false` and `-watch_window <days>`, like
  `social_recommendations`.
- New `smart_lists` module: `smart_lists -a summary -i <slug>` exports a smart list definition
  (name, media type, filters) and `smart_lists -a items -i <slug>` the movies or shows it resolves to, with media
  filters (`-watchnow`, `-genres`, `-subgenres`, `-years`, `-ratings`, `-runtimes`, `-countries`, `-certifications`),
  `-ignore_watched`, `-ignore_watchlisted` and `-ex`.

### Changed

- trakt-sync now exits with status 1 when a command fails (API error, unknown command or flag, invalid
  config); before it always exited with 0.
- Go API: the `str.UserProfile.Userame` field is renamed to `Username` (JSON output unchanged). Also removed an unused
  `sync:remove_watchlist` config entry.
- A flag the module does not know (e.g. a typo) now stops the command with `<module>: flag provided but not
  defined: -x` after the help. Before, the rest of the flags were skipped and the command ran anyway.
- `CHANGELOG.md` now covers every release back to 1.0.0.
- GitHub release notes now show the version's `CHANGELOG.md` section (Added / Changed / Fixed) instead of a list of
  commits.

### Fixed

- `users -a add_list` treated a created list (201) as an error and did not write its result file;
  `sync -a remove_playback` reported a successful removal (204) as an error.
- `-o` was ignored by `sync` and `users`, which always wrote to their generated file name (e.g.
  `export_sync_history_movies.json`), including the actions that write a result file (`add_*`, `remove_*`,
  `reorder_*`, `users -a update_list`, `users -a add_list`, `users -a watching`, ...); `-o` now sets the output file.
  Without `-o` the file names are unchanged.
- `users -a lists -i <id>` wrote the list items over the lists overview (`export_users_lists.json`); the items now go
  to `-o` (default `export_users_lists_<type>.json`) and the overview keeps its own file. Scripts that read the
  items from `export_users_lists.json` should read `export_users_lists_<type>.json` (or pass `-o`).
- `-h` / `-help` after a module name printed the help and then ran the command; it now only prints the help.
- `users -a update_list -description "..."`: `-description` was not a `users` flag, so the description was never
  sent; `users` now accepts it.
- A crash inside any command was caught but then reported nothing, as if the command had succeeded; it now
  prints `panic error:<reason>` (or `fatal error`).
- `checkin -a movie -trakt_id <id>` and `checkin -a show_episode -trakt_id <id>` looked the item up without its ID
  (`GET /movies/`, `GET /shows/`) since 1.8.0, so the checkin could not work; they now use `-trakt_id`.
- `checkin -a show_episode`: an active checkin (409) without an expiry time in the response now reports the
  existing checkin instead of crashing.
- `users -a history -start_at <date> -end_at <date>` failed with `flag provided but not defined: -start_at` and
  silently used the default window; `users` now accepts both flags.
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

[Unreleased]: https://github.com/mfederowicz/trakt-sync/compare/v1.24.0...HEAD
[1.24.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.23.0...v1.24.0
[1.23.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.20.0...v1.21.0
[1.20.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.19.1...v1.20.0
[1.19.1]: https://github.com/mfederowicz/trakt-sync/compare/v1.19.0...v1.19.1
[1.19.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.18.0...v1.19.0
[1.18.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.17.0...v1.18.0
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
