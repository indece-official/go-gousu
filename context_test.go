package gousu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()

	assert.NotNil(t, ctx)
}

func TestRegisterServiceNoName(t *testing.T) {
	logFatalfOrg := logFatalf
	logFatalfCalled := 0
	logFatalf = func(format string, v ...interface{}) { logFatalfCalled++ }

	service := NewMockService()
	service.NameFunc = func() string { return "" }

	ctx := NewContext()
	ctx.RegisterService(service)

	retServices := ctx.GetServices()
	assert.Len(t, retServices, 0)
	assert.Equal(t, 1, logFatalfCalled)

	logFatalf = logFatalfOrg
}

func TestRegisterServiceSameName(t *testing.T) {
	logFatalfOrg := logFatalf
	logFatalfCalled := 0
	logFatalf = func(format string, v ...interface{}) { logFatalfCalled++ }

	service0 := NewMockService()
	service0.NameFunc = func() string { return "mock" }
	service1 := NewMockService()
	service1.NameFunc = func() string { return "mock" }

	ctx := NewContext()

	ctx.RegisterService(service0)
	assert.Len(t, ctx.GetServices(), 1)
	assert.Equal(t, 0, logFatalfCalled)

	ctx.RegisterService(service1)
	assert.Len(t, ctx.GetServices(), 1)
	assert.Equal(t, 1, logFatalfCalled)

	logFatalf = logFatalfOrg
}

func TestRegisterAndGetService(t *testing.T) {
	service0 := NewMockService()
	service0.NameFunc = func() string { return "mock0" }
	service1 := NewMockService()
	service1.NameFunc = func() string { return "mock1" }

	ctx := NewContext()
	retServices := ctx.GetServices()
	assert.Len(t, retServices, 0)

	ctx.RegisterService(service0)
	ctx.RegisterService(service1)
	retService0 := ctx.GetService("mock0")
	retService1 := ctx.GetService("mock1")

	assert.Equal(t, retService0, service0)
	assert.Equal(t, retService1, service1)

	retServices = ctx.GetServices()
	assert.Len(t, retServices, 2)
	assert.Contains(t, retServices, service0)
	assert.Contains(t, retServices, service1)
}

func TestRegisterControllerNoName(t *testing.T) {
	logFatalfOrg := logFatalf
	logFatalfCalled := 0
	logFatalf = func(format string, v ...interface{}) { logFatalfCalled++ }

	controller := NewMockController()
	controller.NameFunc = func() string { return "" }

	ctx := NewContext()
	ctx.RegisterController(controller)

	retControllers := ctx.GetControllers()
	assert.Len(t, retControllers, 0)
	assert.Equal(t, 1, logFatalfCalled)

	logFatalf = logFatalfOrg
}

func TestRegisterControllerSameName(t *testing.T) {
	logFatalfOrg := logFatalf
	logFatalfCalled := 0
	logFatalf = func(format string, v ...interface{}) { logFatalfCalled++ }

	controller0 := NewMockController()
	controller0.NameFunc = func() string { return "mock" }
	controller1 := NewMockController()
	controller1.NameFunc = func() string { return "mock" }

	ctx := NewContext()

	ctx.RegisterController(controller0)
	assert.Len(t, ctx.GetControllers(), 1)
	assert.Equal(t, 0, logFatalfCalled)

	ctx.RegisterController(controller1)
	assert.Len(t, ctx.GetControllers(), 1)
	assert.Equal(t, 1, logFatalfCalled)

	logFatalf = logFatalfOrg
}

func TestRegisterAndGetController(t *testing.T) {
	controller0 := NewMockController()
	controller0.NameFunc = func() string { return "mock0" }
	controller1 := NewMockController()
	controller1.NameFunc = func() string { return "mock1" }

	ctx := NewContext()
	retControllers := ctx.GetControllers()
	assert.Len(t, retControllers, 0)

	ctx.RegisterController(controller0)
	ctx.RegisterController(controller1)
	retController0 := ctx.GetController("mock0")
	retController1 := ctx.GetController("mock1")

	assert.Equal(t, retController0, controller0)
	assert.Equal(t, retController1, controller1)

	retControllers = ctx.GetControllers()
	assert.Len(t, retControllers, 2)
	assert.Contains(t, retControllers, controller0)
	assert.Contains(t, retControllers, controller1)
}
