package gousujwt

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/go-gousu/v2/gousu/siem"
	"github.com/namsral/flag"
)

var (
	jwksURL                     = flag.String("jwks_url", "", "JWKS-URL for Verifier")
	jwksRefreshInterval         = flag.Int("jwks_refresh_interval", 3600, "Interval for JWKS-Refresh [s]")
	jwksRefreshRateLimit        = flag.Int("jwks_refresh_rate_limit", 300, "Rate-Limit for JWKS-Refresh [s]")
	jwksRefreshTimeout          = flag.Int("jwks_refresh_timeout", 10, "Timeout for JWKS-Refresh [s]")
	jwksRefreshUnknownKID       = flag.Bool("jwks_refresh_unknown_kid", true, "Do JWKS-Refresh for unknown KID")
	jwtPublicKeyFile            = flag.String("jwt_publickey", "", "JWT-Public-Key for Verifier")
	jwtVerifyAudience           = flag.String("jwt_verify_audience", "", "JWT-Audience for Verifier")
	jwtVerifyAlgorithm          = flag.String("jwt_verify_algorithm", "", "JWT-Algorithm for Verifier")
	jwtVerifyNoSuccessSiemEvent = flag.Bool("jwt_verify_no_success_siem_event", true, "Don't log success siem event from JWT-Verifier")
)

// IJWTVerifier defines the interface of JWTVerifier
type IVerifier interface {
	Verify(w http.ResponseWriter, r *http.Request, groups []string) (*CustomClaims, error)
	VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims ICustomClaims) (ICustomClaims, error)
}

// JWTVerifier provides a simple JWT verification utility
//
// Used flags:
//   - jwt_publickey Filename of JWT-Public-Key-File ()
type Verifier struct {
	log       *logger.Log
	publicKey *ecdsa.PublicKey
	jwks      *keyfunc.JWKS
	keyFunc   jwt.Keyfunc
}

var _ IVerifier = (*Verifier)(nil)

// ICustomClaims definies the generic interface for custom jwt claims
type ICustomClaims interface {
	jwt.Claims
	GetUserID() int64
	GetGroups() []string
	GetAudiences() []string
}

// CustomClaims definies the format of the JWT payload
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID int64    `json:"user_id"`
	Groups []string `json:"groups"`
}

// GetUserID returns the UserID
func (j *CustomClaims) GetUserID() int64 {
	return j.UserID
}

// GetGroups returns all Groups
func (j *CustomClaims) GetGroups() []string {
	return j.Groups
}

// GetAudiences returns the audiences
func (j *CustomClaims) GetAudiences() []string {
	return j.Audience
}

var _ ICustomClaims = (*CustomClaims)(nil)

func (j *Verifier) load() error {
	var err error

	if *jwtPublicKeyFile != "" {
		j.log.Infof("Using certificate %s for JWT verification", *jwtPublicKeyFile)

		publicKeyPem, err := os.ReadFile(*jwtPublicKeyFile)
		if err != nil {
			return err
		}

		j.publicKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyPem)
		if err != nil {
			return err
		}

		j.keyFunc = func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return j.publicKey, nil
		}
	} else if *jwksURL != "" {
		j.log.Infof("Using jwks from %s for JWT verification", *jwksURL)

		jwksOptions := keyfunc.Options{
			RefreshErrorHandler: func(err error) {
				j.log.Errorf("Error loading jwks certificates: %s", err)
			},
			RefreshInterval:   time.Second * time.Duration(*jwksRefreshInterval),
			RefreshRateLimit:  time.Second * time.Duration(*jwksRefreshRateLimit),
			RefreshTimeout:    time.Second * time.Duration(*jwksRefreshTimeout),
			RefreshUnknownKID: *jwksRefreshUnknownKID,
		}

		j.jwks, err = keyfunc.Get(*jwksURL, jwksOptions)
		if err != nil {
			return err
		}

		j.keyFunc = j.jwks.Keyfunc
	} else {
		return fmt.Errorf("neither keyfile nor jwks url specified")
	}

	return nil
}

