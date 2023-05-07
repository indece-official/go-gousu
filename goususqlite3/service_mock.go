package goususqlite3

import (
	"database/sql"

	"github.com/indece-official/go-gousu/v2/gousu"
)

// MockService for simply mocking IService
type MockService struct {
	gousu.MockService

	GetDBFunc           func() *sql.DB
	GetDBSafeFunc       func() (*sql.DB, error)
	GetDBFuncCalled     int
	GetDBSafeFuncCalled int
}

// MockService implements IService
var _ (IService) = (*MockService)(nil)

// GetDB calls GetDBFunc and increases GetDBFuncCalled
func (s *MockService) GetDB() *sql.DB {
	s.GetDBFuncCalled++

	return s.GetDBFunc()
}

// GetDBSafe calls GetDBSafeFunc and increases GetDBSafeFuncCalled
func (s *MockService) GetDBSafe() (*sql.DB, error) {
	s.GetDBSafeFuncCalled++

	return s.GetDBSafeFunc()
}

// NewMockService creates a new initialized instance of MockService
func NewMockService() *MockService {
	return &MockService{
		MockService: gousu.MockService{
			NameFunc: func() string {
				return ServiceName
			},
		},

		GetDBFunc: func() *sql.DB {
			return &sql.DB{}
		},
		GetDBSafeFunc: func() (*sql.DB, error) {
			return &sql.DB{}, nil
		},
	}
}
