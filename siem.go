package gousu

import "gopkg.in/guregu/null.v4"

// Fields used in logs for siem events
const (
	SiemEventFieldType           = "siem_type"
	SiemEventFieldLevel          = "siem_level"
	SiemEventFieldUserIdentifier = "siem_user_identifier"
	SiemEventFieldSourceIP       = "siem_source_ip"
	SiemEventFieldSourceRealIP   = "siem_source_real_ip"
)

// SiemEventType specifies the type of the siem event
type SiemEventType = string

// All types of our siem events
const (
	SiemEventTypeLoginSuccess               SiemEventType = "login_success"
	SiemEventTypeLoginFailed                SiemEventType = "login_failed"
	SiemEventTypeLoginFailedAttact          SiemEventType = "login_failed_attact"
	SiemEventTypeLogoutSuccess              SiemEventType = "logout_success"
	SiemEventTypeLogoutFailed               SiemEventType = "logout_failed"
	SiemEventTypeUserCreated                SiemEventType = "user_created"
	SiemEventTypeUserDeleted                SiemEventType = "user_deleted"
	SiemEventTypeUserLocked                 SiemEventType = "user_locked"
	SiemEventTypeUserUnlocked               SiemEventType = "user_unlocked"
	SiemEventTypeAuthenticationSuccess      SiemEventType = "authentication_success"
	SiemEventTypeAuthenticationFailed       SiemEventType = "authentication_failed"
	SiemEventTypeAuthenticationFailedAttact SiemEventType = "authentication_failed_attact"
)

// SiemEventLevel specifies the level of the siem event
type SiemEventLevel string

// All levels of our siem events
const (
	SiemEventLevelInfo     SiemEventLevel = "info"
	SiemEventLevelWarn     SiemEventLevel = "warn"
	SiemEventLevelCritical SiemEventLevel = "critical"
)

// SiemEventLevels is a map of the matching level for each siem event type
var SiemEventLevels = map[SiemEventType]SiemEventLevel{
	SiemEventTypeLoginSuccess:               SiemEventLevelInfo,
	SiemEventTypeLoginFailed:                SiemEventLevelWarn,
	SiemEventTypeLoginFailedAttact:          SiemEventLevelCritical,
	SiemEventTypeLogoutSuccess:              SiemEventLevelInfo,
	SiemEventTypeLogoutFailed:               SiemEventLevelWarn,
	SiemEventTypeUserCreated:                SiemEventLevelInfo,
	SiemEventTypeUserDeleted:                SiemEventLevelInfo,
	SiemEventTypeUserLocked:                 SiemEventLevelInfo,
	SiemEventTypeUserUnlocked:               SiemEventLevelInfo,
	SiemEventTypeAuthenticationSuccess:      SiemEventLevelInfo,
	SiemEventTypeAuthenticationFailed:       SiemEventLevelWarn,
	SiemEventTypeAuthenticationFailedAttact: SiemEventLevelCritical,
}

// SiemEvent is the basic struct for siem events
type SiemEvent struct {
	Type           SiemEventType
	UserIdentifier null.String
	SourceIP       null.String
	SourceRealIP   null.String
}

// Level retuns the level to the siem event type (defaults to "critical")
func (s *SiemEvent) Level() SiemEventLevel {
	level, ok := SiemEventLevels[s.Type]
	if !ok {
		level = SiemEventLevelCritical
	}

	return level
}
