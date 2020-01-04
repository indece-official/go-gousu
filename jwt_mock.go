package gousu

import "net/http"

// MockJWTVerifier provides a simple mock for JWTVerifier
type MockJWTVerifier struct {
	VerifyFunc       func(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool)
	VerifyFuncCalled int
}

var _ IJWTVerifier = (*JWTVerifier)(nil)

// Verify calls VerifyFunc and increases VerifyFuncCalled
func (m *MockJWTVerifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool) {
	m.VerifyFuncCalled++

	return m.VerifyFunc(w, r, groups)
}

// NewMockJWTVerifier creates a new initialized instance of MockJWTVerifier
func NewMockJWTVerifier() *MockJWTVerifier {
	return &MockJWTVerifier{
		VerifyFunc: func(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool) {
			return &JWTCustomClaims{}, true
		},
		VerifyFuncCalled: 0,
	}
}
