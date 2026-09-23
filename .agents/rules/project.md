---
description: 'Project overview, structure, tooling, restrictions, and commit standards for trakt-sync. Apply to all files.'
applyTo: '**'
---

# Project Guidelines

## Tech Stack

Go (`go 1.21` in `go.mod`; CI tests on `oldstable` and `stable`). Dependencies
are vendored (`vendor/`). Linting with [revive](https://github.com/mgechev/revive)
(`revive.toml`), formatting with `gofmt`, releases with GoReleaser.

## Project Structure

A request flows top to bottom through these packages:

```
main.go             # flags, config, token, then cmds.ModulesRuntime
cmds/               # one command_<module>.go per CLI module; registered in runtime.go
handlers/           # <module>_<action>_handler.go, implements Handler.Handle(*str.Options, *internal.Client) error
                    # commons.go / commons_helpers.go hold shared CommonLogic
internal/           # Client (client.go), <module>_service.go, typed HTTP errors (*_error.go)
str/                # JSON request/response types, one type per file
uri/                # query option structs (`url:` tags), AddQuery, ranges
```

Supporting packages: `cfg/` (TOML config and `str.Options`), `consts/` (named
constants and usage strings), `printer/` (all console output), `writer/`
(JSON file output), `buffer/`, `cli/` (OAuth device flow, token, version),
`test/` (shared test helpers), `testdata/` (fixtures), `docs/<module>.md`
(per-module user docs).

## Tooling

Run all of these before opening a PR; CI enforces them.

- `make cleanup` - `gofmt -w` on every `.go` file (CI fails on any diff).
- `make linter` - `revive --config ./revive.toml --formatter friendly ./...`.
- `go test -race -shuffle=on ./...` - same flags as CI (`make test` runs
  `-race` without shuffle).
- `go mod tidy` must leave `go.mod` and `go.sum` unchanged.
- `make install` (`go mod vendor`) after any dependency change.
- `make build` - builds with version ldflags.

## Restrictions

Hard limits for every contributor, human or agent. Ask before crossing one.

- **Never hand-edit `vendor/`.** Change `go.mod` and re-vendor.
- **No new dependencies without asking.**
- **Do not weaken checks to go green.** No loosening `revive.toml`, no
  `//revive:disable` comments, no skipped or deleted tests. Fix the code.
- **No secrets in the repo.** Never commit `trakt-sync.toml`, `token.json`,
  client IDs/secrets, or real user data.
- **Do not commit build or run output:** `*.json` exports, the `trakt-sync`
  binary, `dist/`, `coverage.out` (all gitignored - keep it that way).
- **Do not touch the release flow** (`.goreleaser.yaml`, tags, version ldflags)
  unless the task is about releasing.
- **Surgical diffs.** Every changed line traces to the request; do not reformat
  or "improve" adjacent code. Match existing style even if you would do it
  differently.

## Commit Standards

- Subject: `<module>: <what changed>`, lowercase, imperative, short - e.g.
  `watched: add pagination`, `sync: handle sync endpoints`. Use the CLI module
  name (or area such as `readme`, `linter`, `actions`). GitHub appends `(#N)` on
  squash merge.
- One module or feature per PR. Never commit directly to `main`; branch first
  and let CI gate the PR.
- No `Co-Authored-By` trailers and no "Generated with" footers.

## Pull Requests

- Short title, short description - a reviewer should grasp the PR in seconds.
- If an AI tool helped with code, tests or prose, say briefly what it did in the
  PR description. You remain responsible for every line.
- Review your own diff in the GitHub UI before asking for review; remove
  irrelevant generated changes.
- Bug fixes include a test that fails without the fix and passes with it.
- When a reviewer flags one occurrence of an issue, fix every similar
  occurrence in the PR.

## Documentation

- A new module or action updates `docs/<module>.md` and the command list in
  `README.md`.
- Keep examples in docs runnable against the current flags.

## Rule Files

Rules live in `.agents/rules/`; `AGENTS.md` loads the core ones and `CLAUDE.md`
routes to the rest. A new rule file needs a routing line in both. When you
establish a pattern that diverges from or extends these rules, update the
matching rule file in the same PR.

## Tone

Direct, concise, technical. State assumptions before non-trivial changes.
