package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConcurrencyLimiter middleware to limit concurrent requests
func ConcurrencyLimiter(limit int) gin.HandlerFunc {
	semaphore := make(chan struct{}, limit)

	return func(c *gin.Context) {
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			c.Next()
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "server is at capacity, please try again later",
			})
			c.Abort()
		}
	}
}

