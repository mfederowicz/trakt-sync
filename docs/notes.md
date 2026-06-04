#### Notes:
##### Adding notes for media types:
```console
$ ./trakt-sync notes -a notes -t movie -i the-sopranos -notes "xyz"
```
```console
$ ./trakt-sync notes -a notes -t show -i breaking-bad -notes "greate show"
```
```console
$ ./trakt-sync notes -a notes -t season -i 250341 -notes "greate season"
```
```console
$ ./trakt-sync notes -a notes -t episode -i 250341 -notes "greate episode"
```
```console
$ ./trakt-sync notes -a notes -t person -i john-wayne -notes "greate person"
```
```console
$ ./trakt-sync notes -a notes -t history -i 1234567 -notes "history note"
```
##### Adding notes depends on activities:
```console
$ ./trakt-sync notes -a notes -t collection -item episode -i 73629 -notes "great episode"
```
```console
$ ./trakt-sync notes -a notes -t collection -item movie -i despicable-me-4-2024 -notes "great animation"
```
```console
$ ./trakt-sync notes -a notes -t rating -item movie -i despicable-me-4-2024 -notes "great animation"
```
```console
$ ./trakt-sync notes -a notes -t rating -item episode -i 73629 -notes "overall 10/10"
```
```console
$ ./trakt-sync notes -a notes -t rating -item movie -i the-gorge-2025 -notes "overall 7/10"
```
```console
$ ./trakt-sync notes -a notes -t rating -item season -i 3961 -notes "overall 9/10"
```
```console
$ ./trakt-sync notes -a notes -t rating -item show -i the-sopranos -notes "overall 9/10"
```
##### Manage notes get/modify/delete:
```console
$ ./trakt-sync notes -a note -i 97857
```
```console
$ ./trakt-sync notes -a note -i 97857 -notes "super 10/10" -privacy public -spoiler
```
```console
$ ./trakt-sync notes -a note -i 97857 -delete
```
##### Get items attachment to note:
```console
$ ./trakt-sync notes -a item -i 97854
```

