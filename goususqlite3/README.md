# SQLite3-Integration for Go Universal Service Utilities

Full docu for go-gousu on https://github.com/indece-official/go-gousu

## Usage
### Config flags
| Flag | Env-Var | Type | Default | Description |
| --- | --- | --- | --- | --- |
| _sqlite3\_filename_ | _SQLITE3_FILENAME_ | string | "db.sqlite" | Filename of sqlite database |
| _sqlite3\_cache_ | _SQLITE3_CACHE_ | string | "private" | SQLite3 cache mode: shared | private |
| _sqlite3\_mode_ | _SQLITE3_MODE_ | string | "rwc" | SQLite3 access mode: ro | rw | rwc | memory |
| _sqlite3\_max\_retries_ | _SQLITE3\_MAX\_RETRIES_ | int | 10 | Max number of retries before connecting to db fails  |
| _sqlite3\_retry\_interval_ | _SQLITE3\_RETRY\_INTERVAL_ | int | 6 | Interval in seconds between retrying to connect to db |
| _sqlite3\_max\_idle\_conns_ | _SQLITE3\_MAX\_IDLE\_CONNS_ | int | 0 | Maximum number of idle connections (0 is unlimited) |
| _sqlite3\_max\_open\_conns_ | _SQLITE3\_MAX\_OPEN\_CONNS_ | int | 0 | Maximum number of open connections (0 is unlimited) |
