package main

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultConcurrencyLimit = 20
	CacheMaxSize            = 50
	CacheTTL                = 10 * time.Minute
	RedisTTL                = 30 * time.Minute
	RedisTimeout            = 2 * time.Second
)

// getConcurrencyLimit returns the concurrency limit from env var or default
func getConcurrencyLimit() int {
	if limit := os.Getenv("CONCURRENCY_LIMIT"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			return l
		}
	}
	return DefaultConcurrencyLimit
}

