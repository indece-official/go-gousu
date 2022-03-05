package logger

import (
	"fmt"
	"strings"

	"github.com/chakrit/go-bunyan"
	"github.com/indece-official/go-gousu/v2/gousu/siem"
	"github.com/namsral/flag"
)

var (
	loglevel    = flag.String("loglevel", "INFO", "")
	siemEnabled = flag.Bool("siem_enabled", true, "")
	logDisabled = false
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

// SiemEvent logs an siem event
func (l *Log) SiemEvent(event *siem.Event, msg string, args ...interface{}) {
	if !*siemEnabled {
		return
	}

	log := l.Log.
		Record(siem.EventFieldType, event.Type).
		Record(siem.EventFieldLevel, event.Level())

	if event.UserIdentifier.Valid {
		log = log.Record(siem.EventFieldUserIdentifier, event.UserIdentifier)
	}

	if event.SourceIP.Valid {
		log = log.Record(siem.EventFieldSourceIP, event.SourceIP)
	}

	if event.SourceRealIP.Valid {
		log = log.Record(siem.EventFieldSourceRealIP, event.SourceRealIP)
	}

	switch event.Level() {
	case siem.EventLevelInfo:
		log.Infof(msg, args...)
	case siem.EventLevelWarn:
		log.Warnf(msg, args...)
	case siem.EventLevelCritical:
		log.Errorf(msg, args...)
	default:
		log.Errorf(msg, args...)
	}
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

// DisableLogger disabled the logger (must be called before InitLogger())
func DisableLogger() {
	logDisabled = true
}

// InitLogger initializes the parent logger and sets the project's name
func InitLogger(projectName string) {
	level, ok := mapLevels[strings.ToUpper(*loglevel)]
	if !ok {
		level = bunyan.INFO
	}

	var sink bunyan.Sink

	if logDisabled {
		sink = bunyan.NilSink()
	} else {
		sink = bunyan.FilterSink(level, bunyan.StdoutSink())
	}

	parentLogger = &Log{bunyan.NewStdLogger(projectName, sink)}

	if !*siemEnabled {
		parentLogger.Warnf("SIEM-Event logging is disabled")
	}
}

// GetLogger returns a logger for a specific component
func GetLogger(componentName string) *Log {
	if parentLogger == nil {
		InitLogger("test")
	}

	return &Log{parentLogger.Record("component", componentName)}
}
