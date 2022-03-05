module github.com/indece-official/go-gousu/v2/gousujwt

go 1.17

replace github.com/indece-official/go-gousu/v2 => ../

require github.com/indece-official/go-gousu/v2 v2.0.1

require (
	github.com/golang-jwt/jwt/v4 v4.3.0
	github.com/namsral/flag v1.7.4-pre
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
)
