// copyright

package redispool

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisDialer struct {
	addr         string
	connTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (rd *redisDialer) dial() (redis.Conn, error) {
	return redis.DialTimeout("tcp", rd.addr, rd.connTimeout, rd.readTimeout, rd.writeTimeout)
}

func (rd *redisDialer) testOnBorrow(c redis.Conn, t time.Time) error {
	_, err := c.Do("PING")
	return err
}
