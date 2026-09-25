# API coverage

Trakt API routes (from the contract) and whether trakt-sync implements them.

- Source: [trakt/trakt-api](https://github.com/trakt/trakt-api) contracts, via the generated spec at <https://developer.trakt.tv/openapi.json>.
- Snapshot date: 2026-09-23. Update rows when the contract changes.
- ✅ implemented: a service method in `internal/` calls this exact route.
- 🟡 needs checking: a service method calls this route only through a generic path parameter (for example `sync/collection/%s`); confirm the CLI accepts this value.
- ⬜ missing: no service method calls this route.
- ⚠️ not served: listed in the contract, but the live API does not serve it (404, or another route answers; see [Findings](#findings)). The Go method column shows whether trakt-sync implements it anyway.
- ➖ not used: the CLI has no use for this route; the Go method column says why. Not work to pick.
- Update this file in the same PR that adds or removes an endpoint.

## Summary

| Domain | ✅ | 🟡 | ⬜ | ⚠️ | ➖ | Total |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| [`calendars`](#calendars) | 12 | 0 | 0 | 0 | 0 | 12 |
| [`certifications`](#certifications) | 3 | 0 | 0 | 0 | 0 | 3 |
| [`checkin`](#checkin) | 2 | 0 | 0 | 0 | 0 | 2 |
| [`comments`](#comments) | 18 | 0 | 0 | 0 | 0 | 18 |
| [`countries`](#countries) | 1 | 0 | 0 | 0 | 0 | 1 |
| [`episodes`](#episodes) | 2 | 0 | 0 | 0 | 0 | 2 |
| [`genres`](#genres) | 1 | 0 | 0 | 0 | 0 | 1 |
| [`languages`](#languages) | 1 | 0 | 0 | 0 | 0 | 1 |
| [`lists`](#lists) | 15 | 0 | 0 | 0 | 0 | 15 |
| [`media`](#media) | 3 | 0 | 0 | 0 | 0 | 3 |
| [`movies`](#movies) | 29 | 0 | 0 | 2 | 0 | 31 |
| [`networks`](#networks) | 1 | 0 | 0 | 0 | 0 | 1 |
| [`notes`](#notes) | 5 | 0 | 0 | 0 | 0 | 5 |
| [`oauth`](#oauth) | 3 | 0 | 0 | 0 | 2 | 5 |
| [`people`](#people) | 8 | 0 | 0 | 0 | 0 | 8 |
| [`recommendations`](#recommendations) | 4 | 0 | 0 | 0 | 0 | 4 |
| [`scrobble`](#scrobble) | 3 | 0 | 0 | 0 | 0 | 3 |
| [`search`](#search) | 6 | 0 | 0 | 0 | 0 | 6 |
| [`seasons`](#seasons) | 1 | 0 | 0 | 0 | 0 | 1 |
| [`shows`](#shows) | 58 | 0 | 0 | 2 | 0 | 60 |
| [`smart-lists`](#smart-lists) | 0 | 0 | 2 | 0 | 0 | 2 |
| [`social_recommendations`](#social_recommendations) | 2 | 0 | 0 | 0 | 0 | 2 |
| [`sync`](#sync) | 37 | 0 | 0 | 0 | 0 | 37 |
| [`team`](#team) | 1 | 0 | 0 | 0 | 0 | 1 |
| [`users`](#users) | 63 | 0 | 41 | 0 | 0 | 104 |
| [`watchnow`](#watchnow) | 2 | 0 | 0 | 0 | 0 | 2 |
| [`younify`](#younify) | 0 | 0 | 5 | 0 | 0 | 5 |
| **Total** | **281** | **0** | **48** | **4** | **2** | **335** |

## calendars

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/calendars/releases/hot/finales/{start_date}/{days}` | Get hot finales | `CalendarsService.GetHotFinales` |
| ✅ | GET | `/calendars/releases/hot/new/{start_date}/{days}` | Get hot new shows | `CalendarsService.GetHotNewShows` |
| ✅ | GET | `/calendars/releases/hot/premieres/{start_date}/{days}` | Get hot premieres | `CalendarsService.GetHotPremieres` |
| ✅ | GET | `/calendars/releases/hot/{start_date}/{days}` | Get hot releases | `CalendarsService.GetHotReleases` |
| ✅ | GET | `/calendars/{target}/dvd/{start_date}/{days}` | Get DVD releases | `CalendarsService.GetDVDReleases` |
| ✅ | GET | `/calendars/{target}/media/{start_date}/{days}` | Get media | `CalendarsService.GetMedia` |
| ✅ | GET | `/calendars/{target}/movies/{start_date}/{days}` | Get movies | `CalendarsService.GetMovies` |
| ✅ | GET | `/calendars/{target}/shows/finales/{start_date}/{days}` | Get finales | `CalendarsService.GetFinales` |
| ✅ | GET | `/calendars/{target}/shows/new/{start_date}/{days}` | Get new shows | `CalendarsService.GetNewShows` |
| ✅ | GET | `/calendars/{target}/shows/premieres/{start_date}/{days}` | Get season premieres | `CalendarsService.GetSeasonPremieres` |
| ✅ | GET | `/calendars/{target}/shows/{start_date}/{days}` | Get shows | `CalendarsService.GetShows` |
| ✅ | GET | `/calendars/{target}/streaming/{start_date}/{days}` | Get streaming releases | `CalendarsService.GetStreamingReleases` |

## certifications

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/certifications/movies` | Get movie certifications | `CertificationsService.GetCertifications` |
| ✅ | GET | `/certifications/shows` | Get show certifications | `CertificationsService.GetCertifications` |
| ✅ | GET | `/certifications/{type}` | Get certifications | `CertificationsService.GetCertifications` |

## checkin

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/checkin` | Check into an item | `CheckinService.CheckintoAnItem` |
| ✅ | DELETE | `/checkin` | Delete any active checkins | `CheckinService.DeleteAnyActiveCheckins` |

## comments

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/comments/` | Post a comment | `CommentsService.PostAComment` |
| ✅ | GET | `/comments/recent/{comment_type}/{type}` | Get recently created comments | `CommentsService.GetRecentComments` |
| ✅ | GET | `/comments/trending/{comment_type}/{type}` | Get trending comments | `CommentsService.GetTrendingComments` |
| ✅ | GET | `/comments/updates/{comment_type}/{type}` | Get recently updated comments | `CommentsService.GetUpdatedComments` |
| ✅ | GET | `/comments/{id}` | Get a comment or reply | `CommentsService.GetComment` |
| ✅ | PUT | `/comments/{id}/` | Update a comment or reply | `CommentsService.UpdateComment` |
| ✅ | DELETE | `/comments/{id}/` | Delete a comment or reply | `CommentsService.DeleteComment` |
| ✅ | GET | `/comments/{id}/item` | Get the attached media item | `CommentsService.GetCommentItem` |
| ✅ | POST | `/comments/{id}/like` | Like a comment | `CommentsService.LikeComment` |
| ✅ | DELETE | `/comments/{id}/like` | Remove like on a comment | `CommentsService.RemoveLikeComment` |
| ✅ | GET | `/comments/{id}/likes` | Get all users who liked a comment | `CommentsService.GetCommentUserLikes` |
| ✅ | GET | `/comments/{id}/reactions/` | Get comment reactions | `CommentsService.GetCommentReactions` |
| ✅ | GET | `/comments/{id}/reactions/summary` | Get reaction summary | `CommentsService.GetCommentReactionsSummary` |
| ✅ | POST | `/comments/{id}/reactions/{reaction_type}` | Add comment reaction | `CommentsService.AddCommentReaction` |
| ✅ | DELETE | `/comments/{id}/reactions/{reaction_type}` | Remove comment reaction | `CommentsService.RemoveCommentReaction` |
| ✅ | GET | `/comments/{id}/replies` | Get replies for a comment | `CommentsService.GetRepliesForComment` |
| ✅ | POST | `/comments/{id}/replies` | Post a reply for a comment | `CommentsService.ReplyAComment` |
| ✅ | POST | `/comments/{id}/report` | Report a comment | `CommentsService.ReportComment` |

## countries

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/countries/{type}` | Get countries | `CountriesService.GetCountries` |

## episodes

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/episodes/{id}/report` | Report an episode | `EpisodesService.ReportEpisode` |
| ✅ | GET | `/episodes/{id}/watchnow/{country}` | Get episode watch now sources | `EpisodesService.GetEpisodeWatchNow` |

## genres

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/genres/{type}` | Get genres | `GenresService.GetGenres` |

## languages

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/languages/{type}` | Get languages | `LanguagesService.GetLanguages` |

## lists

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/lists/popular` | Get popular lists | `ListsService.GetPopularLists` |
| ✅ | GET | `/lists/popular/{type}` | Get popular lists | `ListsService.GetPopularListsByType` |
| ✅ | GET | `/lists/trending` | Get trending lists | `ListsService.GetTrendingLists` |
| ✅ | GET | `/lists/trending/{type}` | Get trending lists | `ListsService.GetTrendingListsByType` |
| ✅ | GET | `/lists/{id}` | Get list | `ListsService.GetList` |
| ✅ | GET | `/lists/{id}/comments/{sort}` | Get all list comments | `ListsService.GetListComments` |
| ✅ | GET | `/lists/{id}/items/movie` | Get movie list items | `ListsService.GetListItems` |
| ✅ | GET | `/lists/{id}/items/movie,show` | Get media list items | `ListsService.GetListItems` |
| ✅ | GET | `/lists/{id}/items/movie,show,episode,season` | Get all list items | `ListsService.GetListItems` |
| ✅ | GET | `/lists/{id}/items/show` | Get show list items | `ListsService.GetListItems` |
| ✅ | GET | `/lists/{id}/items/{type}/{sort_by}/{sort_how}` | Get items on a list | `ListsService.GetListItems` |
| ✅ | POST | `/lists/{id}/like` | Like a list | `ListsService.LikeList` |
| ✅ | DELETE | `/lists/{id}/like` | Remove like on a list | `ListsService.RemoveLikeList` |
| ✅ | GET | `/lists/{id}/likes` | Get all users who liked a list | `ListsService.GetAllUsersWhoLikedList` |
| ✅ | POST | `/lists/{id}/report` | Report a list | `ListsService.ReportList` |

## media

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/media/anticipated` | Get anticipated media | `MediaService.GetAnticipatedMedia` |
| ✅ | GET | `/media/popular` | Get popular media | `MediaService.GetPopularMedia` |
| ✅ | GET | `/media/trending` | Get trending media | `MediaService.GetTrendingMedia` |

## movies

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/movies/anticipated` | Get the most anticipated movies | `MoviesService.GetAnticipatedMovies` |
| ✅ | GET | `/movies/boxoffice` | Get the weekend box office | `MoviesService.GetBoxoffice` |
| ✅ | GET | `/movies/collected/{period}` | Get the most collected movies | `MoviesService.GetCollectedMovies` |
| ✅ | GET | `/movies/favorited/{period}` | Get the most favorited movies | `MoviesService.GetFavoritedMovies` |
| ⚠️ | GET | `/movies/hot` | Get hot movies | `MoviesService.GetHotMovies` |
| ✅ | GET | `/movies/played/{period}` | Get the most played movies | `MoviesService.GetPlayedMovies` |
| ✅ | GET | `/movies/popular` | Get popular movies | `MoviesService.GetPopularMovies` |
| ⚠️ | GET | `/movies/streaming/{period}` | Get streaming movies | `MoviesService.GetStreamingMovies` |
| ✅ | GET | `/movies/trending` | Get trending movies | `MoviesService.GetTrendingMovies` |
| ✅ | GET | `/movies/updates/id/{start_date}` | Get recently updated movie Trakt IDs | `MoviesService.GetRecentlyUpdatedMoviesTraktIDs` |
| ✅ | GET | `/movies/updates/{start_date}` | Get recently updated movies | `MoviesService.GetRecentlyUpdatedMovies` |
| ✅ | GET | `/movies/watched/{period}` | Get the most watched movies | `MoviesService.GetWatchedMovies` |
| ✅ | GET | `/movies/{id}` | Get a movie | `MoviesService.GetMovie` |
| ✅ | GET | `/movies/{id}/aliases` | Get all movie aliases | `MoviesService.GetAllMovieAliases` |
| ✅ | GET | `/movies/{id}/comments/{sort}` | Get all movie comments | `MoviesService.GetAllMovieComments` |
| ✅ | GET | `/movies/{id}/lists/{type}/{sort}` | Get lists containing this movie | `MoviesService.GetListsContainingMovie` |
| ✅ | GET | `/movies/{id}/people` | Get all people for a movie | `MoviesService.GetAllPeopleForMovie` |
| ✅ | GET | `/movies/{id}/ratings` | Get movie ratings | `MoviesService.GetMovieRatings` |
| ✅ | POST | `/movies/{id}/refresh` | Refresh movie metadata | `MoviesService.RefreshMovieMetadata` |
| ✅ | POST | `/movies/{id}/refresh/justwatch` | Refresh movie JustWatch links | `MoviesService.RefreshMovieJustwatch` |
| ✅ | GET | `/movies/{id}/related` | Get related movies | `MoviesService.GetRelatedMovies` |
| ✅ | GET | `/movies/{id}/releases/{country}` | Get all movie releases | `MoviesService.GetAllMovieReleases` |
| ✅ | POST | `/movies/{id}/report` | Report a movie | `MoviesService.ReportMovie` |
| ✅ | GET | `/movies/{id}/sentiments` | Get movie sentiments | `MoviesService.GetMovieSentiments` |
| ✅ | GET | `/movies/{id}/stats` | Get movie stats | `MoviesService.GetMovieStats` |
| ✅ | GET | `/movies/{id}/studios` | Get movie studios | `MoviesService.GetMovieStudios` |
| ✅ | GET | `/movies/{id}/translations` | Get all movie translations | `MoviesService.GetAllMovieTranslations` |
| ✅ | GET | `/movies/{id}/videos` | Get all videos | `MoviesService.GetMovieVideos` |
| ✅ | GET | `/movies/{id}/watching` | Get users watching right now | `MoviesService.GetMovieWatching` |
| ✅ | GET | `/movies/{id}/watchnow/justwatch_links/{country}` | Get movie JustWatch links | `MoviesService.GetMovieJustwatchLinks` |
| ✅ | GET | `/movies/{id}/watchnow/{country}` | Get movie watch now sources | `MoviesService.GetMovieWatchNow` |

## networks

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/networks` | Get networks | `NetworksService.GetNetworksList` |

## notes

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/notes` | Add notes | `NotesService.AddNotes` |
| ✅ | GET | `/notes/{id}` | Get a note | `NotesService.GetNotes` |
| ✅ | PUT | `/notes/{id}` | Update a note | `NotesService.UpdateNotes` |
| ✅ | DELETE | `/notes/{id}` | Delete a note | `NotesService.DeleteNotes` |
| ✅ | GET | `/notes/{id}/item` | Get the attached item | `NotesService.GetNotesItem` |

## oauth

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ➖ | GET | `/oauth/authorize` | Authorize Application | not used: browser redirect flow for web apps; the CLI logs in with the device flow |
| ✅ | POST | `/oauth/device/code` | Generate new device codes | `OauthService.GenerateNewDeviceCodes` |
| ✅ | POST | `/oauth/device/token` | Poll for the access_token | `OauthService.PoolForTheAccessToken` |
| ➖ | POST | `/oauth/revoke` | Revoke an access_token | not used: the CLI refreshes an expired token or starts a new device login |
| ✅ | POST | `/oauth/token` | Exchange a token | `OauthService.ExchangeRefreshTokenForAccessToken` |

## people

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/people/updates/id/{start_date}` | Get recently updated people Trakt IDs | `PeopleService.GetRecentlyUpdatedPeopleTraktIDs` |
| ✅ | GET | `/people/updates/{start_date}` | Get recently updated people | `PeopleService.GetRecentlyUpdatedPeople` |
| ✅ | GET | `/people/{id}/` | Get a single person | `PeopleService.GetSinglePerson` |
| ✅ | GET | `/people/{id}/lists/{type}/{sort}` | Get lists containing this person | `PeopleService.GetListsContainingThisPerson` |
| ✅ | GET | `/people/{id}/movies` | Get movie credits | `PeopleService.GetMovieCredits` |
| ✅ | POST | `/people/{id}/refresh` | Refresh person metadata | `PeopleService.RefreshPersonMetadata` |
| ✅ | POST | `/people/{id}/report` | Report a person | `PeopleService.ReportPerson` |
| ✅ | GET | `/people/{id}/shows` | Get show credits | `PeopleService.GetShowCredits` |

## recommendations

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/recommendations/movies/` | Get movie recommendations | `RecommendationsService.GetMovieRecommendations` |
| ✅ | DELETE | `/recommendations/movies/{id}` | Hide a movie recommendation | `RecommendationsService.HideMovieRecommendation` |
| ✅ | GET | `/recommendations/shows/` | Get show recommendations | `RecommendationsService.GetShowRecommendations` |
| ✅ | DELETE | `/recommendations/shows/{id}` | Hide a show recommendation | `RecommendationsService.HideShowRecommendation` |

## scrobble

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/scrobble/pause` | Pause watching in a media center | `ScrobbleService.PauseScrobble` |
| ✅ | POST | `/scrobble/start` | Start watching in a media center | `ScrobbleService.StartScrobble` |
| ✅ | POST | `/scrobble/stop` | Stop or finish watching in a media center | `ScrobbleService.StopScrobble` |

## search

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/search/recent/` | Add recent search | `SearchService.AddRecentSearch` |
| ✅ | POST | `/search/recent/remove` | Remove recent search | `SearchService.RemoveRecentSearch` |
| ✅ | GET | `/search/recent_by_id/global/{type}` | Get trending search results | `SearchService.GetTrendingSearches` |
| ✅ | GET | `/search/{id_type}/{id}` | Get ID lookup results | `SearchService.GetIDLookupResults` |
| ✅ | GET | `/search/{type}` | Get text query results | `SearchService.GetTextQueryResults` |
| ✅ | GET | `/search/{type}/exact` | Get exact text query results | `SearchService.GetExactTextQueryResults` |

## seasons

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/seasons/{id}/report` | Report a season | `SeasonsService.ReportSeason` |

## shows

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/shows/anticipated` | Get the most anticipated shows | `ShowsService.GetAnticipatedShows` |
| ✅ | GET | `/shows/collected/{period}` | Get the most collected shows | `ShowsService.GetCollectedShows` |
| ✅ | GET | `/shows/favorited/{period}` | Get the most favorited shows | `ShowsService.GetFavoritedShows` |
| ⚠️ | GET | `/shows/hot` | Get hot shows |  |
| ✅ | GET | `/shows/played/{period}` | Get the most played shows | `ShowsService.GetPlayedShows` |
| ✅ | GET | `/shows/popular` | Get popular shows | `ShowsService.GetPopularShows` |
| ⚠️ | GET | `/shows/streaming/{period}` | Get streaming shows |  |
| ✅ | GET | `/shows/trending` | Get trending shows | `ShowsService.GetTrendingShows` |
| ✅ | GET | `/shows/updates/id/{start_date}` | Get recently updated show Trakt IDs | `ShowsService.GetRecentlyUpdatedShowsTraktIDs` |
| ✅ | GET | `/shows/updates/{start_date}` | Get recently updated shows | `ShowsService.GetRecentlyUpdatedShows` |
| ✅ | GET | `/shows/watched/{period}` | Get the most watched shows | `ShowsService.GetWatchedShows` |
| ✅ | GET | `/shows/{id}` | Get a single show | `ShowsService.GetShow` |
| ✅ | GET | `/shows/{id}/aliases` | Get all show aliases | `ShowsService.GetAllShowAliases` |
| ✅ | GET | `/shows/{id}/certifications` | Get all show certifications | `ShowsService.GetAllShowCertifications` |
| ✅ | GET | `/shows/{id}/comments/{sort}` | Get all show comments | `ShowsService.GetAllShowComments` |
| ✅ | GET | `/shows/{id}/last_episode` | Get last episode | `ShowsService.GetLastEpisode` |
| ✅ | GET | `/shows/{id}/lists/{type}/{sort}` | Get lists containing this show | `ShowsService.GetListsContainingShow` |
| ✅ | GET | `/shows/{id}/next_episode` | Get next episode | `ShowsService.GetNextEpisode` |
| ✅ | GET | `/shows/{id}/people` | Get all people for a show | `PeopleService.GetAllPeopleForShow` |
| ✅ | GET | `/shows/{id}/progress/collection` | Get show collection progress | `ShowsService.GetShowCollectionProgress` |
| ✅ | GET | `/shows/{id}/progress/watched` | Get show watched progress | `ShowsService.GetShowWatchedProgress` |
| ✅ | POST | `/shows/{id}/progress/watched/reset` | Reset show progress | `ShowsService.ResetShowProgress` |
| ✅ | DELETE | `/shows/{id}/progress/watched/reset` | Undo reset show progress | `ShowsService.UndoResetShowProgress` |
| ✅ | GET | `/shows/{id}/ratings` | Get show ratings | `ShowsService.GetShowRatings` |
| ✅ | POST | `/shows/{id}/refresh` | Refresh show metadata | `ShowsService.RefreshShowMetadata` |
| ✅ | POST | `/shows/{id}/refresh/justwatch` | Refresh show JustWatch links | `ShowsService.RefreshShowJustwatch` |
| ✅ | GET | `/shows/{id}/related` | Get related shows | `ShowsService.GetRelatedShows` |
| ✅ | POST | `/shows/{id}/report` | Report a show | `ShowsService.ReportShow` |
| ✅ | GET | `/shows/{id}/seasons` | Get all seasons for a show | `ShowsService.GetAllSeasonsForShow` |
| ✅ | GET | `/shows/{id}/seasons/{season}` | Get all episodes for a single season | `ShowsService.GetAllEpisodesForSingleSeason` |
| ✅ | GET | `/shows/{id}/seasons/{season}/comments/{sort}` | Get all season comments | `ShowsService.GetAllSeasonComments` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}` | Get a single episode for a show | `ShowsService.GetSingleEpisodeForShow` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/comments/{sort}` | Get all episode comments | `ShowsService.GetAllEpisodeComments` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/lists/{type}/{sort}` | Get lists containing this episode | `ShowsService.GetListsContainingEpisode` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/people` | Get all people for an episode | `ShowsService.GetAllPeopleForEpisode` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/ratings` | Get episode ratings | `ShowsService.GetEpisodeRatings` |
| ✅ | POST | `/shows/{id}/seasons/{season}/episodes/{episode}/report` | Report an episode | `ShowsService.ReportEpisode` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/stats` | Get episode stats | `ShowsService.GetEpisodeStats` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/translations` | Get all episode translations | `ShowsService.GetAllEpisodeTranslations` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/videos` | Get all videos | `ShowsService.GetEpisodeVideos` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/watching` | Get users watching right now | `ShowsService.GetEpisodesWatching` |
| ✅ | GET | `/shows/{id}/seasons/{season}/episodes/{episode}/watchnow/{country}` | Get episode watch now sources | `ShowsService.GetEpisodeWatchNow` |
| ✅ | GET | `/shows/{id}/seasons/{season}/info` | Get single seasons for a show | `ShowsService.GetSingleSeasonsForShow` |
| ✅ | GET | `/shows/{id}/seasons/{season}/lists/{type}/{sort}` | Get lists containing this season | `ShowsService.GetListsContainingSeason` |
| ✅ | GET | `/shows/{id}/seasons/{season}/people` | Get all people for a season | `ShowsService.GetAllPeopleForSeason` |
| ✅ | GET | `/shows/{id}/seasons/{season}/ratings` | Get season ratings | `ShowsService.GetSeasonRatings` |
| ✅ | POST | `/shows/{id}/seasons/{season}/report` | Report a season | `ShowsService.ReportSeason` |
| ✅ | GET | `/shows/{id}/seasons/{season}/stats` | Get season stats | `ShowsService.GetSeasonStats` |
| ✅ | GET | `/shows/{id}/seasons/{season}/translations` | Get all season translations | `ShowsService.GetAllSeasonTranslations` |
| ✅ | GET | `/shows/{id}/seasons/{season}/videos` | Get all videos | `ShowsService.GetSeasonsVideos` |
| ✅ | GET | `/shows/{id}/seasons/{season}/watching` | Get users watching right now | `ShowsService.GetSeasonsWatching` |
| ✅ | GET | `/shows/{id}/seasons/{season}/watchnow/justwatch_links/{country}` | Get season JustWatch links | `ShowsService.GetSeasonJustwatchLinks` |
| ✅ | GET | `/shows/{id}/sentiments` | Get show sentiments | `ShowsService.GetShowSentiments` |
| ✅ | GET | `/shows/{id}/stats` | Get show stats | `ShowsService.GetShowStats` |
| ✅ | GET | `/shows/{id}/studios` | Get show studios | `ShowsService.GetShowStudios` |
| ✅ | GET | `/shows/{id}/translations` | Get all show translations | `ShowsService.GetAllShowTranslations` |
| ✅ | GET | `/shows/{id}/videos` | Get all videos | `ShowsService.GetShowVideos` |
| ✅ | GET | `/shows/{id}/watching` | Get users watching right now | `ShowsService.GetShowWatching` |
| ✅ | GET | `/shows/{id}/watchnow/justwatch_links/{country}` | Get show JustWatch links | `ShowsService.GetShowJustwatchLinks` |
| ✅ | GET | `/shows/{id}/watchnow/{country}` | Get show watch now sources | `ShowsService.GetShowWatchNow` |

## smart-lists

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ⬜ | GET | `/smart-lists/{list_id}` | Get smart list |  |
| ⬜ | GET | `/smart-lists/{list_id}/items` | Get smart list items |  |

## social_recommendations

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/social_recommendations/movies/` | Get social movie recommendations | `SocialRecommendationsService.GetSocialMovieRecommendations` |
| ✅ | GET | `/social_recommendations/shows/` | Get social show recommendations | `SocialRecommendationsService.GetSocialShowRecommendations` |

## sync

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | POST | `/sync/collection` | Add items to collection | `SyncService.AddItemsToCollection` |
| ✅ | GET | `/sync/collection/episodes` | Get episode collection | `SyncService.GetCollection` |
| ✅ | GET | `/sync/collection/media` | Get media collection | `SyncService.GetCollection` |
| ✅ | GET | `/sync/collection/minimal/episodes` | Get minimal episode collection | `SyncService.GetMinimalCollection` |
| ✅ | GET | `/sync/collection/minimal/movies` | Get minimal movie collection | `SyncService.GetMinimalCollection` |
| ✅ | GET | `/sync/collection/minimal/shows` | Get minimal show collection | `SyncService.GetMinimalShowCollection` |
| ✅ | GET | `/sync/collection/movies` | Get movie collection | `SyncService.GetCollection` |
| ✅ | POST | `/sync/collection/remove` | Remove items from collection | `SyncService.RemoveItemsFromCollection` |
| ✅ | GET | `/sync/collection/shows` | Get show collection | `SyncService.GetCollection` |
| ✅ | GET | `/sync/collection/{type}` | Get collection | `SyncService.GetCollection` |
| ✅ | POST | `/sync/favorites` | Add items to favorites | `SyncService.AddItemsToFavorites` |
| ✅ | PUT | `/sync/favorites` | Update favorites | `SyncService.UpdateFavorites` |
| ✅ | POST | `/sync/favorites/remove` | Remove items from favorites | `SyncService.RemoveItemsFromFavorites` |
| ✅ | POST | `/sync/favorites/reorder` | Reorder favorited items | `SyncService.ReorderFavoritesItems` |
| ✅ | PUT | `/sync/favorites/{list_item_id}` | Update a favorite item | `SyncService.UpdateFavoriteItem` |
| ✅ | GET | `/sync/favorites/{type}/{sort_by}/{sort_how}` | Get favorites | `SyncService.GetFavorites` |
| ✅ | POST | `/sync/history` | Add items to watched history | `SyncService.AddItemsToHistory` |
| ✅ | POST | `/sync/history/remove` | Remove items from history | `SyncService.RemoveItemsFromHistory` |
| ✅ | GET | `/sync/history/{type}/{id}` | Get watched history | `SyncService.GetWatchedHistory` |
| ✅ | GET | `/sync/last_activities` | Get last activity | `SyncService.GetLastActivity` |
| ✅ | GET | `/sync/playback` | Get playback progress | `SyncService.GetPlaybackProgress` |
| ✅ | GET | `/sync/playback/episodes` | Get episode playback progress | `SyncService.GetPlaybackProgress` |
| ✅ | GET | `/sync/playback/movies` | Get movie playback progress | `SyncService.GetPlaybackProgress` |
| ✅ | DELETE | `/sync/playback/{id}` | Remove a playback item | `SyncService.RemovePlaybackItem` |
| ✅ | GET | `/sync/progress/up_next` | Get up next | `SyncService.GetUpNext` |
| ✅ | GET | `/sync/progress/up_next_nitro` | Get up next nitro | `SyncService.GetUpNextNitro` |
| ✅ | GET | `/sync/progress/watched` | Get watched progress | `SyncService.GetWatchedProgress` |
| ✅ | POST | `/sync/ratings` | Add new ratings | `SyncService.AddItemsToRatings` |
| ✅ | POST | `/sync/ratings/remove` | Remove ratings | `SyncService.RemoveItemsFromRatings` |
| ✅ | GET | `/sync/ratings/{type}/{rating}` | Get ratings | `SyncService.GetRatings` |
| ✅ | GET | `/sync/watched/{type}` | Get watched | `SyncService.GetWatched` |
| ✅ | POST | `/sync/watchlist` | Add items to watchlist | `SyncService.AddItemsToWatchlist` |
| ✅ | PUT | `/sync/watchlist` | Update watchlist | `SyncService.UpdateWatchlist` |
| ✅ | POST | `/sync/watchlist/remove` | Remove items from watchlist | `SyncService.RemoveItemsFromWatchlist` |
| ✅ | POST | `/sync/watchlist/reorder` | Reorder watchlist items | `SyncService.ReorderWatchlistItems` |
| ✅ | PUT | `/sync/watchlist/{list_item_id}` | Update a watchlist item | `SyncService.UpdateWatchlistItem` |
| ✅ | GET | `/sync/watchlist/{type}/{sort_by}/{sort_how}` | Get watchlist | `SyncService.GetWatchlist` |

## team

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/team/` | Get team members | `TeamService.GetTeamMembers` |

## users

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ⬜ | PUT | `/users/avatar` | Update avatar |  |
| ✅ | GET | `/users/blocked` | Get blocked users | `UsersService.GetBlockedUsers` |
| ✅ | POST | `/users/hidden/calendar/remove` | Remove hidden calendar items | `UsersService.RemoveHiddenItems` |
| ✅ | GET | `/users/hidden/dropped` | Get dropped shows | `UsersService.GetHiddenItems` |
| ✅ | GET | `/users/hidden/progress_watched` | Get hidden progress items | `UsersService.GetHiddenItems` |
| ✅ | POST | `/users/hidden/progress_watched/remove` | Remove hidden progress items | `UsersService.RemoveHiddenItems` |
| ✅ | GET | `/users/hidden/{section}` | Get hidden items | `UsersService.GetHiddenItems` |
| ✅ | POST | `/users/hidden/{section}` | Add hidden items | `UsersService.AddHiddenItems` |
| ✅ | POST | `/users/hidden/{section}/remove` | Remove hidden items | `UsersService.RemoveHiddenItems` |
| ⬜ | GET | `/users/reactions/comments` | Get comment reactions |  |
| ✅ | GET | `/users/requests/` | Get follow requests | `UsersService.GetFollowRequests` |
| ✅ | GET | `/users/requests/following` | Get pending following requests | `UsersService.GetPendingFollowingRequests` |
| ✅ | POST | `/users/requests/{id}` | Approve follow request | `UsersService.ApproveFollowRequest` |
| ✅ | DELETE | `/users/requests/{id}` | Deny follow request | `UsersService.DenyFollowRequest` |
| ⬜ | POST | `/users/saved_filters` | Add saved filters |  |
| ⬜ | DELETE | `/users/saved_filters/{id}` | Delete saved filter |  |
| ✅ | GET | `/users/saved_filters/{section}` | Get saved filters | `UsersService.GetSavedFilters` |
| ⬜ | PUT | `/users/set_cover` | Update cover image |  |
| ✅ | GET | `/users/settings` | Retrieve settings | `UsersService.RetrieveSettings` |
| ⬜ | PUT | `/users/settings` | Update settings |  |
| ⬜ | GET | `/users/settings/plex/` | Get Plex settings |  |
| ⬜ | PUT | `/users/settings/plex/` | Update Plex settings |  |
| ⬜ | POST | `/users/settings/plex/connect` | Connect Plex |  |
| ⬜ | DELETE | `/users/settings/plex/connect` | Disconnect Plex |  |
| ⬜ | GET | `/users/settings/plex/servers` | Get Plex servers |  |
| ⬜ | GET | `/users/settings/plex/servers/{server_id}` | Get Plex server accounts and libraries |  |
| ⬜ | POST | `/users/settings/plex/sync` | Sync Plex now |  |
| ⬜ | GET | `/users/syncs/` | Get data syncs |  |
| ⬜ | GET | `/users/syncs/{id}` | Get a data sync |  |
| ⬜ | DELETE | `/users/syncs/{id}` | Undo a data sync |  |
| ⬜ | GET | `/users/syncs/{id}/paused` | Get paused sync items |  |
| ⬜ | GET | `/users/syncs/{id}/skipped` | Get skipped sync items |  |
| ⬜ | GET | `/users/syncs/{type}` | Get data syncs by type |  |
| ✅ | GET | `/users/{id}/` | Get user profile | `UsersService.GetProfile`, `UsersService.GetUserProfile` |
| ✅ | POST | `/users/{id}/block` | Block this user | `UsersService.Block` |
| ✅ | DELETE | `/users/{id}/block` | Unblock this user | `UsersService.Unblock` |
| ✅ | GET | `/users/{id}/collection/{type}` | Get collection | `UsersService.GetCollection` |
| ✅ | GET | `/users/{id}/comments/{comment_type}/{type}` | Get comments | `UsersService.GetComments` |
| ✅ | GET | `/users/{id}/favorites/comments/{sort}` | Get all favorites comments | `UsersService.GetFavoritesComments` |
| ⬜ | GET | `/users/{id}/favorites/media/{sort}` | Get favorite media |  |
| ⬜ | GET | `/users/{id}/favorites/movies/{sort}` | Get favorite movies |  |
| ⬜ | GET | `/users/{id}/favorites/shows/{sort}` | Get favorite shows |  |
| ✅ | GET | `/users/{id}/favorites/{type}/{sort_by}/{sort_how}` | Get favorites | `UsersService.GetFavorites` |
| ✅ | POST | `/users/{id}/follow` | Follow this user | `UsersService.Follow` |
| ✅ | DELETE | `/users/{id}/follow` | Unfollow this user | `UsersService.Unfollow` |
| ✅ | GET | `/users/{id}/followers` | Get followers | `UsersService.GetFollowers` |
| ✅ | GET | `/users/{id}/following` | Get following | `UsersService.GetFollowing` |
| ✅ | GET | `/users/{id}/friends` | Get friends | `UsersService.GetFriends` |
| ✅ | GET | `/users/{id}/history/` | Get watched history | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/episodes` | Get episode watched history | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/episodes/{item_id}` | Get history for an episode | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/movies` | Get movie watched history | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/movies/{item_id}` | Get history for a movie | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/shows` | Get show watched history | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/shows/{item_id}` | Get history for a show | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/history/{type}/{item_id}` | Get watched history | `UsersService.GetHistory` |
| ✅ | GET | `/users/{id}/likes/{type}` | Get likes | `UsersService.GetLikes` |
| ✅ | GET | `/users/{id}/lists` | Get a user's personal lists | `UsersService.GetUsersPersonalLists` |
| ✅ | POST | `/users/{id}/lists` | Create personal list | `UsersService.AddPersonalList` |
| ✅ | GET | `/users/{id}/lists/collaborations` | Get all lists a user can collaborate on | `UsersService.GetCollaborations` |
| ✅ | POST | `/users/{id}/lists/reorder` | Reorder a user's lists | `UsersService.ReorderLists` |
| ✅ | GET | `/users/{id}/lists/{list_id}/` | Get personal list | `UsersService.GetList` |
| ✅ | PUT | `/users/{id}/lists/{list_id}/` | Update personal list | `UsersService.UpdateList` |
| ✅ | DELETE | `/users/{id}/lists/{list_id}/` | Delete a user's personal list | `UsersService.DeleteList` |
| ✅ | GET | `/users/{id}/lists/{list_id}/comments/{sort}` | Get all list comments | `UsersService.GetListComments` |
| ✅ | POST | `/users/{id}/lists/{list_id}/items` | Add items to personal list | `UsersService.AddListItems` |
| ⬜ | GET | `/users/{id}/lists/{list_id}/items/movie` | Get movie list items |  |
| ⬜ | GET | `/users/{id}/lists/{list_id}/items/movie,show` | Get media list items |  |
| ⬜ | GET | `/users/{id}/lists/{list_id}/items/movie,show,season,episode` | Get all list items |  |
| ✅ | POST | `/users/{id}/lists/{list_id}/items/remove` | Remove items from personal list | `UsersService.RemoveListItems` |
| ✅ | POST | `/users/{id}/lists/{list_id}/items/reorder` | Reorder items on a list | `UsersService.ReorderListItems` |
| ⬜ | GET | `/users/{id}/lists/{list_id}/items/show` | Get show list items |  |
| ✅ | PUT | `/users/{id}/lists/{list_id}/items/{list_item_id}` | Update a list item | `UsersService.UpdateListItem` |
| ✅ | GET | `/users/{id}/lists/{list_id}/items/{type}/{sort_by}/{sort_how}` | Get items on a personal list | `UsersService.GetItemstOnAPersonalList`, `UsersService.GetListItems` |
| ✅ | POST | `/users/{id}/lists/{list_id}/like` | Like a list | `UsersService.ListLike` |
| ✅ | DELETE | `/users/{id}/lists/{list_id}/like` | Remove like on a list | `UsersService.RemoveListLike` |
| ✅ | GET | `/users/{id}/lists/{list_id}/likes` | Get all users who liked a list | `UsersService.GetListLikes` |
| ⬜ | POST | `/users/{id}/lists/{list_id}/reorder` | Reorder items on a list |  |
| ✅ | POST | `/users/{id}/lists/{list_id}/report` | Report a user's list | `UsersService.ListReport` |
| ⬜ | GET | `/users/{id}/mir/{year}/{month}` | Get month in review |  |
| ✅ | GET | `/users/{id}/notes/{type}` | Get notes | `UsersService.GetNotes` |
| ✅ | GET | `/users/{id}/ratings/` | Get all ratings | `UsersService.GetRatings` |
| ⬜ | GET | `/users/{id}/ratings/episodes` | Get episode ratings |  |
| ⬜ | GET | `/users/{id}/ratings/movies` | Get movie ratings |  |
| ⬜ | GET | `/users/{id}/ratings/shows` | Get show ratings |  |
| ✅ | GET | `/users/{id}/ratings/{type}/{rating}` | Get ratings | `UsersService.GetRatings` |
| ✅ | POST | `/users/{id}/report` | Report a user | `UsersService.Report` |
| ⬜ | GET | `/users/{id}/smart-lists` | Get a user's smart lists |  |
| ⬜ | POST | `/users/{id}/smart-lists` | Create smart list |  |
| ⬜ | GET | `/users/{id}/smart-lists/{list_id}/` | Get smart list |  |
| ⬜ | PUT | `/users/{id}/smart-lists/{list_id}/` | Update smart list |  |
| ⬜ | DELETE | `/users/{id}/smart-lists/{list_id}/` | Delete a user's smart list |  |
| ✅ | GET | `/users/{id}/stats` | Get stats | `UsersService.GetStats` |
| ✅ | GET | `/users/{id}/watched/movies` | Get watched movies | `UsersService.GetWatched` |
| ✅ | GET | `/users/{id}/watched/shows` | Get watched shows | `UsersService.GetWatched` |
| ✅ | GET | `/users/{id}/watched/{type}` | Get watched | `UsersService.GetWatched` |
| ✅ | GET | `/users/{id}/watching` | Get watching | `UsersService.Watching` |
| ✅ | GET | `/users/{id}/watchlist/comments/{sort}` | Get all watchlist comments | `UsersService.GetWatchlistComments` |
| ⬜ | GET | `/users/{id}/watchlist/movie,show/{sort}` | Get media watchlist |  |
| ⬜ | GET | `/users/{id}/watchlist/movies/{sort}` | Get movie watchlist |  |
| ⬜ | GET | `/users/{id}/watchlist/shows/{sort}` | Get show watchlist |  |
| ✅ | GET | `/users/{id}/watchlist/{type}/{sort_by}/{sort_how}` | Get watchlist | `UsersService.GetWatchlist` |
| ⬜ | GET | `/users/{id}/yir/{year}` | Get year in review |  |
| ⬜ | GET | `/users/{id}/{type}/activities` | Get social activity |  |

## watchnow

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ✅ | GET | `/watchnow/sources` | Get watch now sources | `WatchNowService.GetWatchNowSources` |
| ✅ | GET | `/watchnow/sources/{countryCode}` | Get watch now sources by country | `WatchNowService.GetWatchNowSourcesByCountry` |

## younify

| Status | Method | Path | Summary | Go method |
| :---: | --- | --- | --- | --- |
| ⬜ | POST | `/younify/connect` | Create a streaming connection |  |
| ⬜ | GET | `/younify/connections` | Get streaming connections |  |
| ⬜ | POST | `/younify/users/refresh/{service_id}` | Refresh a streaming service |  |
| ⬜ | POST | `/younify/users/refresh/{service_id}/{all_data}` | Refresh a streaming service (full re-sync) |  |
| ⬜ | DELETE | `/younify/users/services/{service_id}` | Unlink a streaming service |  |

## Findings

Differences between the service code and the contract, found while building this map. Each one should be checked and fixed in its own small PR.

| Go method | Route | Note |
| --- | --- | --- |
| `MoviesService.GetHotMovies` | `GET /movies/hot` | live API returns 404 (checked 2026-09-24); the CLI explains the 404. Upstream issue: TBD |
| - | `GET /shows/hot` | live API routes it to `GET /shows/{id}` and returns the show with slug `hot` (checked 2026-09-24). Upstream issue: TBD |
| `MoviesService.GetStreamingMovies` | `GET /movies/streaming/{period}` | live API returns 404 `{"error":"endpoint removed"}` (checked 2026-09-24); the CLI explains the 404. Upstream issue: TBD |
| - | `GET /shows/streaming/{period}` | live API returns 404 `{"error":"endpoint removed"}` (checked 2026-09-24). Upstream issue: TBD |
