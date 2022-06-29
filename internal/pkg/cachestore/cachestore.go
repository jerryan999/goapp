package cachestore

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	// ErrCacheMiss is the error returned when the requested item is not available in cache
	ErrCacheMiss = errors.New("not found in cache")
	// ErrCacheNotInitialized is the error returned when the cache handler is not initialized
	ErrCacheNotInitialized = errors.New("not initialized")
)

// Config holds all the configuration required for this package
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	StoreName string `json:"store_name"`
	Username  string `json:"user_name"`
	Password  string `json:"password"`

	PoolSize     int `json:"pool_size"`
	IdleTimeout  int `json:"idle_timeout"`
	ReadTimeout  int `json:"read_timeout"`
	WriteTimeout int `json:"write_timeout"`
	DialTimeout  int `json:"dial_timeout"`
}

// NewService returns an instance of redis.Pool with all the required configurations set
func NewService(cfg *Config) (*redis.Pool, error) {
	db, _ := strconv.Atoi(cfg.StoreName)
	rpool := &redis.Pool{
		MaxIdle:         cfg.PoolSize,
		MaxActive:       cfg.PoolSize,
		IdleTimeout:     time.Duration(cfg.IdleTimeout) * time.Second,
		Wait:            true,
		MaxConnLifetime: time.Duration(cfg.DialTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
				redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Second),
				redis.DialPassword(cfg.Password),
				redis.DialConnectTimeout(time.Duration(cfg.DialTimeout)*time.Second),
				redis.DialDatabase(db),
			)
		},
	}

	conn := rpool.Get()
	rep, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	pong, _ := rep.(string)
	if pong != "PONG" {
		return nil, errors.New("ping failed")
	}
	conn.Close()

	return rpool, nil
}
