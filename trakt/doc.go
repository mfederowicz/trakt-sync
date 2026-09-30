/*
Package trakt is a Go client for the Trakt API (https://trakt.tv).

It is the API layer of the trakt-sync CLI, published as a library. The package is experimental: its API may
still change in minor releases until it is declared stable.

	import "github.com/mfederowicz/trakt-sync/trakt"

Request and response types live in the str package, and query options in the uri package.

# Creating a client

Every request needs the client id of your Trakt API app. Set a User-Agent that names your app:

	client := trakt.NewClient(nil).
		WithClientID("your-client-id").
		WithUserAgent("my-app/1.0")

	movies, _, err := client.Movies.GetTrendingMovies(ctx, &uri.ListOptions{Limit: 10})

NewClient takes an optional *http.Client, for example one with a timeout. WithClientID, WithAuthToken and
WithUserAgent return a copy of the client; the original is never modified, so one base client can hand out
variants (for example one per user token).

The API is split into services that follow the Trakt API sections: client.Movies, client.Shows, client.Sync,
client.Users, client.Search and so on. Methods take a context.Context first and return the decoded result, a
*str.Response and an error.

# Authentication

Public data (movies, shows, people, trending lists) needs only the client id. Anything that reads or changes a
user's account needs an OAuth access token:

	userClient := client.WithAuthToken(accessToken)

Command line and TV-style apps get a token with the device flow: client.Oauth.GenerateNewDeviceCodes, then poll
client.Oauth.PollForAccessToken until the user approves the code. See example/devicecode in the repository.
Tokens expire; exchange the refresh token with client.Oauth.ExchangeRefreshTokenForAccessToken.
The app's client secret is needed only by these two token calls and goes in their request body
(str.NewDeviceToken, str.CurrentDeviceToken); the client never stores it or sends it on other requests.

# Pagination

List methods take a *uri.ListOptions with Page and Limit. The pagination headers of the response tell whether
more pages exist:

	opts := &uri.ListOptions{Page: 1, Limit: 100}
	for {
		items, resp, err := client.Sync.GetWatchedHistory(ctx, &noID, nil, opts)
		// ... use items, handle err
		if !client.HavePages(opts.Page, resp, 0) {
			break
		}
		opts.Page++
	}

The last argument of HavePages caps the number of pages; 0 means no cap.

# Errors

Every non-2xx response is returned as an error. The common statuses have their own types, which errors.As can
pick out: *BadRequestError (400), *InvalidUserError (401), *ForbiddenError (403), *NotFoundError (404),
*ConflictError (409), *PreconditionFailedRequestError (412), *UpgradeUserLimitsError (420), *ValidationError
(422), *UpgradeRequiredError (426, Trakt VIP only), *AbuseRateLimitError (429) and *ServerError (500). Other
statuses come as *str.ErrorResponse. The error message has the method, the URL, the status and the API message.
See example/errors in the repository.

# Rate limits

The rate limit of the last response is in resp.Rate. After a 429 the client returns *AbuseRateLimitError, with
RetryAfter when Trakt sends it, and refuses later requests locally until Client.RateLimitReset, so it does not
hammer the API.

# Timezones

Timestamps in responses are UTC. To get them in a local timezone, put it in the context:

	ctx := trakt.WithTimezone(context.Background(), time.Local)

# Debugging

The client never prints. Set DebugLogger to see each request (method and URL, never headers or tokens) and
the services' debug notes:

	client.DebugLogger = func(v ...any) { log.Println(v...) }

# API use

Every app has to follow the Trakt API Use Policy:
https://developer.trakt.tv/?section=guides&guide=api-use-policy
*/
package trakt
