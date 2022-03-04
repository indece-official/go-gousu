package gousuredis

// From https://github.com/go-redsync/redsync/

import (
	"context"
	"fmt"
	"strings"
	"time"

	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
)

type redsyncPool struct {
	cluster   *redisc.Cluster
	pool      *redis.Pool
	withRetry bool
}

func (p *redsyncPool) Get(ctx context.Context) (redsyncredis.Conn, error) {
	var c redis.Conn
	var err error

	if p.cluster != nil {
		c = p.cluster.Get()
	} else {
		if ctx != nil {
			c, err = p.pool.GetContext(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			c = p.pool.Get()
		}
	}

	if p.withRetry {
		rc, err := redisc.RetryConn(c, 3, 100*time.Millisecond)
		if err != nil {
			return nil, fmt.Errorf("retry failed: %s", err)
		}

		return &conn{rc}, nil
	}

	return &conn{c}, nil
}

func newRedsyncPoolFromPool(p *redis.Pool) redsyncredis.Pool {
	return &redsyncPool{
		cluster:   nil,
		pool:      p,
		withRetry: false,
	}
}

func newRedsyncPoolFromCluster(c *redisc.Cluster) redsyncredis.Pool {
	return &redsyncPool{
		pool:      nil,
		cluster:   c,
		withRetry: true,
	}
}

type conn struct {
	delegate redis.Conn
}

func (c *conn) Get(name string) (string, error) {
	value, err := redis.String(c.delegate.Do("GET", name))
	return value, noErrNil(err)
}

func (c *conn) Set(name string, value string) (bool, error) {
	reply, err := redis.String(c.delegate.Do("SET", name, value))
	return reply == "OK", noErrNil(err)
}

func (c *conn) SetNX(name string, value string, expiry time.Duration) (bool, error) {
	reply, err := redis.String(c.delegate.Do("SET", name, value, "NX", "PX", int(expiry/time.Millisecond)))
	return reply == "OK", noErrNil(err)
}

func (c *conn) PTTL(name string) (time.Duration, error) {
	expiry, err := redis.Int64(c.delegate.Do("PTTL", name))
	return time.Duration(expiry) * time.Millisecond, noErrNil(err)
}

func (c *conn) Eval(script *redsyncredis.Script, keysAndArgs ...interface{}) (interface{}, error) {
	v, err := c.delegate.Do("EVALSHA", args(script, script.Hash, keysAndArgs)...)
	if e, ok := err.(redis.Error); ok && strings.HasPrefix(string(e), "NOSCRIPT ") {
		v, err = c.delegate.Do("EVAL", args(script, script.Src, keysAndArgs)...)
	}
	return v, noErrNil(err)
}

func (c *conn) Close() error {
	err := c.delegate.Close()
	return noErrNil(err)
}

func noErrNil(err error) error {
	if err == redis.ErrNil {
		return nil
	}

	return err
}

func args(script *redsyncredis.Script, spec string, keysAndArgs []interface{}) []interface{} {
	var args []interface{}
	if script.KeyCount < 0 {
		args = make([]interface{}, 1+len(keysAndArgs))
		args[0] = spec
		copy(args[1:], keysAndArgs)
	} else {
		args = make([]interface{}, 2+len(keysAndArgs))
		args[0] = spec
		args[1] = script.KeyCount
		copy(args[2:], keysAndArgs)
	}
	return args
}
