package gousuredis

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// XAdd adds an stream event
func (s *Service) XAdd(key string, data map[string]string) (string, error) {
	conn, err := s.openConn(true)
	if err != nil {
		return "", fmt.Errorf("can't connect to redis: %s", err)
	}
	defer conn.Close()

	args := []interface{}{
		key,
		"*",
	}

	for key, value := range data {
		args = append(args, key, value)
	}

	return redis.String(conn.Do("XADD", args...))
}

type XGroupCreateOffset = string

const (
	XGroupCreateOffsetLast  XGroupCreateOffset = "$"
	XGroupCreateOffsetFirst XGroupCreateOffset = "0"
)

// XGroupCreate adds an stream event
func (s *Service) XGroupCreate(groupName string, key string, offset XGroupCreateOffset, mkStream bool, ignoreBusy bool) error {
	conn, err := s.openConn(true)
	if err != nil {
		return fmt.Errorf("can't connect to redis: %s", err)
	}
	defer conn.Close()

	args := []interface{}{
		"CREATE",
		key,
		groupName,
		offset,
	}

	if mkStream {
		args = append(args, "MKSTREAM")
	}

	_, err = conn.Do("XGROUP", args...)

	if err != nil && strings.HasPrefix(err.Error(), "BUSYGROUP") && ignoreBusy {
		return nil
	}

	return err
}

type XReadGroupStreamID = string

const (
	XReadGroupIDStreamNew     XReadGroupStreamID = ">"
	XReadGroupIDStreamPending XReadGroupStreamID = "0"
)

type XEvent struct {
	Key  string
	ID   string
	Data map[string]string
}

// XReadGroup waits for one new item in a stream (blocking with timeout)
func (s *Service) XReadGroup(groupName string, consumerName string, key string, timeout time.Duration, streamID XReadGroupStreamID) (*XEvent, error) {
	conn, err := s.openConn(true)
	if err != nil {
		return nil, fmt.Errorf("can't connect to redis: %s", err)
	}
	defer conn.Close()

	result, err := redis.Values(conn.Do("XREADGROUP", "GROUP", groupName, consumerName, "COUNT", 1, "BLOCK", int(timeout/time.Millisecond), "STREAMS", key, streamID))
	if err != nil {
		return nil, err
	}

	if result == nil || len(result) < 1 || result[0] == nil {
		return nil, ErrNil
	}

	resultArr, err := redis.Values(result[0], nil)
	if err != nil {
		return nil, fmt.Errorf("parsing result failed: %s", err)
	}

	if len(resultArr) < 1 || resultArr[0] == nil {
		return nil, ErrNil
	}

	evt := &XEvent{}

	evt.Key, err = redis.String(resultArr[0], nil)
	if err != nil {
		return nil, fmt.Errorf("parsing key from result failed: %s", err)
	}

	resultEvents, err := redis.Values(resultArr[1], nil)
	if err != nil {
		return nil, fmt.Errorf("parsing events from result failed: %s", err)
	}

	if len(resultEvents) < 1 {
		return nil, ErrNil
	}

	resultEvent, err := redis.Values(resultEvents[0], nil)
	if err != nil {
		return nil, fmt.Errorf("parsing event from result failed: %s", err)
	}

	if len(resultEvent) < 2 {
		return nil, fmt.Errorf("malformed result event: %v", resultEvent)
	}

	evt.ID, err = redis.String(resultEvent[0], nil)
	if err != nil {
		return nil, fmt.Errorf("parsing event id from result failed: %s", err)
	}

	evt.Data, err = redis.StringMap(resultEvent[1], nil)
	if err != nil {
		return nil, fmt.Errorf("parsing event payload from result failed: %s", err)
	}

	return evt, err
}

// XAck acknowledges stream event
func (s *Service) XAck(groupName string, key string, id string) (int, error) {
	conn, err := s.openConn(true)
	if err != nil {
		return 0, fmt.Errorf("can't connect to redis: %s", err)
	}
	defer conn.Close()

	return redis.Int(conn.Do("XACK", key, groupName, id))
}
