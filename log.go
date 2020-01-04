package gousu

import (
	"fmt"
	"strings"

	"github.com/chakrit/go-bunyan"
	"github.com/namsral/flag"
)

var (
	loglevel = flag.String("loglevel", "INFO", "")
)

var (
	mapLevels = map[string]bunyan.Level{
		"EVERYTHING": bunyan.EVERYTHING,
		"TRACE":      bunyan.TRACE,
		"DEBUG":      bunyan.DEBUG,
		"INFO":       bunyan.INFO,
		"WARN":       bunyan.WARN,
		"ERROR":      bunyan.ERROR,
		"FATAL":      bunyan.FATAL,
	}
)

// Log provides the base structure for a extended logger
//
// The loglevel can be controller via the config property "loglevel" or the environment
// variable LOGLEVEL
type Log struct {
	bunyan.Log
}

// ErrorfX logs an error and returns it
func (l *Log) ErrorfX(msg string, args ...interface{}) error {
	l.Log.Errorf(msg, args...)

	return fmt.Errorf(msg, args...)
}

// RecordX returns a new Logger with a specified Record assigned
func (l *Log) RecordX(key string, value interface{}) *Log {
	return &Log{
		l.Log.Record(key, value),
	}
}

// RecordfX returns a new Logger with a specified formatted Record assigned
func (l *Log) RecordfX(key, value string, args ...interface{}) *Log {
	return &Log{
		l.Log.Recordf(key, value, args...),
	}
}

var parentLogger *Log

// InitLogger initializes the parent logger and sets the project's name
func InitLogger(projectName string) {
	level, ok := mapLevels[strings.ToUpper(*loglevel)]
	if !ok {
		level = bunyan.INFO
	}

	sink := bunyan.FilterSink(level, bunyan.StdoutSink())

	parentLogger = &Log{bunyan.NewStdLogger(projectName, sink)}
}

// GetLogger returns a logger for a specific component
func GetLogger(componentName string) *Log {
	if parentLogger == nil {
		InitLogger("test")
	}

	return &Log{parentLogger.Record("component", componentName)}
}
