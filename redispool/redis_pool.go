// copyright

// Package redispool shall be used over twemproxy
package redispool

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

type RedisPoolConfig struct {
	RedisAddrs                         []string
	ConnTimeout                        time.Duration
	ReadTimeout                        time.Duration
	WriteTimeout                       time.Duration
	MaxIdleConnNumEach, MaxConnNumEach int
	IdleTimeout                        time.Duration
}

type RedisPool struct {
	config *RedisPoolConfig
	rcp    []*redis.Pool
}

func NewRedisPool(config *RedisPoolConfig) (*RedisPool, error) {
	rp := &RedisPool{
		config: config,
	}
	err := rp.init()
	if err != nil {
		return nil, err
	}
	return rp, nil
}

func (rp *RedisPool) init() error {
	for _, addr := range rp.config.RedisAddrs {
		rd := redisDialer{
			addr:         addr,
			connTimeout:  rp.config.ConnTimeout,
			readTimeout:  rp.config.ReadTimeout,
			writeTimeout: rp.config.WriteTimeout,
		}
		p := &redis.Pool{
			MaxIdle:      rp.config.MaxIdleConnNumEach,
			MaxActive:    rp.config.MaxConnNumEach,
			IdleTimeout:  rp.config.IdleTimeout,
			Dial:         rd.dial,
			TestOnBorrow: rd.testOnBorrow,
		}
		rp.rcp = append(rp.rcp, p)
	}
	if len(rp.rcp) == 0 {
		return errors.New("no addrs in RedisPoolConfig")
	}
	return nil
}

func (rp *RedisPool) Get() (redis.Conn, error) {
	pn := len(rp.rcp)
	for i, j := rand.Int31n(int32(pn)), 0; j < pn; j++ {
		p := rp.rcp[i]
		c := p.Get()
		if c.Err() == nil {
			return c, nil
		}
		i = (i + 1) % int32(pn)
	}
	return nil, errors.New("no valid conn in RedisPool")
}

func (rp *RedisPool) Put(c redis.Conn) {
	c.Close()
}

func (rp *RedisPool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c, err := rp.Get()
	if err != nil {
		return nil, err
	}
	defer rp.Put(c)
	return c.Do(commandName, args...)
}

func (rp *RedisPool) Close() {
	rp.rcp.Close()
}
