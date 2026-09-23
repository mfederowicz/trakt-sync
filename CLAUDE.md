Before implementing anything, identify which area you are working in and read
the corresponding rule file from `.agents/rules/` (only the core rules auto-load
via AGENTS.md; domain rules load on demand):

- New or changed Trakt API endpoint, service method, handler, command action or
  `str/` type: read `endpoints.md` - contract -> Go mapping checklist.
- Tests (`*_test.go`, `test/`, `testdata/`): read `testing.md`.
- Everything else: `project.md` and `go-style.md` (always-on baseline, already
  loaded as core).

@AGENTS.md
