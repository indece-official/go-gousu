module github.com/indece-official/go-gousu/gousukafka/v2

go 1.22

toolchain go1.22.4

replace github.com/indece-official/go-gousu/v2 => ../

require (
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/indece-official/go-gousu/v2 v2.2.0
	github.com/namsral/flag v1.7.4-pre
)

require (
	github.com/chakrit/go-bunyan v0.0.0-20140303180041-5a9b5e7b1765 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
)
