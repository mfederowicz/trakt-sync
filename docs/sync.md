#### Sync:
##### Last activities
```console
$ ./trakt-sync sync -a last_activities
```
##### Playback progress from last 60 days
```console
$ ./trakt-sync sync -a playback
```
##### Playback progress from last 60 days (movies or episodes)
```console
$ ./trakt-sync sync -a playback -t movies
```
##### Playback progress from 7 days
```console
$ ./trakt-sync sync -a playback -start_at 2025-06-01 -end_at 2025-06-07
```
##### Remove playback
```console
$ ./trakt-sync sync -a remove_playback -playback_id 12345
```
##### Get collection - movies
```console
$ ./trakt-sync sync -a get_collection -t movies -ex metadata
```
##### Get collection - shows
```console
$ ./trakt-sync sync -a get_collection -t shows -ex metadata
```
##### Get collection - episodes
```console
$ ./trakt-sync sync -a get_collection -t episodes -ex metadata
```
##### Get collection - seasons
```console
$ ./trakt-sync sync -a get_collection -t seasons -ex metadata
```
##### Add to collection - via -items flag
```console
$ ./trakt-sync sync -t movies -a add_to_collection -items export_sync_collection_movies.json
```
```console
$ ./trakt-sync sync -t shows -a add_to_collection -items export_sync_collection_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a add_to_collection -items export_sync_collection_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a add_to_collection -items export_sync_collection_seasons.json
```
##### Add to collection - via stdin
```console
$ cat export_sync_collection_movies.json | ./trakt-sync sync -t movies -a add_to_collection
```
```console
$ cat export_sync_collection_shows.json | ./trakt-sync sync -t shows -a add_to_collection
```
```console
$ cat export_sync_collection_episodes.json | ./trakt-sync sync -t episodes -a add_to_collection
```
```console
$ cat export_sync_collection_seasons.json | ./trakt-sync sync -t seasons -a add_to_collection
```
##### Remove from collection - via -items flag
```console
$ ./trakt-sync sync -t movies -a remove_from_collection -items export_sync_collection_movies.json
```
```console
$ ./trakt-sync sync -t shows -a remove_from_collection -items export_sync_collection_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a remove_from_collection -items export_sync_collection_episode.json
```
```console
$ ./trakt-sync sync -t seasons -a remove_from_collection -items export_sync_collection_seasons.json
```
##### Remove from collection - via stdin
```console
$ cat export_sync_collection_movies.json | ./trakt-sync sync -t movies -a remove_from_collection
```
```console
$ cat export_sync_collection_shows.json | ./trakt-sync sync -t shows -a remove_from_collection
```
```console
$ cat export_sync_collection_episodes.json | ./trakt-sync sync -t episodes -a remove_from_collection
```
```console
$ cat export_sync_collection_seasons.json | ./trakt-sync sync -t seasons -a remove_from_collection
```
##### Get watched - movies
```console
$ ./trakt-sync sync -a get_watched -t movies
```
##### Get watched - shows
```console
$ ./trakt-sync sync -a get_watched -t shows
```
##### Get watched - shows - noseasons
```console
$ ./trakt-sync sync -a get_watched -t shows -ex noseasons
```
##### Get watched - episodes
```console
$ ./trakt-sync sync -a get_watched -t episodes
```
##### Get history - movies
```console
$ ./trakt-sync sync -a get_history -t movies -start_at 2025-07-01 -end_at 2025-07-06
```
##### Get history - shows - with trakt_id
```console
$ ./trakt-sync sync -a get_history -t shows -i 1388
```
##### Get history - shows
```console
$ ./trakt-sync sync -a get_history -t shows
```
##### Get history - episodes
```console
$ ./trakt-sync sync -a get_history -t episodes
```
##### Remove and Add to history - via -items flag
```console
$ ./trakt-sync sync -t movies -a add_to_history -items export_sync_history_movies.json
```
```console
$ ./trakt-sync sync -t shows -a add_to_history -items export_sync_history_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a add_to_history -items export_sync_history_episode.json
```
```console
$ ./trakt-sync sync -t seasons -a add_to_history -items export_sync_history_seasons.json
```
##### Remove and Add to history - via stdin
```console
$ cat export_sync_history_movies.json | ./trakt-sync sync -t movies -a add_to_history
```
```console
$ cat export_sync_history_shows.json | ./trakt-sync sync -t shows -a add_to_history
```
```console
$ cat export_sync_history_episodes.json | ./trakt-sync sync -t episodes -a add_to_history
```
```console
$ cat export_sync_history_seasons.json | ./trakt-sync sync -t seasons -a add_to_history
```
##### Remove from history - via -items flag
```console
$ ./trakt-sync sync -t movies -a remove_from_history -items export_sync_history_movies.json
```
```console
$ ./trakt-sync sync -t shows -a remove_from_history -items export_sync_history_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a remove_from_history -items export_sync_history_episode.json
```
```console
$ ./trakt-sync sync -t seasons -a remove_from_history -items export_sync_history_seasons.json
```
##### Remove from history - via stdin
```console
$ cat export_sync_history_movies.json | ./trakt-sync sync -t movies -a remove_from_history
```
```console
$ cat export_sync_history_shows.json | ./trakt-sync sync -t shows -a remove_from_history
```
```console
$ cat export_sync_history_episodes.json | ./trakt-sync sync -t episodes -a remove_from_history
```
```console
$ cat export_sync_history_seasons.json | ./trakt-sync sync -t seasons -a remove_from_history
```
##### Get ratings - movies - all ratings
```console
$ ./trakt-sync sync -a get_ratings -t movies
```
##### Get ratings - shows - filter for specific rating from 1 to 10
```console
$ ./trakt-sync sync -a get_ratings -t shows -rating 1,2,3
```
##### Get ratings - seasons
```console
$ ./trakt-sync sync -a get_ratings -t seasons
```
##### Get ratings - episodes
```console
$ ./trakt-sync sync -a get_ratings -t episodes
```
##### Get ratings - all ratings
```console
$ ./trakt-sync sync -a get_ratings -t all
```
##### Add ratings - via -items flag
```console
$ ./trakt-sync sync -t movies -a add_to_ratings -items export_sync_ratings_movies.json
```
```console
$ ./trakt-sync sync -t shows -a add_to_ratings -items export_sync_ratings_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a add_to_ratings -items export_sync_ratings_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a add_to_ratings -items export_sync_ratings_seasons.json
```
##### Add ratings - via stdin
```console
$ cat export_sync_ratings_movies.json |
```
```console
$ cat export_sync_ratings_shows.json | ./trakt-sync sync -t shows -a add_to_ratings
```
```console
$ cat export_sync_ratings_episodes.json | ./trakt-sync sync -t episodes -a add_to_ratings
```
```console
$ cat export_sync_ratings_seasons.json | ./trakt-sync sync -t seasons -a add_to_ratings
```
```console
$ cat export_sync_ratings_all.json | ./trakt-sync sync -t seasons -a add_to_ratings
```
##### Remove from ratings - via -items flag
```console
$ ./trakt-sync sync -t movies -a remove_from_ratings -items export_sync_ratings_movies.json
```
```console
$ ./trakt-sync sync -t shows -a remove_from_ratings -items export_sync_ratings_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a remove_from_ratings -items export_sync_ratings_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a remove_from_ratings -items export_sync_ratings_seasons.json
```
```console
$ ./trakt-sync sync -t seasons -a remove_from_ratings -items export_sync_ratings_all.json
```
##### Remove from ratings - via stdin
```console
$ cat export_sync_ratings_movies.json | ./trakt-sync sync -t movies -a remove_from_ratings
```
```console
$ cat export_sync_ratings_shows.json | ./trakt-sync sync -t shows -a remove_from_ratings
```
```console
$ cat export_sync_ratings_episodes.json | ./trakt-sync sync -t episodes -a remove_from_ratings
```
```console
$ cat export_sync_ratings_seasons.json | ./trakt-sync sync -t seasons -a remove_from_ratings
```
```console
$ cat export_sync_ratings_all.json | ./trakt-sync sync -t seasons -a remove_from_ratings
```
##### Get Watchlist
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_how asc
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_how des
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by rank
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by added
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by title
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by released
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by runtime
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by popularity
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by random
```
```console
$ ./trakt-sync sync -a get_watchlist -t movies -sort_by percentage
```

```console
🔥VIP Only including imdb_rating, tmdb_rating, rt_tomatometer, rt_audience, metascore, votes,
imdb_votes, and tmdb_votes. If sent for a non VIP, the items will fall back to rank.
```
##### Update Watchlist
```console
$ ./trakt-sync sync -a update_watchlist -description "short watchlist description" -sort_how asc
```
```console
$ ./trakt-sync sync -a update_watchlist -description "short watchlist description" -sort_by added
```
```console
$ ./trakt-sync sync -a update_watchlist -sort_by added -sort_how desc
```
##### Add items to Watchlist - via -items flag
```console
$ ./trakt-sync sync -t movies -a add_to_watchlist -items export_sync_watchlist_movies.json
```
```console
$ ./trakt-sync sync -t shows -a add_to_watchlist -items export_sync_watchlist_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a add_to_watchlist -items export_sync_watchlist_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a add_to_watchlist -items export_sync_watchlist_seasons.json
```
##### Add items to Watchlist - via stdin
```console
$ cat export_sync_watchlist_movies.json | ./trakt-sync sync -t movies -a add_to_watchlist
```
```console
$ cat export_sync_watchlist_shows.json | ./trakt-sync sync -t shows -a add_to_watchlist
```
```console
$ cat export_sync_watchlist_episodes.json | ./trakt-sync sync -t episodes -a add_to_watchlist
```
```console
$ cat export_sync_watchlist_seasons.json | ./trakt-sync sync -t seasons -a add_to_watchlist
```
```console
$ cat export_sync_watchlist_all.json | ./trakt-sync sync -t seasons -a add_to_watchlist
```
##### Remove from Watchlist - via -items flag
```console
$ ./trakt-sync sync -t movies -a remove_from_watchlist -items export_sync_watchlist_movies.json
```
```console
$ ./trakt-sync sync -t shows -a remove_from_watchlist -items export_sync_watchlist_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a remove_from_watchlist -items export_sync_watchlist_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a remove_from_watchlist -items export_sync_watchlist_seasons.json
```
##### Remove from Watchlist - via stdin
```console
$ cat export_sync_watchlist_movies.json | ./trakt-sync sync -t movies -a remove_from_watchlist
```
```console
$ cat export_sync_watchlist_shows.json | ./trakt-sync sync -t shows -a remove_from_watchlist
```
```console
$ cat export_sync_watchlist_episodes.json | ./trakt-sync sync -t episodes -a remove_from_watchlist
```
```console
$ cat export_sync_watchlist_seasons.json | ./trakt-sync sync -t seasons -a remove_from_watchlist
```
```console
$ cat export_sync_watchlist_all.json | ./trakt-sync sync -t seasons -a remove_from_watchlist
```
##### Reorder Watchlist - via -items flag
```console
$ ./trakt-sync sync -a reorder_watchlist -items export_sync_watchlist_movies.json
```
##### Reorder Watchlist - via stdin
```console
$ cat export_sync_watchlist_movies.json | ./trakt-sync sync -a reorder_watchlist
```
##### Update Watchlist item
```console
$ ./trakt-sync sync -a update_watchlist_item -list_item_id 97857 -notes "super 10/10"
```
##### Get Favorites
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_how asc
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_how des
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by rank
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by added
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by title
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by released
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by runtime
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by popularity
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by random
```
```console
$ ./trakt-sync sync -a get_favorites -t movies -sort_by percentage
```
```console

🔥VIP Only including imdb_rating, tmdb_rating, rt_tomatometer, rt_audience, metascore, votes,
imdb_votes, and tmdb_votes. If sent for a non VIP, the items will fall back to rank.
```

