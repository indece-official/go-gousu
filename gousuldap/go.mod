module github.com/indece-official/go-gousu/gousuldap/v2

go 1.17

replace github.com/indece-official/go-gousu/v2 => ../

require (
	github.com/go-ldap/ldap/v3 v3.4.4
	github.com/indece-official/go-gousu/v2 v2.0.3
	github.com/namsral/flag v1.7.4-pre
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
)
