#### Younify (streaming service connections):
Connect streaming services (Netflix, Hulu, ...) to your Trakt account so Trakt can sync from them. All actions need OAuth; some services need Trakt VIP (`vip: true` / `connectable: false` in `connections`).

> **Note:** as of 2026-09-25 Trakt answers `younify -a connections` with 401 for API apps (the developer portal too), and the other
> actions are expected to do the same; they fail with `younify is not open to API apps yet`. They are kept for when Trakt opens younify to API apps.

##### Streaming services and your connection status:
```console
$ ./trakt-sync younify -a connections -> export_younify_connections.json
```
##### Connect a service:
Prints a signed web auth URL; open it in a browser to link the service. `-return_url` is where the web auth returns afterwards and must be `trakt://...` or `https://*.trakt.tv` (default `https://trakt.tv`).
```console
$ ./trakt-sync younify -a connect -service_id netflix -> export_younify_connect_results.json
```
##### Queue a re-sync of a connected service:
```console
$ ./trakt-sync younify -a refresh -service_id netflix
```
Full re-sync of all data instead of an incremental one:
```console
$ ./trakt-sync younify -a refresh -service_id netflix -all_data
```
##### Unlink a service:
```console
$ ./trakt-sync younify -a disconnect -service_id netflix
```
