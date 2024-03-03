# schizgophrenia-got
Attempt to migrate from terrifying javascript to distinguished golang. (God forgive me for my foolish decisions of the past.)

## Environment variables

| Variable                   | Required? | Possible values                                 | Description                                                                                     |
| -------------------------- | --------- | ----------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `CLIENT_TOKEN`             | yes       | `string`                                        | Discord bot token                                                                               |
| `CLIENT_ID`                | yes       | `string`                                        | Discord bot user ID                                                                             |
| `GUILD_ID`                 | no        | `string`                                        | Required for registering commands on single server, leave empty for global command registration |
| `SQLITE_DATABASE_FILEPATH` | no        | `string`                                        | SQLite database location; will be set to "./userdata.sqlite3.db" if not specified               |
| `REGISTER_COMMANDS`        | no        | `1`/`0`, `true`/`false`, `yes`/`no`, `on`/`off` | Wether to register commands; defaults to disabled                                               |
| `REMOVE_COMMANDS`          | no        | `1`/`0`, `true`/`false`, `yes`/`no`, `on`/`off` | Wether to remove registered commands; defaults to disabled                                      |
| `GREETING_TIME_INTERVAL`   | no        | `int`                                           | Time interval between greetings in utix seconds; defaults to 604800 (one week)                  |