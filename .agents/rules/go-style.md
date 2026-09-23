---
description: 'Go code conventions for trakt-sync: comments, naming, types, constants, errors, context, output. Apply to all Go code.'
applyTo: '**/*.go'
---

# Go Style

Baseline: [Effective Go](https://go.dev/doc/effective_go), [Go Code Review
Comments](https://go.dev/wiki/CodeReviewComments) and the rules enabled in
`revive.toml`. The points below are what this repo adds or enforces.

## Comments

- Every file starts with the package comment already used by its package, e.g.
  `// Package handlers used to handle module actions` (`package-comments`).
- Every exported identifier has a doc comment starting with its name
  (`exported`):

```go
// GetNetworksList Get a list of all TV networks, including the name, country, and ids.
func (m *NetworksService) GetNetworksList(ctx context.Context, opts *uri.ListOptions) ([]*str.TvNetwork, *str.Response, error) {
```

- Do not repeat what the Trakt docs already say; a one-line summary is enough.
- `// comment` with a space after the slashes (`comment-spacings`).

## Files

- File names are lowercase snake_case: `^[_a-z][_a-z0-9]*\.go$`
  (`filename-format`).
- Naming patterns: `cmds/command_<module>.go`,
  `handlers/<module>_<action>_handler.go`, `internal/<module>_service.go`,
  `internal/<status>_error.go`, `str/<type_name>.go` (one type per file).

## Naming

- Receivers are short and consistent within a type (`receiver-naming`); use `_`
  or omit the name when unused (`unused-receiver`).
- Service methods state the action and the resource: `Get...`, `Add...`,
  `Remove...`, `Update...`, `Delete...` (e.g. `GetNetworksList`).
- Common variable names: `ctx` (context), `url` (endpoint path), `opts` (query
  options), `req`, `resp`, `err`, `list` / `result` (decoded data), `options`
  (`*str.Options`), `client` (`*internal.Client`).
- Initialisms keep their case: `ID`, `URL`, `JSON`, `IMDB` (`var-naming`).

## Types

- Optional API fields are pointers with `omitempty`:
  ``Name *string `json:"name,omitempty"` ``. Required fields may be values.
- `str` types get a `String()` method via `Stringify`:

```go
func (n Network) String() string {
	return Stringify(n)
}
```

- Query parameters live in `uri` option structs with `url:"...,omitempty"`
  tags and are applied with `uri.AddQuery`; never concatenate query strings.
- Use `any`, not `interface{}` (`use-any`).
- Maps and slices use literals: `map[string]handlers.Handler{}`, `[]*str.X{}`,
  not `make(...)` for empty values (`enforce-map-style`, `enforce-slice-style`).

## Constants

No magic numbers or strings in logic. Use or extend `consts/`
(`consts.ZeroValue`, `consts.DefaultPage`, `consts.NextPageStep`,
`consts.EmptyString`, `consts.JSONDataFormat`, usage strings in `glob.go`).
Keep each `const` block sorted as it is.

## Errors

- Return early; keep the happy path unindented (`indent-error-flow`,
  `superfluous-else`).
- Wrap with context: `fmt.Errorf("fetch lists error: %w", err)`; use
  `errors.New` when nothing is formatted, never `errors.New(fmt.Sprintf(...))`
  (`errorf`).
- Error strings are lowercase, no trailing punctuation (`error-strings`).
- Error types are named `XxxError`, sentinel values `errXxx` / `ErrXxx`
  (`error-naming`); HTTP status errors live in `internal/*_error.go`.
- `error` is always the last return value (`error-return`).
- Never ignore an error silently in new code; if it is truly safe to ignore,
  say why in a comment.

## Context

- `context.Context` is the first parameter of every service method
  (`context-as-argument`); never store it in a struct.
- Handlers build it with `client.BuildCtxFromOptions(options)` so the user's
  timezone is carried along.
- Context keys use the typed `contextKey`, never plain strings
  (`context-keys-type`).

## Output

- Console output goes through `printer` (`printer.Println`, `printer.Printf`),
  never `fmt.Print*` or the builtin `print`.
- JSON files are written through `writer.WriteJSON(options, data)`.

## Formatting and limits

- `gofmt` clean; imports grouped stdlib, then module, then third-party.
- Max line length 255.
- No `os.Exit` / `log.Fatal` outside `main` (`deep-exit`); return errors.
- No dot imports, no blank imports outside `main`/tests (`dot-imports`,
  `blank-imports`).
- No unused parameters - rename to `_` (`unused-parameter`).
- No bare returns, no empty blocks, no unreachable code, no identical branches.
