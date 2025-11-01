package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	// Initialize Redis
	redisClient = initRedis()
	defer redisClient.Close()

	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("⚠️  Redis connection failed: %v (continuing without Redis)\n", err)
		redisClient = nil
	} else {
		fmt.Println("✅ Redis connected successfully")
	}

	router := gin.Default()

	// Set trusted proxies for production deployments
	router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})

	// Apply concurrency limiter middleware
	router.Use(ConcurrencyLimiter(getConcurrencyLimit()))

	// Routes
	router.GET("/health", healthCheck)
	router.GET("/api/component/:component_type", getLocalizedComponentEndpoint)

	// Start server
	fmt.Println("🚀 Localization Manager Backend starting on :8000")
	fmt.Println("📚 Available components:", strings.Join(getComponentKeys(), ", "))
	fmt.Println("🌍 Supported languages:", strings.Join(getLanguageKeys(), ", "))

	if err := router.Run(":8000"); err != nil {
		panic(err)
	}
}
