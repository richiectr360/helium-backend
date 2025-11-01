package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"
)

// Global Redis client and context
var redisClient *redis.Client
var ctx = context.Background()

// Global cache instance
var componentCache = NewTTLCache(CacheMaxSize, CacheTTL)

// initRedis initializes the Redis client
func initRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return client
}

// getFromRedis retrieves a component from Redis
func getFromRedis(key string) (*LocalizedComponent, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisTimeout)
	defer cancel()

	val, err := redisClient.Get(timeoutCtx, key).Result()
	if err != nil {
		return nil, err
	}

	var component LocalizedComponent
	if err := json.Unmarshal([]byte(val), &component); err != nil {
		return nil, err
	}

	return &component, nil
}

// setInRedis stores a component in Redis with TTL
func setInRedis(key string, component *LocalizedComponent) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisTimeout)
	defer cancel()

	data, err := json.Marshal(component)
	if err != nil {
		return err
	}

	return redisClient.Set(timeoutCtx, key, data, RedisTTL).Err()
}