func (j *Verifier) generateSiemEvent(r *http.Request, eventType siem.EventType, claims ICustomClaims) *siem.Event {
	evt := &siem.Event{}

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
func (j *Verifier) VerifyWithCustomClaims(w http.ResponseWriter, r *http.Request, groups []string, claims ICustomClaims) (ICustomClaims, error) {
	authorizationHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	// Check if valid Bearer-Header
	if len(authorizationHeader) != 2 || authorizationHeader[1] == "" {
		j.log.SiemEvent(
			j.generateSiemEvent(
				r,
				siem.EventTypeAuthenticationFailed,
				nil,
			),
			"Invalid or empty authorization header: %s",
			r.Header.Get("Authorization"),
		)

		return nil, fmt.Errorf("invalid authorization header")
	}

	// validate the token
	token, err := jwt.ParseWithClaims(authorizationHeader[1], claims, j.keyFunc)

	// branch out into the possible error from signing
	switch typedErr := err.(type) {
	case nil: // no error
		if !token.Valid { // but may still be invalid
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s",
				authorizationHeader[1],
			)

			return nil, fmt.Errorf("authorization failed: invalid token")
		}

		err = token.Claims.Valid()
		if err != nil {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt claims in token %s: %s",
				authorizationHeader[1],
				err,
			)

			return nil, fmt.Errorf("authorization failed: invalid claims in token: %s", err)
		}

		if token.Method.Alg() != *jwtVerifyAlgorithm {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s - missmatching algorithm: got %s, expected %s",
				authorizationHeader[1],
				token.Method.Alg(),
				*jwtVerifyAlgorithm,
			)

			return nil, fmt.Errorf("missmatching algorithm: got %s, expected %s", token.Method.Alg(), *jwtVerifyAlgorithm)
		}

		customClaims, ok := token.Claims.(ICustomClaims)
		if !ok {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s - casting jwt custom claims failed",
				*jwtVerifyAlgorithm,
			)

			return nil, fmt.Errorf("casting jwt custom claims failed")
		}

		if !gousu.ContainsString(customClaims.GetAudiences(), *jwtVerifyAudience) {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailedAttact,
					nil,
				),
				"Invalid jwt token: %s - missmatching audience: got %v, expected %s",
				authorizationHeader[1],
				customClaims.GetAudiences(),
				*jwtVerifyAudience,
			)

			return nil, fmt.Errorf("missmatching audience: got %v, expected %s", customClaims.GetAudiences(), *jwtVerifyAudience)
		}

		for i := range groups {
			if !gousu.ContainsString(customClaims.GetGroups(), groups[i]) {
				j.log.SiemEvent(
					j.generateSiemEvent(
						r,
						siem.EventTypeAuthenticationFailedAttact,
						nil,
					),
					"Authorization failed: missing group in token %s: %s",
					authorizationHeader[1],
					groups[i],
				)

				return nil, fmt.Errorf("authorization failed: missing group %s", groups[i])
			}
		}

		if !*jwtVerifyNoSuccessSiemEvent {
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationSuccess,
					customClaims,
				),
				"JWT-Verification succeded for user %d",
				customClaims.GetUserID(),
			)
		}

		return customClaims, nil

	case *jwt.ValidationError: // something was wrong during the validation
		switch typedErr.Errors {
		case jwt.ValidationErrorExpired:
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailed,
					nil,
				),
				"Authorization failed: jwt expired: %s",
				authorizationHeader[1],
			)

			return nil, fmt.Errorf("authorization failed: JWT expired")

		default:
			j.log.SiemEvent(
				j.generateSiemEvent(
					r,
					siem.EventTypeAuthenticationFailedAttact,
					nil,
				),
				"Authorization failed: invalid jwt %s: %s",
				authorizationHeader[1],
				err,
			)

			return nil, fmt.Errorf("authorization failed: invalid JWT: %s", err)
		}

	default: // something else went wrong
		j.log.SiemEvent(
			j.generateSiemEvent(
				r,
				siem.EventTypeAuthenticationFailedAttact,
				nil,
			),
			"Authorization failed for token %s: %s",
			authorizationHeader[1],
			err,
		)

		return nil, fmt.Errorf("authorization failed: %s", err)
	}
}

// Verify validates the JWT from the authorization header and checks if the
// required groups are fullfiled
//
// If the authorization fails it sends a 403 and returns
func (j *Verifier) Verify(w http.ResponseWriter, r *http.Request, groups []string) (*CustomClaims, error) {
	claims, err := j.VerifyWithCustomClaims(w, r, groups, &CustomClaims{})
	if err != nil {
		return nil, err
	}

	return claims.(*CustomClaims), nil
}

// NewVerifier creates a new initilized instance of JWTVerifier
// and loads the public key from the file specified by the flag 'jwt_publickey'
func NewVerifier() (*Verifier, error) {
	j := &Verifier{
		log: logger.GetLogger("gousujwt"),
	}

	err := j.load()
	if err != nil {
		j.log.Errorf("Error initializing jwt utiliy: %s", err)
		return nil, err
	}

	return j, nil
}
