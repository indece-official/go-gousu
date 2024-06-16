module github.com/indece-official/go-gousu/gousupostgres/v2

go 1.22

toolchain go1.22.4

replace github.com/indece-official/go-gousu/v2 => ../

require (
	github.com/indece-official/go-gousu/v2 v2.2.0
	github.com/lib/pq v1.10.9
	github.com/namsral/flag v1.7.4-pre
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
