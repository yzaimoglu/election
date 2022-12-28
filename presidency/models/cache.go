package models

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yzaimoglu/election/presidency/utilities"
)

// RedisConnection Options struct
type RedisConnectionOptions struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

// RedisConnection pool
var redisConnectionPool *redis.Pool

// Setup the cache
func SetupCache() {
	redisConnectionPool = newRedisPool()
}

// Initialize new redis pool
func newRedisPool() *redis.Pool {
	// Set options
	options := RedisConnectionOptions{
		Host:        utilities.GetEnv("CB_CACHE_HOST", "localhost") + ":" + utilities.GetEnv("CB_CACHE_PORT", "6379"),
		Password:    utilities.GetEnv("CB_CACHE_PASSWORD", ""),
		MaxIdle:     80,
		MaxActive:   12000,
		IdleTimeout: 240 * time.Second,
	}

	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", options.Host)
			if err != nil {
				return nil, err
			}
			if options.Password != "" {
				if _, err := c.Do("AUTH", options.Password); err != nil {
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

// Set the value for a key in redis
func RedisSet(key string, value interface{}) error {
	client := redisConnectionPool.Get()
	defer client.Close()
	_, err := client.Do("SET", key, value)
	if err != nil {
		return err
	}
	return nil
}

// Set the ttl for a key in redis
func RedisTTL(key string, ttl int) error {
	client := redisConnectionPool.Get()
	defer client.Close()
	_, err := client.Do("EXPIRE", key, ttl)
	if err != nil {
		return err
	}
	return nil
}

// Get a value from redis
func RedisGet(key string) ([]byte, error) {
	client := redisConnectionPool.Get()
	defer client.Close()
	value, err := redis.Bytes(client.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Get a slice of values from redis
func RedisGetSlice(key string) ([][]byte, error) {
	client := redisConnectionPool.Get()
	defer client.Close()
	value, err := redis.ByteSlices(client.Do("LRANGE", key, 0, -1))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Check if a key exists in redis
func RedisExists(key string) (bool, error) {
	client := redisConnectionPool.Get()
	defer client.Close()
	value, err := redis.Bool(client.Do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return value, nil
}
