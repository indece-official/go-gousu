package gousu

import (
	"net/http"
)

// MockJWTVerifier provides a simple mock for JWTVerifier
type MockJWTVerifier struct {
	VerifyFunc                       func(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, error)
	VerifyFuncCalled                 int
	VerifyWithCustomClaimsFunc       func(w http.ResponseWriter, r *http.Request, groups []string, claims IJWTCustomClaims) (IJWTCustomClaims, error)
	VerifyWithCustomClaimsFuncCalled int
}

var _ IJWTVerifier = (*MockJWTVerifier)(nil)

// Verify calls VerifyFunc and increases VerifyFuncCalled
func (m *MockJWTVerifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, error) {
	m.VerifyFuncCalled++

	return m.VerifyFunc(w, r, groups)
}

// VerifyWithCustomClaims calls VerifyWithCustomClaimsFunc and increases VerifyWithCustomClaimsFuncCalled
func (m *MockJWTVerifier) VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims IJWTCustomClaims) (IJWTCustomClaims, error) {
	m.VerifyWithCustomClaimsFuncCalled++

	return m.VerifyWithCustomClaimsFunc(w, r, groups, claims)
}

// NewMockJWTVerifier creates a new initialized instance of MockJWTVerifier
func NewMockJWTVerifier() *MockJWTVerifier {
	return &MockJWTVerifier{
		VerifyFunc: func(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, error) {
			return &JWTCustomClaims{}, nil
		},
		VerifyWithCustomClaimsFunc: func(w http.ResponseWriter, r *http.Request, groups []string, claims IJWTCustomClaims) (IJWTCustomClaims, error) {
			return &JWTCustomClaims{}, nil
		},
	}
}
