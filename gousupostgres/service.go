package gousupostgres

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/namsral/flag"

	// Use postgres driver for database/sql
	_ "github.com/lib/pq"
)

// ServiceName defines the name of postgres service used for dependency injection
const ServiceName = "postgres"

var (
	postgresHost          = flag.String("postgres_host", "localhost", "")
	postgresPort          = flag.Int("postgres_port", 5432, "")
	postgresUser          = flag.String("postgres_user", "", "")
	postgresPassword      = flag.String("postgres_password", "", "")
	postgresDatabase      = flag.String("postgres_database", "", "")
	postgresMaxRetries    = flag.Int("postgres_max_retries", 10, "")
	postgresRetryInterval = flag.Int("postgres_retry_interval", 6, "")
	postgresMaxIdleConns  = flag.Int("postgres_max_idle_conns", 0, "")
	postgresMaxOpenConns  = flag.Int("postgres_max_open_conns", 0, "")
)

// Options can contain parameters passed to the postgres service
type Options struct {
	// SetupSQL can contain the content of a sql-file for updating the
	// database on startup
	SetupSQL string

	// UpdateSQL can contain the content of a sql-file for updating the
	// database on startup
	UpdateSQL string

	// GetDBRevisionSQL can be used for retrieving the revision of the database
	// used, must return/select one integer field
	GetDBRevisionSQL string

	// OpenFunc can be used to override the default sql.Open
	OpenFunc func(driverName string, dataSourceName string) (*sql.DB, error)
}

// IService defined the interface of the postgresql database service
type IService interface {
	gousu.IService

	GetDB() *sql.DB
	GetDBSafe() (*sql.DB, error)
}

// Service provides the interaction with the postgresql database
type Service struct {
	error                error
	log                  *logger.Log
	db                   *sql.DB
	options              *Options
	waitGroupReconnected sync.WaitGroup
	reconnecting         bool
}

var _ IService = (*Service)(nil)

// Name returns the name of redis service from ServiceName
func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) connect() error {
	var err error

	openFunc := sql.Open

	if s.options != nil && s.options.OpenFunc != nil {
		openFunc = s.options.OpenFunc
	}

	if s.reconnecting {
		s.waitGroupReconnected.Wait()
		if s.db == nil {
			return fmt.Errorf("no connection accomplished")
		}

		return nil
	}

	s.reconnecting = true
	s.waitGroupReconnected.Add(1)
	defer func() {
		s.reconnecting = false
		s.waitGroupReconnected.Done()
	}()

	if s.db != nil {
		s.log.Infof("Disconnecting from postgres database on %s:%d ...", *postgresHost, *postgresPort)

		s.db.Close()
		s.db = nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", *postgresUser, *postgresPassword, *postgresHost, *postgresPort, *postgresDatabase)

	retries := 0

	for retries < *postgresMaxRetries {
		s.log.Infof("Connecting to postgres database on %s:%d ...", *postgresHost, *postgresPort)

		s.db, err = openFunc("postgres", connStr)
		if err == nil {
			s.db.SetMaxIdleConns(*postgresMaxIdleConns)
			s.db.SetMaxOpenConns(*postgresMaxOpenConns)

			err = s.db.Ping()
			if err == nil {
				s.log.Infof("Connected to postgres database on %s:%d", *postgresHost, *postgresPort)

				return nil
			}
		}

		s.log.Errorf("Can't connect to postgres on %s:%d: %s", *postgresHost, *postgresPort, err)

		time.Sleep(time.Second * time.Duration(*postgresRetryInterval))
		retries++
	}

	return err
}

// GetDB returns the postgres db connection
//
// Deprecated: Use GetDBSafe() instead
func (s *Service) GetDB() *sql.DB {
	return s.db
}

// GetDBSafe returns the postgres db connection after verifying it is alive
func (s *Service) GetDBSafe() (*sql.DB, error) {
	var err error

	if s.db == nil {
		err = s.connect()
		if err != nil {
			return nil, err
		}

		return s.db, nil
	}

	err = s.db.Ping()
	if err != nil {
		err = s.connect()
		if err != nil {
			return nil, err
		}

		return s.db, nil
	}

	return s.db, nil
}

// Start initializes the connection to the postgres database and executed both setup.sql and update.sql
// after connecting
func (s *Service) Start() error {
	var err error

	s.error = s.connect()

	if s.error != nil {
		s.log.Errorf("Can't connect to postgres on %s:%d after %d attempts: %s", *postgresHost, *postgresPort, *postgresMaxRetries, err)

		return s.error
	}

	if s.options.SetupSQL != "" {
		s.log.Infof("Executing setup SQL ...")

		_, err = s.db.Exec(s.options.SetupSQL)
		if err != nil {
			s.log.Errorf("Error executing setup SQL: %s", err)

			return err
		}
	}

	if s.options.UpdateSQL != "" {
		s.log.Infof("Executing update SQL ...")

		_, err = s.db.Exec(s.options.UpdateSQL)
		if err != nil {
			s.log.Errorf("Error executing update SQL: %s", err)

			return err
		}
	}

	if s.options.GetDBRevisionSQL != "" {
		var rev int
		err = s.db.QueryRow(s.options.GetDBRevisionSQL).Scan(&rev)
		if err != nil {
			s.log.Errorf("Retrieving revision from database failed: %s", err)

			return err
		}

		s.log.Infof("Using database rev.%d", rev)
	}

	return nil
}

// Stop currently does nothing
func (s *Service) Stop() error {
	return nil
}

// Health checks the health of the postgres-service by pinging the postgres database
func (s *Service) Health() error {
	if s.error != nil {
		return s.error
	}

	return s.db.Ping()
}

// NewServiceBase creates a new instance of postgres-service, should be used instead
//  of generating it manually
func NewServiceBase(ctx gousu.IContext, options *Options) *Service {
	if options == nil {
		options = &Options{}
	}

	return &Service{
		options:      options,
		log:          logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
		reconnecting: false,
	}
}
