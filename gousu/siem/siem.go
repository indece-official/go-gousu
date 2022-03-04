package siem

import "gopkg.in/guregu/null.v4"

// Fields used in logs for siem events
const (
	EventFieldType           = "siem_type"
	EventFieldLevel          = "siem_level"
	EventFieldUserIdentifier = "siem_user_identifier"
	EventFieldSourceIP       = "siem_source_ip"
	EventFieldSourceRealIP   = "siem_source_real_ip"
)

// EventType specifies the type of the siem event
type EventType = string

// All types of our siem events
const (
	EventTypeLoginSuccess               EventType = "login_success"
	EventTypeLoginFailed                EventType = "login_failed"
	EventTypeLoginFailedAttact          EventType = "login_failed_attact"
	EventTypeLogoutSuccess              EventType = "logout_success"
	EventTypeLogoutFailed               EventType = "logout_failed"
	EventTypeUserCreated                EventType = "user_created"
	EventTypeUserDeleted                EventType = "user_deleted"
	EventTypeUserLocked                 EventType = "user_locked"
	EventTypeUserUnlocked               EventType = "user_unlocked"
	EventTypeAuthenticationSuccess      EventType = "authentication_success"
	EventTypeAuthenticationFailed       EventType = "authentication_failed"
	EventTypeAuthenticationFailedAttact EventType = "authentication_failed_attact"
)

// EventLevel specifies the level of the siem event
type EventLevel string

// All levels of our siem events
const (
	EventLevelInfo     EventLevel = "info"
	EventLevelWarn     EventLevel = "warn"
	EventLevelCritical EventLevel = "critical"
)

// EventLevels is a map of the matching level for each siem event type
var EventLevels = map[EventType]EventLevel{
	EventTypeLoginSuccess:               EventLevelInfo,
	EventTypeLoginFailed:                EventLevelWarn,
	EventTypeLoginFailedAttact:          EventLevelCritical,
	EventTypeLogoutSuccess:              EventLevelInfo,
	EventTypeLogoutFailed:               EventLevelWarn,
	EventTypeUserCreated:                EventLevelInfo,
	EventTypeUserDeleted:                EventLevelInfo,
	EventTypeUserLocked:                 EventLevelInfo,
	EventTypeUserUnlocked:               EventLevelInfo,
	EventTypeAuthenticationSuccess:      EventLevelInfo,
	EventTypeAuthenticationFailed:       EventLevelWarn,
	EventTypeAuthenticationFailedAttact: EventLevelCritical,
}

// Event is the basic struct for siem events
type Event struct {
	Type           EventType
	UserIdentifier null.String
	SourceIP       null.String
	SourceRealIP   null.String
}

// Level retuns the level to the siem event type (defaults to "critical")
func (s *Event) Level() EventLevel {
	level, ok := EventLevels[s.Type]
	if !ok {
		level = EventLevelCritical
	}

	return level
}
