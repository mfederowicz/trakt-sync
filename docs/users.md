#### Users:
##### Fetch settings for current user:
```console
$ ./trakt-sync users -a settings
```
##### Following requests:
```console
$ ./trakt-sync users -a following_requests
```
##### Follower requests:
```console
$ ./trakt-sync users -a follower_requests
```
##### Approve Follower requests (with request id):
```console
$ ./trakt-sync users -a follower_requests -follower_request 123
```
##### Deny Follower requests (with request id):
```console
$ ./trakt-sync users -a follower_requests -follower_request 123 -deny
```
##### Fetch saved filters for selected user:
```console
$ ./trakt-sync users -a saved_filters -u username
```
##### Export hidden items for a section (movie,show,season,user):
```console
$ ./trakt-sync users -a hidden_items -t show -section calendar
```
```console
$ ./trakt-sync users -a hidden_items -t show -section progress_watched
```
```console
$ ./trakt-sync users -a hidden_items -t show -section progress_watched_reset
```
```console
$ ./trakt-sync users -a hidden_items -t show -section progress_collected
```
```console
$ ./trakt-sync users -a hidden_items -t show -section recommendations
```
```console
$ ./trakt-sync users -a hidden_items -t show -section comments
```
```console
$ ./trakt-sync users -a hidden_items -t show -section dropped
```

