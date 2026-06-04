#### Scrobble:
##### Scrobble start/pause/stop movie:
```console
$ ./trakt-sync scrobble -a start -t movie -progress 3.45 -i guardians-of-the-galaxy-2014
```
```console
$ ./trakt-sync scrobble -a pause -t movie -progress 3.45 -i guardians-of-the-galaxy-2014
```
```console
$ ./trakt-sync scrobble -a stop -t movie -progress 3.45 -i guardians-of-the-galaxy-2014
```
##### Scrobble start/pause/stop episode:
```console
$ ./trakt-sync scrobble -a start -t episode -i 73629 -progress 10.25
```
```console
$ ./trakt-sync scrobble -a pause -t episode -i 73629 -progress 10.25
```
```console
$ ./trakt-sync scrobble -a stop -t episode -i 73629 -progress 50.25
```
##### Scrobble start/pause/stop show by episode code (format season x episode):
```console
$ ./trakt-sync scrobble -a start -t show_episode -i 136121 -episode_code 1x5 -progress 3.45
```
```console
$ ./trakt-sync scrobble -a pause -t show_episode -i 136121 -episode_code 1x5 -progress 3.45
```
```console
$ ./trakt-sync scrobble -a stop -t show_episode -i 136121 -episode_code 1x5 -progress 3.45
```
##### Scrobble start/pause/stop show by episode abs number (useful for Anime and Donghua):
```console
$ ./trakt-sync scrobble -a start -t show_episode -i 37696 -episode_abs 164 -progress 50
```
```console
$ ./trakt-sync scrobble -a pause -t show_episode -i 37696 -episode_abs 164 -progress 50
```
```console
$ ./trakt-sync scrobble -a stop -t show_episode -i 37696 -episode_abs 164 -progress 60
```
