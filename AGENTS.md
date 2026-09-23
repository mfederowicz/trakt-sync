# AGENTS.md

Guidance for AI coding agents (and humans) working in this repository.

`trakt-sync` is an open-source Go CLI and API client for [Trakt](https://trakt.tv)
(module `github.com/mfederowicz/trakt-sync`). Each CLI module (`movies`, `sync`,
`users`, ...) maps to one domain of the Trakt API and exports results as JSON.

## API source of truth

The Trakt API contract lives in [trakt/trakt-api](https://github.com/trakt/trakt-api):
`projects/api/src/contracts/<domain>/` (ts-rest routers + Zod schemas). Use it for
endpoint paths, methods, path/query params, request bodies and response shapes;
use the developer portal <https://developer.trakt.tv> for guides. If the contract and our code disagree,
follow the contract and mention it in the PR.

Only the contracts are relevant here. trakt-api is a Deno/TypeScript project;
its code style, tooling and commit rules do not apply to this Go repo.

# Always-loaded core

@.agents/rules/project.md

@.agents/rules/go-style.md

# Domain rules - load on demand

Read these with the Read tool when the work enters the domain (CLAUDE.md routes
the mapping; files live in `.agents/rules/`):

- `endpoints.md` - adding or changing an endpoint: contract -> `str/` type ->
  `uri/` options -> `internal/` service -> `handlers/` -> `cmds/` -> `docs/`.
- `testing.md` - writing or changing tests, test helpers and fixtures.

Re-read after long gaps if context was compacted.
