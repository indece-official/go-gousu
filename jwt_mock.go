package gousu

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// MockJWTVerifier provides a simple mock for JWTVerifier
type MockJWTVerifier struct {
	VerifyFunc                       func(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool)
	VerifyFuncCalled                 int
	VerifyWithCustomClaimsFunc       func(w http.ResponseWriter, r *http.Request, groups []string, claims jwt.Claims) (jwt.Claims, bool)
	VerifyWithCustomClaimsFuncCalled int
}

var _ IJWTVerifier = (*JWTVerifier)(nil)

// Verify calls VerifyFunc and increases VerifyFuncCalled
func (m *MockJWTVerifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool) {
	m.VerifyFuncCalled++

	return m.VerifyFunc(w, r, groups)
}

// VerifyWithCustomClaims calls VerifyWithCustomClaimsFunc and increases VerifyWithCustomClaimsFuncCalled
func (m *MockJWTVerifier) VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims jwt.Claims) (jwt.Claims, bool) {
	m.VerifyWithCustomClaimsFuncCalled++

	return m.VerifyWithCustomClaimsFunc(w, r, groups, claims)
}

// NewMockJWTVerifier creates a new initialized instance of MockJWTVerifier
func NewMockJWTVerifier() *MockJWTVerifier {
	return &MockJWTVerifier{
		VerifyFunc: func(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool) {
			return &JWTCustomClaims{}, true
		},
		VerifyWithCustomClaimsFunc: func(w http.ResponseWriter, r *http.Request, groups []string, claims jwt.Claims) (jwt.Claims, bool) {
			return &JWTCustomClaims{}, true
		},
	}
}
