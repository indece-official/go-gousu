package goususmtp

import "github.com/indece-official/go-gousu/v2/gousu"

// MockService for simply mocking IService
type MockService struct {
	gousu.MockService

	SendEmailFunc       func(m *Email) error
	PingFunc            func() error
	SendEmailFuncCalled int
	PingFuncCalled      int
}

// MockService implements IService
var _ (IService) = (*MockService)(nil)

// SendEmail calls SendEmailFunc and increases SendEmailFuncCalled
func (s *MockService) SendEmail(m *Email) error {
	s.SendEmailFuncCalled++

	return s.SendEmailFunc(m)
}

// Ping calls PingFunc and increases PingFuncCalled
func (s *MockService) Ping() error {
	s.PingFuncCalled++

	return s.PingFunc()
}

// NewMockService creates a new initialized instance of MockService
func NewMockService() *MockService {
	return &MockService{
		MockService: gousu.MockService{
			NameFunc: func() string {
				return ServiceName
			},
		},
		SendEmailFunc: func(m *Email) error {
			return nil
		},
		PingFunc: func() error {
			return nil
		},
	}
}
