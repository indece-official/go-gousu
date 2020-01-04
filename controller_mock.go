package gousu

// MockController for simply mocking IController
type MockController struct {
	StartFunc        func() error
	StopFunc         func() error
	HealthFunc       func() error
	StartFuncCalled  int
	StopFuncCalled   int
	HealthFuncCalled int
}

// MockController must implement IController
var _ IController = (*MockController)(nil)

// Start calls StartFunc and increases StartFuncCalled
func (c *MockController) Start() error {
	c.StartFuncCalled++

	return c.StartFunc()
}

// Stop calls StopFunc and increases StopFuncCalled
func (c *MockController) Stop() error {
	c.StopFuncCalled++

	return c.StopFunc()
}

// Health calls HealthFunc and increases HealthFuncCalled
func (c *MockController) Health() error {
	c.HealthFuncCalled++

	return c.HealthFunc()
}

// NewMockController creates a new initialized instance of MockController
func NewMockController() *MockController {
	return &MockController{
		StartFunc: func() error {
			return nil
		},
		StopFunc: func() error {
			return nil
		},
		HealthFunc: func() error {
			return nil
		},
		StartFuncCalled:  0,
		StopFuncCalled:   0,
		HealthFuncCalled: 0,
	}
}
