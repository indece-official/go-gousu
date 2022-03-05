package gousuldap

import (
	"github.com/indece-official/go-gousu/v2/gousu"
)

// MockService for simply mocking ldap IService
type MockService struct {
	gousu.MockService

	SimpleLoginFunc       func(username string, password string, attributes []string) (*map[string][]string, error)
	ReconnectFunc         func() error
	SimpleLoginFuncCalled int
	ReconnectFuncCalled   int
}

// MockService implements IService
var _ (IService) = (*MockService)(nil)

// SimpleLogin calls SimpleLoginFunc and increases SimpleLoginFuncCalled
func (s *MockService) SimpleLogin(username string, password string, attributes []string) (*map[string][]string, error) {
	s.SimpleLoginFuncCalled++

	return s.SimpleLoginFunc(username, password, attributes)
}

// Reconnect calls ReconnectFunc and increases ReconnectFuncCalled
func (s *MockService) Reconnect() error {
	s.ReconnectFuncCalled++

	return s.ReconnectFunc()
}

// NewMockService creates a new initialized instance of MockService
func NewMockService() *MockService {
	return &MockService{
		MockService: gousu.MockService{},

		SimpleLoginFunc: func(username string, password string, attributes []string) (*map[string][]string, error) {
			return &map[string][]string{}, nil
		},
		ReconnectFunc: func() error {
			return nil
		},
		SimpleLoginFuncCalled: 0,
		ReconnectFuncCalled:   0,
	}
}
