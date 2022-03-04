package gousuldap

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	ldapv3 "github.com/go-ldap/ldap/v3"

	"github.com/indece-official/go-gousu/gousu"
	"github.com/indece-official/go-gousu/gousu/logger"
	"github.com/namsral/flag"
)

var (
	ldapHost          = flag.String("ldap_host", "localhost", "")
	ldapPort          = flag.Int("ldap_port", 389, "")
	ldapMaxRetries    = flag.Int("ldap_max_retries", 5, "")
	ldapRetryInterval = flag.Int("ldap_retry_interval", 6, "")
	ldapBindUser      = flag.String("ldap_binduser", "", "")
	ldapBindPassword  = flag.String("ldap_bindpassword", "", "")
	ldapFilterPattern = flag.String("ldap_filterpattern", "", "")
	ldapLoginPattern  = flag.String("ldap_loginpattern", "", "")
	ldapBaseDn        = flag.String("ldap_basedn", "", "")
)

// ServiceName defines the name of ldap service used for dependency injection
const ServiceName = "ldap"

// IService defines all public functions of the ldap service
type IService interface {
	gousu.IService

	Reconnect() error
	SimpleLogin(username string, password string, attributes []string) (*map[string][]string, error)
}

// Service is the basic struct for the ldap service
type Service struct {
	error         error
	log           *logger.Log
	ldapFilterTpl *template.Template
	ldapLoginTpl  *template.Template
	conn          *ldapv3.Conn
	retries       int
	reconnecting  bool
}

var _ IService = (*Service)(nil)

func (s *Service) connect() error {
	var err error

	s.retries = 0

	connStr := fmt.Sprintf("%s:%d", *ldapHost, *ldapPort)

	for s.retries < *ldapMaxRetries {
		s.log.Infof("Connecting to ldap on %s:%d ...", *ldapHost, *ldapPort)

		s.conn, err = ldapv3.Dial("tcp", connStr)
		if err == nil {
			s.log.Infof("Connected to ldap on %s:%d", *ldapHost, *ldapPort)

			s.retries = 0
			s.error = nil

			return nil
		}

		s.log.Errorf("Can't connect to ldap on %s:%d: %s", *ldapHost, *ldapPort, err)

		time.Sleep(time.Second * time.Duration(*ldapRetryInterval))
		s.retries++
	}

	s.log.Errorf("Can't connect to ldap on %s:%d after %d attempts: %s", *ldapHost, *ldapPort, s.retries, err)

	s.error = err

	return s.error
}

// Reconnect reconnects the ldap connection
func (s *Service) Reconnect() error {
	if s.reconnecting {
		return fmt.Errorf("already reconnecting")
	}

	s.reconnecting = true

	if s.conn != nil {
		s.log.Infof("Closing ldap connection")

		s.conn.Close()

		s.log.Infof("Closed ldap connection")
	}

	s.log.Infof("Reconnecting to ldap ...")

	err := s.connect()

	s.reconnecting = false

	return err
}

// Name returns the name of ldap service from ServiceName
func (s *Service) Name() string {
	return ServiceName
}

// Start starts the ldap service by compiling the ldap patterns and etablishing the ldap connection
func (s *Service) Start() error {
	var err error

	s.ldapFilterTpl, err = template.New("ldap_filter").Parse(*ldapFilterPattern)
	if err != nil {
		return s.log.ErrorfX("Error parsing ldap filter pattern: %s", err)
	}

	s.ldapLoginTpl, err = template.New("ldap_login").Parse(*ldapLoginPattern)
	if err != nil {
		return s.log.ErrorfX("Error parsing ldap login pattern: %s", err)
	}

	s.connect()
	if err != nil {
		return err
	}

	return nil
}

// Health checks if the ldap connection is healthy by executing a bind against ldap
// If the connectuing is not healthy a reconnect is triggered
func (s *Service) Health() error {
	if s.error == nil {
		s.error = s.conn.Bind(*ldapBindUser, *ldapBindPassword)
		if s.error != nil {
			s.log.Errorf("Can't bind to ldap with user '%s': %s", *ldapBindUser, s.error)
		}
	}

	if s.error != nil && ldapv3.IsErrorWithCode(s.error, ldapv3.ErrorNetwork) && s.retries == 0 && !s.reconnecting {
		go s.Reconnect()
	}

	return s.error
}

// Stop closes the ldap connection
func (s *Service) Stop() error {
	if s.conn != nil {
		s.conn.Close()
	}

	return nil
}

// SimpleLogin check a user against ldap and executes a bind with it's credentials
//
// All attributes requested are returned for the matching user, else an error is returned
func (s *Service) SimpleLogin(username string, password string, attributes []string) (*map[string][]string, error) {
	var err error

	for i := 0; i < 2; i++ {
		// Bind with the read only user
		err = s.conn.Bind(*ldapBindUser, *ldapBindPassword)
		if err != nil {
			errExt := s.log.ErrorfX("Can't bind to ldap with user '%s': %s", *ldapBindUser, err)

			if ldapv3.IsErrorWithCode(s.error, ldapv3.ErrorNetwork) {
				s.Reconnect()

				continue
			}

			return nil, errExt
		}

		break
	}

	// Escape username
	username = ldapv3.EscapeFilter(username)

	filterBuf := new(bytes.Buffer)
	err = s.ldapFilterTpl.Execute(filterBuf, map[string]string{
		"username": username,
	})
	if err != nil {
		return nil, s.log.ErrorfX("Error preparing ldap filter pattern: %s", err)
	}

	// Search for the given username
	searchRequest := ldapv3.NewSearchRequest(
		*ldapBaseDn,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases, 0, 0, false,
		filterBuf.String(),
		attributes,
		nil,
	)

	sr, err := s.conn.Search(searchRequest)
	if err != nil {
		return nil, s.log.ErrorfX("Error executing ldap search query: %s", err)
	}

	if len(sr.Entries) != 1 {
		s.log.Infof("LDAP-User not found")

		return nil, nil
	}

	user := map[string][]string{}

	for i := range sr.Entries[0].Attributes {
		if len(sr.Entries[0].Attributes[i].Values) == 0 {
			continue
		}

		user[sr.Entries[0].Attributes[i].Name] = sr.Entries[0].Attributes[i].Values
	}

	loginBuf := new(bytes.Buffer)
	err = s.ldapLoginTpl.Execute(loginBuf, map[string]string{
		"username": username,
	})
	if err != nil {
		return nil, s.log.ErrorfX("Error preparing ldap login pattern: %s", err)
	}

	// Bind as the user to verify their password
	err = s.conn.Bind(loginBuf.String(), password)
	if err != nil {
		s.log.Infof("Error ldap login failed: %s", err)

		return nil, nil
	}

	return &user, nil
}

// NewService is the ServiceFactory for ldap service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

// Assert NewService fullfills gousu.ServiceFactory
var _ (gousu.ServiceFactory) = NewService
