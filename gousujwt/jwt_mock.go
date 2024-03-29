package gousujwt

import (
	"net/http"
)

// MockVerifier provides a simple mock for JWTVerifier
type MockVerifier struct {
	VerifyFunc                            func(w http.ResponseWriter, r *http.Request, groups []string) (*CustomClaims, error)
	VerifyFuncCalled                      int
	VerifyWithCustomClaimsFunc            func(w http.ResponseWriter, r *http.Request, groups []string, claims ICustomClaims) (ICustomClaims, error)
	VerifyWithCustomClaimsFuncCalled      int
	VerifyTokenWithCustomClaimsFunc       func(r *http.Request, authToken string, groups []string, claims ICustomClaims) (ICustomClaims, error)
	VerifyTokenWithCustomClaimsFuncCalled int
}

var _ IVerifier = (*MockVerifier)(nil)

// Verify calls VerifyFunc and increases VerifyFuncCalled
func (m *MockVerifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*CustomClaims, error) {
	m.VerifyFuncCalled++

	return m.VerifyFunc(w, r, groups)
}

// VerifyWithCustomClaims calls VerifyWithCustomClaimsFunc and increases VerifyWithCustomClaimsFuncCalled
func (m *MockVerifier) VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims ICustomClaims) (ICustomClaims, error) {
	m.VerifyWithCustomClaimsFuncCalled++

	return m.VerifyWithCustomClaimsFunc(w, r, groups, claims)
}

// VerifyWithCustomClaims calls VerifyTokenWithCustomClaimsFunc and increases VerifyTokenWithCustomClaimsFuncCalled
func (m *MockVerifier) VerifyTokenWithCustomClaims(r *http.Request, authToken string, groups []string, claims ICustomClaims) (ICustomClaims, error) {
	m.VerifyTokenWithCustomClaimsFuncCalled++

	return m.VerifyTokenWithCustomClaimsFunc(r, authToken, groups, claims)
}

// NewMockVerifier creates a new initialized instance of MockVerifier
func NewMockVerifier() *MockVerifier {
	return &MockVerifier{
		VerifyFunc: func(w http.ResponseWriter, r *http.Request, groups []string) (*CustomClaims, error) {
			return &CustomClaims{}, nil
		},
		VerifyWithCustomClaimsFunc: func(w http.ResponseWriter, r *http.Request, groups []string, claims ICustomClaims) (ICustomClaims, error) {
			return &CustomClaims{}, nil
		},
		VerifyTokenWithCustomClaimsFunc: func(r *http.Request, authToken string, groups []string, claims ICustomClaims) (ICustomClaims, error) {
			return &CustomClaims{}, nil
		},
	}
}
