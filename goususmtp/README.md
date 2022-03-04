# SMTP-Integration for Go Universal Service Utilities

Full docu for go-gousu on https://github.com/indece-official/go-gousu

[![GoDoc](https://godoc.org/github.com/indece-official/go-gousu-smtp?status.svg)](https://godoc.org/github.com/indece-official/go-gousu-smtp)

## Usage
### Config flags
| Flag | Env-Var | Type | Default | Description |
| --- | --- | --- | --- | --- |
| _smtp\_host_ | _SMTP\_HOST_ | string | "127.0.0.1" | SMTP-Server hostname or ip |
| _smtp\_port_ | _SMTP\_PORT_ | int | 587 | SMTP-Server port |
| _smtp\_user_ | _SMTP\_USER_ | string | "" | User for authentication against SMTP-Server |
| _smtp\_password_ | _SMTP\_PASSWORD_ | string | "" | Password for authentication against SMTP-Server |
| _smtp\_from_ | _SMTP\_FROM_ | string | "" | Default `From`-Header in emails |

**TODO**