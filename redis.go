package hera

import (
	"fmt"
	"github.com/xcodecraft/hera/redis"
	"time"
)

type RedisSvc struct {
	connType string
	ip       string
	port     string
	db       string
	pwd      string
	pool     *redis.Pool
}

func NewRedis(dbIp, dbPort, dbNo, dbPwd string, enablePool bool) *RedisSvc {
	var dbPool *redis.Pool = nil
	if enablePool {
		dbPool = &redis.Pool{
			MaxIdle:     80,
			MaxActive:   1000,
			IdleTimeout: 240 * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", dbIp, dbPort))
				if err != nil {
					return nil, err
				}

				if dbPwd != "" {
					_, err = c.Do("AUTH", dbPwd)
					if err != nil {
						c.Close()
						return nil, err
					}
				}

				if dbNo != "" {
					_, err = c.Do("SELECT", dbNo)
					if err != nil {
						c.Close()
						return nil, err
					}
				}

				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}
	return &RedisSvc{
		connType: "tcp",
		ip:       dbIp,
		port:     dbPort,
		db:       dbNo,
		pwd:      dbPwd,
		pool:     dbPool,
	}
}

func (this *RedisSvc) getRedisPool() redis.Conn {
	if this.pool != nil {
		return this.pool.Get()
	}
	return nil
}

func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

func (this *RedisSvc) ActiveCount() int {
	if this.pool != nil {
		return this.pool.ActiveCount()
	}
	return 0
}

func (this *RedisSvc) DoCmd(cmd string, args ...interface{}) (interface{}, error) {
	var c redis.Conn
	var err error
	if this.pool != nil {
		c = this.getRedisPool()
	} else {
		c, err = redis.Dial(this.connType, fmt.Sprintf("%s:%s", this.ip, this.port))
		if err != nil {
			return nil, err
		}
	}
	defer c.Close()
	re, err := c.Do(cmd, args...)
	if err != nil {
		return nil, err
	}
	return re, err
}
