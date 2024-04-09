package gousuchi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
)

type contextKey string

const contextKeyExtras contextKey = "extras"

type AbstractController struct {
	regexpSanitizeString *regexp.Regexp
	log                  *logger.Log
	server               *http.Server
	router               chi.Router
	tlsConfig            *tls.Config
	host                 string
	port                 int
	error                error
}

type HandlerFunction func(w http.ResponseWriter, r *http.Request) IResponse

func (c *AbstractController) Wrap(clb HandlerFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := clb(w, r)

		log := c.GetLog(resp.GetRequest())

		err := resp.Write(w)
		if err != nil {
			err.Write(w)

			err.Log(log)

			return
		}

		resp.Log(log)
	}
}

func (c *AbstractController) UseRouter(router chi.Router) {
	c.router = router
}

func (c *AbstractController) UsePort(port int) {
	c.port = port
}

func (c *AbstractController) UseHost(host string) {
	c.host = host
}

func (c *AbstractController) UseTlsConfig(tlsConfig *tls.Config) {
	c.tlsConfig = tlsConfig
}

func (c *AbstractController) WithExtra(r *http.Request, key string, value interface{}) *http.Request {
	ctx := r.Context()

	extras, ok := ctx.Value(contextKeyExtras).(map[string]interface{})
	if extras == nil && !ok {
		extras := map[string]interface{}{}

		extras[key] = value

		return r.WithContext(context.WithValue(ctx, contextKeyExtras, extras))
	}

	extras[key] = value

	return r
}

func (c *AbstractController) sanitizeHeaderString(str string, maxLength int) string {
	str = string(c.regexpSanitizeString.ReplaceAll([]byte(str), []byte{}))

	if len(str) > maxLength {
		return str[:maxLength-1]
	}

	return str
}

func (c *AbstractController) GetLog(r *http.Request) *logger.Log {
	log := c.log

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		realIP = c.sanitizeHeaderString(realIP, 32)
		log = log.RecordX("x_real_ip", realIP)
	}

	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		forwardedFor = c.sanitizeHeaderString(forwardedFor, 32)
		log = log.RecordX("x_forwarded_for", forwardedFor)
	}

	userAgentID := r.Header.Get("X-User-Agent-ID")
	if userAgentID != "" {
		userAgentID = c.sanitizeHeaderString(userAgentID, 70)
		log = log.RecordX("x_user_agent_id", userAgentID)
	}

	extras, ok := r.Context().Value(contextKeyExtras).(map[string]interface{})
	if extras != nil && ok {
		for key, val := range extras {
			log = log.RecordX(key, val)
		}
	}

	return log
}

// Start starts the api server in a new go-func
func (c *AbstractController) Start() error {
	c.error = nil

	go func() {
		c.server = &http.Server{
			Addr:      fmt.Sprintf("%s:%d", c.host, c.port),
			Handler:   c.router,
			TLSConfig: c.tlsConfig,
		}

		err := c.server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				c.error = err

				c.log.Errorf("Can't start server: %s", err)
			}

			c.server = nil
		}
	}()

	c.log.Infof("Server listening on %s:%d", c.host, c.port)

	return nil
}

// Health checks if the api server has thrown unresolvable internal errors
func (c *AbstractController) Health() error {
	return c.error
}

// Stop currently does nothing
func (c *AbstractController) Stop() error {
	if c.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return c.server.Shutdown(ctx)
}

func NewAbstractController(
	log *logger.Log,
) *AbstractController {
	return &AbstractController{
		router:               chi.NewRouter(),
		log:                  log,
		regexpSanitizeString: regexp.MustCompile(`[\n\t\"\']+`),
		host:                 "",
		port:                 0,
		tlsConfig:            nil,
	}
}
