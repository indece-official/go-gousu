package gousu

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/namsral/flag"
)

var (
	jwtPublicKeyFile = flag.String("jwt_publickey", "", "")
)

// IJWTVerifier defines the interface of JWTVerifier
type IJWTVerifier interface {
	Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, error)
	VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims jwt.Claims) (jwt.Claims, error)
}

// JWTVerifier provides a simple JWT verification utility
//
// Used flags:
//   * jwt_publickey Filename of JWT-Public-Key-File ()
type JWTVerifier struct {
	log       *Log
	publicKey *ecdsa.PublicKey
}

var _ IJWTVerifier = (*JWTVerifier)(nil)

// JWTCustomClaims definies the format of the JWT payload
type JWTCustomClaims struct {
	jwt.StandardClaims
	UserID int      `json:"user_id"`
	Groups []string `json:"groups"`
}

func (j *JWTVerifier) load() error {
	var err error

	j.log.Infof("Using certificate %s for JWT verification", *jwtPublicKeyFile)

	publicKeyPem, err := ioutil.ReadFile(*jwtPublicKeyFile)
	if err != nil {
		return err
	}

	j.publicKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyPem)
	if err != nil {
		return err
	}

	return nil
}

// Verify validates the JWT from the authorization header and checks if the
// required groups are fullfiled
//
// If the authorization fails it sends a 403 and returns
func (j *JWTVerifier) VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims jwt.Claims) (jwt.Claims, error) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	// Check if valid Bearer-Header
	if len(authorizationHeader) != 2 || authorizationHeader[1] == "" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	// validate the token
	token, err := jwt.ParseWithClaims(authorizationHeader[1], claims, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return j.publicKey, nil
	})

	// branch out into the possible error from signing
	switch err.(type) {
	case nil: // no error
		if !token.Valid { // but may still be invalid
			return nil, fmt.Errorf("Authorization failed: invalid token")
		}

		customClaims := token.Claims.(*JWTCustomClaims)

		for i := range groups {
			if !ContainsString(customClaims.Groups, groups[i]) {
				return nil, fmt.Errorf("Authorization failed: missing group %s", groups[i])
			}
		}

		return customClaims, nil

	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			return nil, fmt.Errorf("Authorization failed: JWT expired")

		default:
			return nil, fmt.Errorf("Authorization failed: invalid JWT: %s", err)
		}

	default: // something else went wrong
		return nil, fmt.Errorf("Authorization failed: %s", err)

	}
}

// VerifyWithCustomClaims validates the JWT from the authorization header and checks if the
// required groups are fullfiled
//
// If the authorization fails it sends a 403 and returns
func (j *JWTVerifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, error) {
	claims, err := j.VerifyWithCustomClaims(w, r, groups, &JWTCustomClaims{})
	if err != nil {
		return nil, err
	}

	return claims.(*JWTCustomClaims), nil
}

// NewJWTVerifier creates a new initilized instance of JWTVerifier
// and loads the public key from the file specified by the flag 'jwt_publickey'
func NewJWTVerifier() (*JWTVerifier, error) {
	j := &JWTVerifier{
		log: GetLogger("utils.jwt"),
	}

	err := j.load()
	if err != nil {
		j.log.Fatalf("Error initializing jwt utiliy: %s", err)
		return nil, err
	}

	return j, nil
}
