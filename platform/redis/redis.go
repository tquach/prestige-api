package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// NewRedisPool creates a new redis connection pool.
func NewRedisPool(settings ServerSettings) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", settings.URL)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
