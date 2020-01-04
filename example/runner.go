package example

import "github.com/indece-official/go-gousu"

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
