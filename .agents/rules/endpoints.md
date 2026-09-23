---
description: 'Checklist for adding or changing a Trakt API endpoint: contract -> str type -> uri options -> service -> handler -> command -> docs -> tests. Read when touching endpoints.'
applyTo: 'internal/**,handlers/**,cmds/**,str/**,uri/**'
---

# Endpoints

Every endpoint follows the same chain. Worked reference: `networks`
(`cmds/command_networks.go` -> `handlers/networks_lists_handler.go` ->
`internal/networks_service.go` -> `str.TvNetwork`), backed by the contract
`projects/api/src/contracts/networks/index.ts` in trakt/trakt-api.

## 1. Read the contract

Open `projects/api/src/contracts/<domain>/` in
[trakt/trakt-api](https://github.com/trakt/trakt-api) (`index.ts` plus its
`schema/request` and `schema/response` files). Note for each route:

- `method` and `path` (joined with the router's `pathPrefix`)
- `pathParams`, `query` (pagination, `extended`, filters), `body`
- `responses` by status code

Contract domain names map to CLI modules and services (`networks` ->
`NetworksService` / `networks` command). Only the contract shape matters; ignore
the TypeScript implementation details.

## 2. Types in `str/`

- One type per file, snake_case file name, `String()` via `Stringify`.
- Reuse existing types (`Movie`, `Show`, `Episode`, `Ids`, ...) before adding
  new ones.
- Zod -> Go field mapping:

| Zod                                     | Go                                    |
| --------------------------------------- | ------------------------------------- |
| `z.string()`                            | `*string` + `omitempty` (or `string`) |
| `z.number().int()`                      | `*int` / `*int64`                     |
| `z.number()`                            | `*float32` / `*float64`               |
| `z.boolean()`                           | `*bool`                               |
| `.nullable()` / `.optional()` / `.nullish()` | pointer + `omitempty` (required)  |
| `z.array(x)`                            | `*[]T` or `[]*T`, matching neighbours |
| datetime strings                        | `*Timestamp` / existing date types     |

- JSON tag names are exactly the contract keys (snake_case).
- Polymorphic responses (a list mixing movies/shows/episodes, or shape chosen
  by a `type` param) are **one flat struct** with every variant-specific field
  optional - never one struct per variant. The caller checks `x.Movie != nil`.

## 3. Query options in `uri/`

- Reuse `uri.ListOptions` / `uri.Pagination` for `page`, `limit`, `extended`.
- Endpoint-specific params get an options struct with `url:"...,omitempty"`
  tags; ranges reuse the existing range types.

## 4. Service method in `internal/<module>_service.go`

```go
// GetNetworksList Get a list of all TV networks, including the name, country, and ids.
func (m *NetworksService) GetNetworksList(ctx context.Context, opts *uri.ListOptions) ([]*str.TvNetwork, *str.Response, error) {
	var url = "networks"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	list := []*str.TvNetwork{}
	resp, err := m.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}
	return list, resp, nil
}
```

- Signature: `(ctx context.Context, <path params>, opts *uri.X[, body]) (T, *str.Response, error)`.
- Use `http.Method*` constants; path built with `fmt.Sprintf` from path params.
- Check `err` before touching `resp`.
- Add `uri.AddQuery` to the URL before `NewRequest`; otherwise the query is dropped.
- Doc comment links to the API reference:
  `// API docs: https://docs.trakt.tv/reference/<operationid>`, where
  `<operationid>` is the route's `operationId` from
  <https://developer.trakt.tv/openapi.json> in lowercase (e.g.
  `getUsersRequestsFollow` -> `getusersrequestsfollow`). Check that the page
  opens; unknown IDs return 404. When
  you touch a method that still links to `trakt.docs.apiary.io`, replace that
  link. Change only the methods the PR already touches, one at a time, never in bulk.
- A new service: declare `type XService Service` in its own file, add the field
  to `Client` and the `c.X = (*XService)(&c.common)` line in `client.go`.

## 5. Handler in `handlers/<module>_<action>_handler.go`

- `type XYHandler struct{}` with
  `Handle(options *str.Options, client *internal.Client) error`.
- Validate required options first and return clear errors.
- Paginated endpoints recurse with `client.HavePages(page, resp, options.PagesLimit)`,
  sleeping `consts.SleepNumberOfSeconds` between pages and advancing by
  `consts.NextPageStep` (see `NetworksListsHandler.fetchNetworksList`).
- Output: `json.MarshalIndent(..., consts.EmptyString, consts.JSONDataFormat)`
  then `writer.WriteJSON(options, data)`.
- Shared logic belongs in `CommonLogic` (`handlers/commons.go`), not copied
  between handlers.

## 6. Command in `cmds/command_<module>.go`

- Add the action to the command's `map[string]handlers.Handler` and to the
  action list passed to `cmd.common.GenActionsUsage`.
- New flags use `consts` usage strings and `cfg.DefaultConfig()` defaults.
- A new module: new `XCmd` with `Name`, `Summary`, `Help`, `Run` set in
  `init()`, and an entry in the `cmds/runtime.go` command list (kept A-Z).

## 7. Docs

Update `docs/<module>.md` (actions, flags, example invocations) and the
command list in `README.md`. Mark the route ✅ in `API_COVERAGE.md`, with the Go
method, and update that domain's counts in the summary table. Pick new work
from its ⬜ rows.

## 8. Tests

See `testing.md`. At minimum: a service test against a mock server asserting
method, path, query and decoded result.
