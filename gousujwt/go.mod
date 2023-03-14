module github.com/indece-official/go-gousu/gousujwt/v2

go 1.17

replace github.com/indece-official/go-gousu/v2 => ../

require github.com/indece-official/go-gousu/v2 v2.0.3

require (
	github.com/MicahParks/keyfunc v1.9.0
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/namsral/flag v1.7.4-pre
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
)