##### Add hidden items - via -items flag
```console
$ ./trakt-sync users -a add_hidden_items -items export_users_all.json
```
```console
$ ./trakt-sync users -a add_hidden_items -t movie -items export_users_all.json
```
```console
$ ./trakt-sync users -a add_hidden_items -t show -items export_users_all.json
```
```console
$ ./trakt-sync users -a add_hidden_items -t season -items export_users_all.json
```
```console
$ ./trakt-sync users -a add_hidden_items -t user -items export_users_all.json
```
##### Add hidden items - via stdin
```console
$ cat export_users_all.json | ./trakt-sync users -a add_hidden_items
```
```console
$ cat export_users_all.json | ./trakt-sync users -a add_hidden_items -t movie
```
```console
$ cat export_users_all.json | ./trakt-sync users -a add_hidden_items -t show
```
```console
$ cat export_users_all.json | ./trakt-sync users -a add_hidden_items -t season
```
```console
$ cat export_users_all.json | ./trakt-sync users -a add_hidden_items -t user
```
##### Remove hidden items - via -items flag
```console
$ ./trakt-sync users -a remove_hidden_items -items export_users_all.json
```
```console
$ ./trakt-sync users -a remove_hidden_items -t movie -items export_users_all.json
```
```console
$ ./trakt-sync users -a remove_hidden_items -t show -items export_users_all.json
```
```console
$ ./trakt-sync users -a remove_hidden_items -t season -items export_users_all.json
```
```console
$ ./trakt-sync users -a remove_hidden_items -t user -items export_users_all.json
```
##### Remove hidden items - via stdin
```console
$ cat export_users_all.json | ./trakt-sync users -a remove_hidden_items
```
```console
$ cat export_users_all.json | ./trakt-sync users -a remove_hidden_items -t movie
```
```console
$ cat export_users_all.json | ./trakt-sync users -a remove_hidden_items -t show
```
```console
$ cat export_users_all.json | ./trakt-sync users -a remove_hidden_items -t season
```
```console
$ cat export_users_all.json | ./trakt-sync users -a remove_hidden_items -t user
```
##### Fetch user profile:
```console
$ ./trakt-sync users -a profile -u username
```
##### Fetch user likes:
```console
$ ./trakt-sync users -a likes -u username
```
```console
$ ./trakt-sync users -a likes -u username -ex comments
```
##### Fetch user collection:
```console
$ ./trakt-sync users -a collection -u username
```
```console
$ ./trakt-sync users -a collection -u username -ex comments
```
```console
$ ./trakt-sync users -a collection -u username -t movies
```
```console
$ ./trakt-sync users -a collection -u username -t shows
```
##### Fetch user comments:
```console
$ ./trakt-sync users -a comments
```
```console
$ ./trakt-sync users -a comments -comment_type all
```
```console
$ ./trakt-sync users -a comments -comment_type reviews
```
```console
$ ./trakt-sync users -a comments -comment_type shouts
```
```console
$ ./trakt-sync users -a comments -t all -include_replies true
```
```console
$ ./trakt-sync users -a comments -t movies -include_replies true
```
```console
$ ./trakt-sync users -a comments -t shows -include_replies true
```
```console
$ ./trakt-sync users -a comments -t seasons -include_replies true
```
```console
$ ./trakt-sync users -a comments -t episodes -include_replies true
```
```console
$ ./trakt-sync users -a comments -t lists -comment_type reviews  -include_replies true
```
```console
$ ./trakt-sync users -a comments -t lists -comment_type reviews  -include_replies false
```
```console
$ ./trakt-sync users -a comments -t lists -comment_type reviews  -include_replies only
```
##### Fetch user notes:
```console
$ ./trakt-sync users -a notes
```
```console
$ ./trakt-sync users -a notes -t all
```
```console
$ ./trakt-sync users -a notes -t movies
```
```console
$ ./trakt-sync users -a notes -t shows
```
```console
$ ./trakt-sync users -a notes -t seasons
```
```console
$ ./trakt-sync users -a notes -t episodes
```
```console
$ ./trakt-sync users -a notes -t people
```
```console
$ ./trakt-sync users -a notes -t history
```
```console
$ ./trakt-sync users -a notes -t collection
```
```console
$ ./trakt-sync users -a notes -t ratings
```
##### Export movies or shows or episodes from user lists:
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t episodes
```
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t shows
```
```console
$ ./trakt-sync users -a lists -u username -i 123456 -t movies
```
##### Fetch lists for selected user:
```console
$ ./trakt-sync users -a lists -u username
```
##### Create personal list for selected user - via -item flag:
```console
$ ./trakt-sync users -a add_list -item personal_list.json
```
##### Create personal list for selected user - via stdin:
```console
$ cat personal_list.json | ./trakt-sync users -a add_list
```
##### Reorder lists - via -items flag:
```console
$ ./trakt-sync users -a reorder_lists -items personal_lists.json
```
##### Reorder lists - via stdin
```console
$ cat personal_lists.json | ./trakt-sync users -a reorder_lists
```
##### Fetch lists  a user can collaborate on:
```console
$ ./trakt-sync users -a collaborations -u username
```
##### Fetch single personal list (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a list -u username -i 123456
```
##### Update single personal list (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a update_list -u username -i 123456 -description "short watchlist description" -sort_how asc
```
```console
$ ./trakt-sync users -a update_list -u username -i 123456 -description "short watchlist description" -sort_by added
```
```console
$ ./trakt-sync users -a update_list -u username -i 123456 -sort_by added -sort_how desc
```
##### Delete single personal list (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a delete_list -u username -i 123456
```
##### Fetch all users who liked a list (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a list_likes -u username -i 123456
```
##### Like a list (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a list_like -u username -i 123456
```
##### Remove like on a list (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a list_like -u username -i 123456 -delete
```
##### Get single personal list items (Trakt ID or Trakt slug):
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_how asc
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_how des
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by rank
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by added
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by title
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by released
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by runtime
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by popularity
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by random
```
```console
$ ./trakt-sync users -a list_items -u username -i 123456 -t movies -sort_by percentage
```

```console
🔥VIP Only including imdb_rating, tmdb_rating, rt_tomatometer, rt_audience, metascore, votes,
imdb_votes, and tmdb_votes. If sent for a non VIP, the items will fall back to rank.
```

##### Add items to personal list - via -items flag
```console
$ ./trakt-sync users -u username -i 123456 -t movies -a add_list_items -items export_users_list_movies.json
```
```console
$ ./trakt-sync users -u username -i 123456 -t shows -a add_list_items -items export_users_list_shows.json
```
```console
$ ./trakt-sync users -u username -i 123456 -t episodes -a add_list_items -items export_users_list_episodes.json
```
```console
$ ./trakt-sync users -u username -i 123456 -t seasons -a add_list_items -items export_users_list_seasons.json
```
##### Add items to personal list - via stdin
```console
$ cat export_users_list_movies.json | ./trakt-sync users -u username -i 123456 -t movies -a add_list_items
```
```console
$ cat export_users_list_shows.json | ./trakt-sync users -u username -i 123456 -t shows -a add_list_items
```
```console
$ cat export_users_list_episodes.json | ./trakt-sync users -u username -i 123456 -t episodes -a add_list_items
```
```console
$ cat export_users_list_seasons.json | ./trakt-sync users -u username -i 123456 -t seasons -a add_list_items
```
```console
$ cat export_users_list_all.json | ./trakt-sync users -u username -i 123456 -t seasons -a add_list_items
```
##### Remove items from personal list - via -items flag
```console
$ ./trakt-sync users -u username -i 123456 -t movies -a remove_list_items -items export_users_list_movies.json
```
```console
$ ./trakt-sync users -u username -i 123456 -t shows -a remove_list_items -items export_users_list_shows.json
```
```console
$ ./trakt-sync users -u username -i 123456 -t episodes -a remove_list_items -items export_users_list_episodes.json
```
```console
$ ./trakt-sync users -u username -i 123456 -t seasons -a remove_list_items -items export_users_list_seasons.json
```
##### Remove items from personal list - via stdin
```console
$ cat export_users_list_movies.json | ./trakt-sync users -u username -i 123456 -t movies -a remove_list_items
```
```console
$ cat export_users_list_shows.json | ./trakt-sync users -u username -i 123456 -t shows -a remove_list_items
```
```console
$ cat export_users_list_episodes.json | ./trakt-sync users -u username -i 123456 -t episodes -a remove_list_items
```
```console
$ cat export_users_list_seasons.json | ./trakt-sync users -u username -i 123456 -t seasons -a remove_list_items
```
```console
$ cat export_users_list_all.json | ./trakt-sync users -u username -i 123456 -t seasons -a remove_list_items
```
##### Reorder items from personal list - via -items flag
```console
$ ./trakt-sync users -a reorder_list_items -items export_users_list_movies.json
```
##### Reorder items from personal list - via stdin
```console
$ cat export_users_list_movies.json | ./trakt-sync users -a reorder_list_items
```
##### Update list item (notes):
```console
$ ./trakt-sync users -a update_list_item -u username -i 27316587 -list_item_id 996953153 -notes "best romcom"
```
##### Get all list comments
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s likes
```
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s likes_30
```
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s replies
```
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s replies_30
```
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s plays
```
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s rating
```
```console
$ ./trakt-sync users -a list_comments -i star-wars-in-machete-order -s added
```
##### List report
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r duplicate -message "Duplicate of another list"
```
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r remove -message "Should be removed from Trakt"
```
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r metadata -message "Metadata is wrong (name, description, etc)"
```
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r adult -message "Contains adult content"
```
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r language -message "Not in English"
```
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r spam -message "Spam or self-promotion"
```
```console
$ ./trakt-sync users -a list_report -u username -i star-wars-in-machete-order -r other -message "Anything else"
```
##### Follow user:
```console
$ ./trakt-sync users -a follow -u username
```
##### Unfollow user:
```console
$ ./trakt-sync users -a unfollow -u username
```
##### Get blocked users
```console
$ ./trakt-sync users -a blocked_users
```
##### Block user
```console
$ ./trakt-sync users -a block -u username
```
##### Unblock user
```console
$ ./trakt-sync users -a unblock -u username
```
##### Get followers
```console
$ ./trakt-sync users -a followers -u username
```
##### Get following
```console
$ ./trakt-sync users -a following -u username
```
##### Get friends
```console
$ ./trakt-sync users -a friends -u username
```
##### Get history last 7 days
```console
$ ./trakt-sync users -a history -u username -start_at 2025-06-01 -end_at 2025-06-07
```
##### Get history filter by type
```console
$ ./trakt-sync users -a history -u username -t movies
```
```console
$ ./trakt-sync users -a history -u username -t shows
```
```console
$ ./trakt-sync users -a history -u username -t seasons
```
```console
$ ./trakt-sync users -a history -u username -t episodes
```
##### Get history filter for specific item
```console
$ ./trakt-sync users -a history -u username -item_id 123456
```
##### Get user ratings - movies - all ratings
```console
$ ./trakt-sync users -a ratings -t movies
```
##### Get user ratings - shows - filter for specific rating from 1 to 10
```console
$ ./trakt-sync users -a ratings -t shows -rating 1,2,3
```
##### Get user ratings - seasons
```console
$ ./trakt-sync users -a ratings -t seasons
```
##### Get user ratings - episodes
```console
$ ./trakt-sync users -a ratings -t episodes
```
##### Get user ratings - all ratings
```console
$ ./trakt-sync users -a ratings -t all
```
##### Get user watchlist
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_how asc
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_how des
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by rank
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by added
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by title
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by released
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by runtime
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by popularity
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by random
```
```console
$ ./trakt-sync users -a watchlist -u username -t movies -sort_by percentage
```

