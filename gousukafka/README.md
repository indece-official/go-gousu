# Kafka-Integration for Go Universal Service Utilities

Full docu for go-gousu on https://github.com/indece-official/go-gousu

## Requirements

## Services & Controllers
### Services
| Name | Description | Interface | Factory |
| --- | --- | --- | --- |
| _kafka_ | Kafka client service for consuming & producing kafka messages | gousukafka.IService | gousukafka.NewService | 

### Controllers
| Name | Description | Interface | Struct |
| --- | --- | --- | --- |
| _kafka_ | Base  controller for consuming one kafka topic | gousukafka.IController | gousukafka.NewControllerBase | 

## Usage
**TODO**