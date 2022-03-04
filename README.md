# Go Universal Service Utilities
Golang framework for writing lightweight microservices

[![GoDoc](https://godoc.org/github.com/indece-official/go-gousu?status.svg)](https://godoc.org/github.com/indece-official/go-gousu)

## Modules
| Name | Description |
| --- | --- |
| [gogousu] (https://github.com/indece-official/go-gousu) | Core & utilities Module (dependency injection, logging, siem events, ...) |
| [gogousujwt](https://github.com/indece-official/go-gousu/tree/main/gousujwt) | JWT-Verification |
| [gogousukafka](https://github.com/indece-official/go-gousu/tree/main/gousukafka) | Kafka service & controller |
| [gogousuredis](https://github.com/indece-official/go-gousu/tree/main/gousuredis) | Redis service |
| [gogousupostgres](https://github.com/indece-official/go-gousu/tree/main/gousupostgres) | Postgres service |
| [gogoususmtp](https://github.com/indece-official/go-gousu/tree/main/goususmtp) | SMTP service |
| [gogousuchi](https://github.com/indece-official/go-gousu/tree/main/gousuchi) | Chi controller |

## Usage
### Example
```
package example

import "github.com/indece-official/go-gousu/gousu"

type TestService struct {
}

var _ (gousu.IService) = (*TestService)(nil)

func (s *TestService) Name() string  { return "test" }
func (s *TestService) Start() error  { return nil }
func (s *TestService) Health() error { return nil }
func (s *TestService) Stop() error   { return nil }

func NewTestService(ctx gousu.IContext) gousu.IService {
	return &TestService{}
}

type TestController struct {
	testService *TestService
}

var _ (gousu.IController) = (*TestController)(nil)

func (c *TestController) Name() string  { return "test" }
func (c *TestController) Start() error  { return nil }
func (c *TestController) Health() error { return nil }
func (c *TestController) Stop() error   { return nil }

func NewTestController(ctx gousu.IContext) gousu.IController {
	return &TestController{
		testService: ctx.GetService("test").(*TestService),
	}
}

func ExampleMain() {
	runner := gousu.NewRunner("example", "1.0.0")
	runner.CreateService(NewTestService)
	runner.CreateController(NewTestController)

	runner.Run()
}
```
