package redis

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/keller0/yxi.io/internal"
)

var (
	Pool    *redis.Pool
	apihash = "api_keys"
)

func init() {
	redisHost := internal.GetEnv("REDIS_ADDR", ":6379")
	redisPass := internal.GetEnv("REDIS_PASS", "")
	options := redis.DialPassword(redisPass)

	Pool = newPool(redisHost, options)
	cleanupHook()
}

func newPool(server string, options redis.DialOption) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, options)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