##### Update Favorites
```console
$ ./trakt-sync sync -a update_favorites -description "short favorites description" -sort_how asc
```
```console
$ ./trakt-sync sync -a update_favorites -description "short favorites description" -sort_by added
```
```console
$ ./trakt-sync sync -a update_favorites -sort_by added -sort_how desc
```
##### Add items to Favorites - via -items flag
```console
$ ./trakt-sync sync -t movies -a add_to_favorites -items export_sync_favorites_movies.json
```
```console
$ ./trakt-sync sync -t shows -a add_to_favorites -items export_sync_favorites_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a add_to_favorites -items export_sync_favorites_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a add_to_favorites -items export_sync_favorites_seasons.json
```
##### Add items to Favorites - via stdin
```console
$ cat export_sync_favorites_movies.json | ./trakt-sync sync -t movies -a add_to_favorites
```
```console
$ cat export_sync_favorites_shows.json | ./trakt-sync sync -t shows -a add_to_favorites
```
```console
$ cat export_sync_favorites_episodes.json | ./trakt-sync sync -t episodes -a add_to_favorites
```
```console
$ cat export_sync_favorites_seasons.json | ./trakt-sync sync -t seasons -a add_to_favorites
```
```console
$ cat export_sync_favorites_all.json | ./trakt-sync sync -t seasons -a add_to_favorites
```
##### Remove from Favorites - via -items flag
```console
$ ./trakt-sync sync -t movies -a remove_from_favorites -items export_sync_favorites_movies.json
```
```console
$ ./trakt-sync sync -t shows -a remove_from_favorites -items export_sync_favorites_shows.json
```
```console
$ ./trakt-sync sync -t episodes -a remove_from_favorites -items export_sync_favorites_episodes.json
```
```console
$ ./trakt-sync sync -t seasons -a remove_from_favorites -items export_sync_favorites_seasons.json
```
##### Remove from Favorites - via stdin
```console
$ cat export_sync_favorites_movies.json | ./trakt-sync sync -t movies -a remove_from_favorites
```
```console
$ cat export_sync_favorites_shows.json | ./trakt-sync sync -t shows -a remove_from_favorites
```
```console
$ cat export_sync_favorites_episodes.json | ./trakt-sync sync -t episodes -a remove_from_favorites
```
```console
$ cat export_sync_favorites_seasons.json | ./trakt-sync sync -t seasons -a remove_from_favorites
```
```console
$ cat export_sync_favorites_all.json | ./trakt-sync sync -t seasons -a remove_from_favorites
```
##### Reorder Favorites - via -items flag
```console
$ ./trakt-sync sync -a reorder_favorites -items export_sync_favorites_movies.json
```
##### Reorder Favorites - via stdin
```console
$ cat export_sync_favorites_movies.json | ./trakt-sync sync -a reorder_favorites
```
##### Update Favorite item
```console
$ ./trakt-sync sync -a update_favorite_item -list_item_id 97857 -notes "super 10/10"
```
