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

[Unreleased]: https://github.com/mfederowicz/trakt-sync/compare/v1.15.2...HEAD
