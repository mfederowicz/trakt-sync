#### Comments:
```console
$ ./trakt-sync comments -a comment -comment_id 779883 -comment "minions,minions,minions movie ever ok"
```
```console
$ ./trakt-sync comments -a comment -comment_id 779883 -delete
```
```console
$ ./trakt-sync comments -a comments -t episode -trakt_id 172245 -comment "super episode, interesting plot ok"
```
```console
$ ./trakt-sync comments -a replies -comment_id 779896 -reply "reply msg min 5 words" -spoiler
```
```console
$ ./trakt-sync comments -a replies -comment_id 71340
```
```console
$ ./trakt-sync comments -a item -comment_id 664237 -ex full
```
```console
$ ./trakt-sync comments -a likes -comment_id 773108 -remove
```
```console
$ ./trakt-sync comments -a like -comment_id 773108
```
```console
$ ./trakt-sync comments -a like -comment_id 773108 -remove
```
```console
$ ./trakt-sync comments -a trending -comment_type reviews
```
```console
$ ./trakt-sync comments -a recent -include_replies false
```
```console
$ ./trakt-sync comments -a recent -include_replies true
```
```console
$ ./trakt-sync comments -a updates -include_replies false
```
```console
$ ./trakt-sync comments -a reactions -comment_id 773108 -> export_comments_reactions_773108.json
```
```console
$ ./trakt-sync comments -a reactions_summary -comment_id 773108 -> export_comments_reactions_summary_773108.json
```
```console
$ ./trakt-sync comments -a reaction -comment_id 773108 -reaction love
```
```console
$ ./trakt-sync comments -a reaction -comment_id 773108 -reaction love -remove
```
```console
$ ./trakt-sync comments -a report -comment_id 773108 -r spoilers -message "reveals the ending"
```
