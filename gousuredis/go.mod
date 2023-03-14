module github.com/indece-official/go-gousu/gousuredis/v2

go 1.17

replace github.com/indece-official/go-gousu/v2 => ../

require (
	github.com/go-redsync/redsync/v4 v4.8.1
	github.com/gomodule/redigo v1.8.9
	github.com/indece-official/go-gousu/v2 v2.0.3
	github.com/mna/redisc v1.3.2
	github.com/namsral/flag v1.7.4-pre
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
