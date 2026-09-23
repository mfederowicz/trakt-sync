---
description: 'Testing conventions for trakt-sync: libraries, mock HTTP server, filesystem, fixtures, CI flags. Read when writing or changing tests.'
applyTo: '**/*_test.go,test/**,testdata/**'
---

# Testing

## Basics

- Tests sit next to the code: `<file>_test.go`, same package.
- Use the standard `testing` package with `github.com/google/go-cmp/cmp` for
  deep comparison; `testify` is already a dependency and may be used where the
  surrounding tests use it. Do not add other assertion or mocking libraries.
- Prefer table-driven tests with `t.Run(name, ...)` subtests.
- Mark helpers with `t.Helper()`.
- Tests must pass with `go test -race -shuffle=on ./...`: no shared mutable
  state between tests, no order dependence.

## HTTP

- Never call the real Trakt API.
- Service tests use `internal.Setup()` (`internal/helpers.go`): an `httptest`
  server with a `Mux` to register handlers and a `Client` pointed at it; always
  `defer setup.Teardown()`.
- Handler tests follow `setup(t)` in `handlers/commons_test.go`.
- Assert the request with `test.AssertMethod` and check the path and query;
  use `test.AssertType` and `test.Ptr` from `test/helpers.go`.
- Cover the error paths that matter: non-2xx statuses mapping to the typed
  errors in `internal/*_error.go`, and pagination headers for `HavePages`.

## Filesystem and config

- Use `afero.NewMemMapFs()` for anything touching files (config, token,
  JSON output); never write to the real home directory.
- Fixtures go in `testdata/`; keep them small and free of real user data.

## Bug fixes

A bug fix ships with a test that fails without the fix and passes with it.
