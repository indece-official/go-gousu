package goususqlite3

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/namsral/flag"

	// Use sqlite3 driver for database/sql
	_ "github.com/mattn/go-sqlite3"
)

// ServiceName defines the name of sqlite3 service used for dependency injection
const ServiceName = "sqlite3"

var (
	sqlite3Filename      = flag.String("sqlite3_filename", "db.sqlite", "SQLite3 file name")
	sqlite3Cache         = flag.String("sqlite3_cache", "private", "SQLite3 cache mode: shared | private")
	sqlite3Mode          = flag.String("sqlite3_mode", "rwc", "SQLite3 access mode: ro | rw | rwc | memory")
	sqlite3MaxRetries    = flag.Int("sqlite3_max_retries", 10, "")
	sqlite3RetryInterval = flag.Int("sqlite3_retry_interval", 6, "")
	sqlite3MaxIdleConns  = flag.Int("sqlite3_max_idle_conns", 0, "")
	sqlite3MaxOpenConns  = flag.Int("sqlite3_max_open_conns", 0, "")
)

// Options can contain parameters passed to the sqlite3 service
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

// IService defined the interface of the sqlite3 database service
type IService interface {
	gousu.IService

	GetDB() *sql.DB
	GetDBSafe() (*sql.DB, error)
}

// Service provides the interaction with the sqlite3 database
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
		s.log.Infof("Disconnecting from sqlite3 database %s ...", *sqlite3Filename)

		s.db.Close()
		s.db = nil
	}

	connStr := fmt.Sprintf("file:%s?cache=%s&mode=%s", *sqlite3Filename, *sqlite3Cache, *sqlite3Mode)

	retries := 0

	for retries < *sqlite3MaxRetries {
		s.log.Infof("Connecting to sqlite3 database %s ...", *sqlite3Filename)

		s.db, err = openFunc("sqlite3", connStr)
		if err == nil {
			s.db.SetMaxIdleConns(*sqlite3MaxIdleConns)
			s.db.SetMaxOpenConns(*sqlite3MaxOpenConns)

			err = s.db.Ping()
			if err == nil {
				s.log.Infof("Connected to sqlite3 database %s", *sqlite3Filename)

				return nil
			}
		}

		s.log.Errorf("Can't connect to sqlite3 database %s: %s", *sqlite3Filename, err)

		time.Sleep(time.Second * time.Duration(*sqlite3RetryInterval))
		retries++
	}

	return err
}

// GetDB returns the sqlite3 db connection
//
// Deprecated: Use GetDBSafe() instead
func (s *Service) GetDB() *sql.DB {
	return s.db
}

// GetDBSafe returns the sqlite3 db connection after verifying it is alive
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

// Start initializes the connection to the sqlite3 database and executed both setup.sql and update.sql
// after connecting
func (s *Service) Start() error {
	var err error

	s.error = s.connect()

	if s.error != nil {
		s.log.Errorf("Can't connect to sqlite3 database %s after %d attempts: %s", *sqlite3Filename, *sqlite3MaxRetries, err)

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

// Health checks the health of the sqlite3-service by pinging the sqlite3 database
func (s *Service) Health() error {
	if s.error != nil {
		return s.error
	}

	return s.db.Ping()
}

// NewServiceBase creates a new instance of sqlite3-service, should be used instead
//
//	of generating it manually
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
