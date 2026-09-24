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
schedule. Releases up to v1.15.2 are listed on
[GitHub Releases](https://github.com/mfederowicz/trakt-sync/releases).

## [Unreleased]

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

### Changed

- `search -a text_query` and `search -a id_lookup` replace `text-query` and `id-lookup`, to match the underscore action
  names of the other modules. The old names still work and print a deprecation note.
- `search` with an unknown or missing `-a` prints the same "Available actions" list as the other modules.

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

[Unreleased]: https://github.com/mfederowicz/trakt-sync/compare/v1.16.0...HEAD
[1.16.0]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.3...v1.16.0
[1.15.3]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.2...v1.15.3
