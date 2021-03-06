package gousu

import (
	"os"
	"testing"
)

type testService struct {
}

var _ (IService) = (*testService)(nil)

func (s *testService) Name() string  { return "test" }
func (s *testService) Start() error  { return nil }
func (s *testService) Health() error { return nil }
func (s *testService) Stop() error   { return nil }

func newTestService(ctx IContext) IService {
	return &testService{}
}

type testController struct {
	testService *testService
}

var _ (IController) = (*testController)(nil)

func (c *testController) Name() string  { return "test" }
func (c *testController) Start() error  { return nil }
func (c *testController) Health() error { return nil }
func (c *testController) Stop() error   { return nil }

func newTestController(ctx IContext) IController {
	return &testController{
		testService: ctx.GetService("test").(*testService),
	}
}

type testUIController struct {
	testService *testService
}

var _ (IUIController) = (*testUIController)(nil)

func (c *testUIController) Name() string  { return "test" }
func (c *testUIController) Start() error  { return nil }
func (c *testUIController) Health() error { return nil }
func (c *testUIController) Run(sigTerm chan os.Signal) error {
	<-sigTerm

	return nil
}
func (c *testUIController) Stop() error { return nil }

func newTestUIController(ctx IContext) IUIController {
	return &testUIController{
		testService: ctx.GetService("test").(*testService),
	}
}

func TestRunner(t *testing.T) {
	runner := NewRunner("example", "1.0.0")
	runner.CreateService(newTestService)
	runner.CreateController(newTestController)

	go runner.Run()

	runner.AwaitReady()
	runner.Kill()
}

func TestRunnerWithUIController(t *testing.T) {
	runner := NewRunner("example", "1.0.0")
	runner.CreateService(newTestService)
	runner.CreateController(newTestController)
	runner.CreateUIController(newTestUIController)

	go runner.Run()

	runner.AwaitReady()
	runner.Kill()
}
