package gousu

// MockService for simply mocking IService
type MockService struct {
	NameFunc         func() string
	StartFunc        func() error
	StopFunc         func() error
	HealthFunc       func() error
	NameFuncCalled   int
	StartFuncCalled  int
	StopFuncCalled   int
	HealthFuncCalled int
}

// MockService must implement IService
var _ IService = (*MockService)(nil)

// Name calls NameFunc and increases NameFuncCalled
func (c *MockService) Name() string {
	c.NameFuncCalled++

	return c.NameFunc()
}

// Start calls StartFunc and increases StartFuncCalled
func (c *MockService) Start() error {
	c.StartFuncCalled++

	return c.StartFunc()
}

// Stop calls StopFunc and increases StopFuncCalled
func (c *MockService) Stop() error {
	c.StopFuncCalled++

	return c.StopFunc()
}

// Health calls HealthFunc and increases HealthFuncCalled
func (c *MockService) Health() error {
	c.HealthFuncCalled++

	return c.HealthFunc()
}

// NewMockService creates a new initialized instance of MockService
func NewMockService() *MockService {
	return &MockService{
		NameFunc: func() string {
			return "mockservice"
		},
		StartFunc: func() error {
			return nil
		},
		StopFunc: func() error {
			return nil
		},
		HealthFunc: func() error {
			return nil
		},
		NameFuncCalled:   0,
		StartFuncCalled:  0,
		StopFuncCalled:   0,
		HealthFuncCalled: 0,
	}
}