```console
🔥VIP Only sort_by including imdb_rating, tmdb_rating, rt_tomatometer, rt_audience, metascore, votes,
imdb_votes, and tmdb_votes. If sent for a non VIP, the items will fall back to rank.
```
##### Get user watchlist comments
```console
$ ./trakt-sync users -a watchlist_comments -u username -s likes
```
```console
$ ./trakt-sync users -a watchlist_comments -u username -s likes_30
```
```console
$ ./trakt-sync users -a watchlist_comments -u username -s replies
```
```console
$ ./trakt-sync users -a watchlist_comments -u username -s replies_30
```
```console
$ ./trakt-sync users -a watchlist_comments -u username -s plays
```
```console
$ ./trakt-sync users -a watchlist_comments -u username -s rating
```
```console
$ ./trakt-sync users -a watchlist_comments -u username -s added
```
##### Get user favorites
```console
$ ./trakt-sync users -a favorites -t movies -sort_how asc
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_how des
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by rank
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by added
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by title
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by released
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by runtime
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by popularity
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by random
```
```console
$ ./trakt-sync users -a favorites -t movies -sort_by percentage
```
```console

🔥VIP Only sort_by including imdb_rating, tmdb_rating, rt_tomatometer, rt_audience, metascore, votes,
imdb_votes, and tmdb_votes. If sent for a non VIP, the items will fall back to rank.
```
##### Get user favorites comments
```console
$ ./trakt-sync users -a favorites_comments -u username -s likes
```
```console
$ ./trakt-sync users -a favorites_comments -u username -s likes_30
```
```console
$ ./trakt-sync users -a favorites_comments -u username -s replies
```
```console
$ ./trakt-sync users -a favorites_comments -u username -s replies_30
```
```console
$ ./trakt-sync users -a favorites_comments -u username -s plays
```
```console
$ ./trakt-sync users -a favorites_comments -u username -s rating
```
```console
$ ./trakt-sync users -a favorites_comments -u username -s added
```
##### Get user watching
```console
$ ./trakt-sync users -a watching -u username
```
##### Fetch watched movies for selected user:
```console
$ ./trakt-sync users -a watched -t movies -u sean
```
##### Fetch watched shows for selected user:
```console
$ ./trakt-sync users -a watched -t shows -u sean
```
##### Fetch watched shows for selected user without seasons:
```console
$ ./trakt-sync users -a watched -t shows -u sean --ex noseasons
```
##### Fetch stats for selected user:
```console
$ ./trakt-sync users -a stats -u sean
```
##### Report user:
```console
$ ./trakt-sync users -a report -u username -r duplicate -message "Duplicate of another list"
```
```console
$ ./trakt-sync users -a report -u username -r remove -message "Should be removed from Trakt"
```
```console
$ ./trakt-sync users -a report -u username -r metadata -message "Metadata is wrong (name, description, etc)"
```
```console
$ ./trakt-sync users -a report -u username -r adult -message "Contains adult content"
```
```console
$ ./trakt-sync users -a report -u username -r language -message "Not in English"
```
```console
$ ./trakt-sync users -a report -u username -r spam -message "Spam or self-promotion"
```
```console
$ ./trakt-sync users -a report -u username -r other -message "Anything else"
```
