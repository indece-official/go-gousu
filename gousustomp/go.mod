module github.com/indece-official/go-gousu/gousustomp/v2

go 1.22

toolchain go1.22.4

replace github.com/indece-official/go-gousu/v2 => ../

require (
	github.com/go-stomp/stomp/v3 v3.1.0
	github.com/indece-official/go-gousu/v2 v2.2.0
)

require (
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	github.com/namsral/flag v1.7.4-pre // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
)
