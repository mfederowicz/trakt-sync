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

- `calendars -a {my,all}-media`: movies and shows releasing in the selected period.
- `calendars -a {my,all}-streaming`: movies with a streaming release in the selected period.
- `calendars -a hot-releases|hot-premieres|hot-new-shows|hot-finales`: hot releases, premieres,
  new shows and finales in the selected period.
- New `media` module: `media -a trending|popular|anticipated` exports movies and shows together in one list.
- `comments -a reactions` and `comments -a reactions_summary`: export the reactions on a comment, or their totals by type.
- `comments -a reaction -reaction <type>`: add a reaction to a comment; add `-remove` to take it back.
- `comments -a report -r <reason> [-message "..."]`: report a comment for moderator review.

### Changed

- All modules: a failed action now prints its error as `<module>/<action>: <cause>`, with a space after the colon.

### Fixed

- Fewer nil pointer panics when a request fails before Trakt answers (no network, DNS, timeout): the
  command now prints the error. Fixed in the shared client and in `checkin -a delete`,
  `comments -a like`, `lists -a like|list`, `movies|people|shows -a refresh` and
  `users -a list|update_list|report|list_report`. Some commands can still panic in this case; they
  will be fixed separately.
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

[Unreleased]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.3...HEAD
[1.15.3]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.2...v1.15.3
