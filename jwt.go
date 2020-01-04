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
	Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool)
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
func (j *JWTVerifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, bool) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	// Check if valid Bearer-Header
	if len(authorizationHeader) != 2 || authorizationHeader[1] == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid authorization header")
		return nil, false
	}

	// validate the token
	token, err := jwt.ParseWithClaims(authorizationHeader[1], &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return j.publicKey, nil
	})

	// branch out into the possible error from signing
	switch err.(type) {
	case nil: // no error
		if !token.Valid { // but may still be invalid
			j.log.Debugf("Authorization failed: invalid token")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Authorization failed")
			return nil, false
		}

		claims := token.Claims.(*JWTCustomClaims)

		for i := range groups {
			if !ContainsString(claims.Groups, groups[i]) {
				j.log.Debugf("Authorization failed: missing group %s", groups[i])

				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Authorization failed")
				return nil, false
			}
		}

		return claims, true

	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			j.log.Debugf("Authorization failed: JWT expired")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Authorization failed")
			return nil, false

		default:
			j.log.Debugf("Authorization failed: invalid JWT: %s", err)

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Authorization failed")
			return nil, false
		}

	default: // something else went wrong
		j.log.Debugf("Authorization failed: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal server error")
		return nil, false
	}
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
