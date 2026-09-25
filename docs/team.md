#### Team:
Trakt team members. `-ex full` adds profile details (joined date, location, about, ...), `-ex images` adds the avatar; both can be combined with `-ex full,images`.
```console
$ ./trakt-sync team -a members -> export_team_members.json
```
```console
$ ./trakt-sync team -a members -ex full,images -> export_team_members.json
```
