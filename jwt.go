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
	jwtPublicKeyFile            = flag.String("jwt_publickey", "", "JWT-Public-Key for Verifier")
	jwtVerifyAudience           = flag.String("jwt_verify_audience", "", "JWT-Audience for Verifier")
	jwtVerifyAlgorithm          = flag.String("jwt_verify_algorithm", "", "JWT-Algorithm for Verifier")
	jwtVerifyNoSuccessSiemEvent = flag.Bool("jwt_verify_no_success_siem_event", true, "Don't log success siem event from JWT-Verifier")
)

// IJWTVerifier defines the interface of JWTVerifier
type IJWTVerifier interface {
	Verify(w http.ResponseWriter, r *http.Request, groups []string) (*JWTCustomClaims, error)
	VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims IJWTCustomClaims) (IJWTCustomClaims, error)
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

// IJWTCustomClaims definies the generic interface for custom jwt claims
type IJWTCustomClaims interface {
	jwt.Claims
	GetUserID() int64
	GetGroups() []string
	GetAudience() string
}

// JWTCustomClaims definies the format of the JWT payload
type JWTCustomClaims struct {
	jwt.StandardClaims
	UserID int64    `json:"user_id"`
	Groups []string `json:"groups"`
}

// GetUserID returns the UserID
func (j *JWTCustomClaims) GetUserID() int64 {
	return j.UserID
}

// GetGroups returns all Groups
func (j *JWTCustomClaims) GetGroups() []string {
	return j.Groups
}

// GetAudience returns the audience
func (j *JWTCustomClaims) GetAudience() string {
	return j.Audience
}

var _ IJWTCustomClaims = (*JWTCustomClaims)(nil)

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

func (j *JWTVerifier) generateSiemEvent(r *http.Request, eventType SiemEventType, claims IJWTCustomClaims) *SiemEvent {
	evt := &SiemEvent{}

	evt.Type = eventType
	evt.SourceIP.Scan(r.RemoteAddr)
	evt.SourceRealIP.Scan(r.Header.Get("X-Real-IP"))

	if claims != nil {
		evt.UserIdentifier.Scan(fmt.Sprintf("userid:%d", claims.GetUserID()))
	}

	return evt
}

// VerifyWithCustomClaims validates the JWT from the authorization header and checks if the
// required groups are fullfiled
//
// If the authorization fails it sends a 403 and returns
func (j *JWTVerifier) VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims IJWTCustomClaims) (IJWTCustomClaims, error) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	// Check if valid Bearer-Header
	if len(authorizationHeader) != 2 || authorizationHeader[1] == "" {
		j.log.SiemEvent(
			j.generateSiemEvent(
				r,
				SiemEventTypeAuthenticationFailed,
				nil,
			),
			"Invalid or empty authorization header: %s",
			r.Header.Get("Authorization"),
		)

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
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s",
				authorizationHeader[1],
			)

			return nil, fmt.Errorf("Authorization failed: invalid token")
		}

		err = token.Claims.Valid()
		if err != nil {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt claims in token %s: %s",
				authorizationHeader[1],
				err,
			)

			return nil, fmt.Errorf("Authorization failed: invalid claims in token: %s", err)
		}

		if token.Method.Alg() != *jwtVerifyAlgorithm {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s - missmatching algorithm: got %s, expected %s",
				authorizationHeader[1],
				token.Method.Alg(),
				*jwtVerifyAlgorithm,
			)

			return nil, fmt.Errorf("Missmatching algorithm: got %s, expected %s", token.Method.Alg(), *jwtVerifyAlgorithm)
		}

		customClaims, ok := token.Claims.(IJWTCustomClaims)
		if !ok {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s - casting jwt custom claims failed",
				*jwtVerifyAlgorithm,
			)

			return nil, fmt.Errorf("Casting jwt custom claims failed")
		}

		if customClaims.GetAudience() != *jwtVerifyAudience {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s - missmatching audience: got %s, expected %s",
				authorizationHeader[1],
				customClaims.GetAudience(),
				*jwtVerifyAudience,
			)

			return nil, fmt.Errorf("Missmatching audience: got %s, expected %s", customClaims.GetAudience(), *jwtVerifyAudience)
		}

		for i := range groups {
			if !ContainsString(customClaims.GetGroups(), groups[i]) {
				j.log.SiemEvent(
					j.generateSiemEvent(
						r,
						SiemEventTypeAuthenticationFailedAttact,
						nil,
					),
					"Authorization failed: missing group in token %s: %s",
					authorizationHeader[1],
					groups[i],
				)

				return nil, fmt.Errorf("Authorization failed: missing group %s", groups[i])
			}
		}

		if !*jwtVerifyNoSuccessSiemEvent {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationSuccess,
					customClaims,
				),
				"JWT-Verification succeded for user %d",
				customClaims.GetUserID(),
			)
		}

		return customClaims, nil

	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailed,
					nil,
				),
				"Authorization failed: jwt expired: %s",
				authorizationHeader[1],
			)

			return nil, fmt.Errorf("Authorization failed: JWT expired")

		default:
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					SiemEventTypeAuthenticationFailedAttact,
					nil,
				),
				"Authorization failed: invalid jwt %s: %s",
				authorizationHeader[1],
				err,
			)

			return nil, fmt.Errorf("Authorization failed: invalid JWT: %s", err)
		}

	default: // something else went wrong
		j.log.SiemEvent(
			j.generateSiemEvent(
				r,
				SiemEventTypeAuthenticationFailedAttact,
				nil,
			),
			"Authorization failed for token %s: %s",
			authorizationHeader[1],
			err,
		)

		return nil, fmt.Errorf("Authorization failed: %s", err)
	}
}

// Verify validates the JWT from the authorization header and checks if the
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
		j.log.Errorf("Error initializing jwt utiliy: %s", err)
		return nil, err
	}

	return j, nil
}
