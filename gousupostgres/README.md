# PostgreSQL-Integration for Go Universal Service Utilities

Full docu for go-gousu on https://github.com/indece-official/go-gousu

## Usage
### Config flags
| Flag | Env-Var | Type | Default | Description |
| --- | --- | --- | --- | --- |
| _postgres\_host_ | _POSTGRES\_HOST_ | string | "127.0.0.1" | DB-Server hostname or ip |
| _postgres\_port_ | _POSTGRES\_PORT_ | int | 5432 | DB-Server port |
| _postgres\_user_ | _POSTGRES\_USER_ | string | "" | User for authentication against db server |
| _postgres\_password_ | _POSTGRES\_PASSWORD_ | string | "" | Password for authentication against db server |
| _postgres\_database_ | _POSTGRES\_DATABASE_ | string | "" | Database to use on db server |
| _postgres\_max\_retries_ | _POSTGRES\_MAX\_RETRIES_ | int | 10 | Max number of retries before connecting to db fails  |
| _postgres\_retry\_interval_ | _POSTGRES\_RETRY\_INTERVAL_ | int | 6 | Interval in seconds between retrying to connect to db |
| _postgres\_max\_idle\_conns_ | _POSTGRES\_MAX\_IDLE\_CONNS_ | int | 0 | Maximum number of idle connections (0 is unlimited) |
| _postgres\_max\_open\_conns_ | _POSTGRES\_MAX\_OPEN\_CONNS_ | int | 0 | Maximum number of open connections (0 is unlimited) |
