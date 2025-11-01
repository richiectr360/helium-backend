package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// interpolateTemplate replaces {l10n.key} patterns with actual localized values
// Optimized with single-pass using strings.Builder
func interpolateTemplate(template string, localizedData map[string]string) string {
	var builder strings.Builder
	// Pre-allocate capacity: template size + estimated replacements
	builder.Grow(len(template) + len(localizedData)*50)

	i := 0
	for i < len(template) {
		// Look for {l10n. pattern
		if i+6 < len(template) && template[i] == '{' && template[i:i+6] == "{l10n." {
			// Find the closing }
			startKey := i + 6
			endKey := strings.IndexByte(template[startKey:], '}')
			if endKey == -1 {
				// No closing brace found, write the current char and continue
				builder.WriteByte(template[i])
				i++
				continue
			}

			key := template[startKey : startKey+endKey]
			if value, ok := localizedData[key]; ok {
				// Escape quotes in value
				escapedValue := strings.ReplaceAll(value, `"`, `\"`)
				// Write the replacement
				builder.WriteString(`"`)
				builder.WriteString(escapedValue)
				builder.WriteString(`"`)
				i = startKey + endKey + 1
				continue
			}
		}
		builder.WriteByte(template[i])
		i++
	}
	return builder.String()
}

// getLocalizedComponent generates a localized React component
func getLocalizedComponent(componentType, lang string) (*LocalizedComponent, error) {
	template, exists := componentTemplates[componentType]
	if !exists {
		return nil, fmt.Errorf("component type '%s' not found", componentType)
	}

	// Get localized strings, fallback to English
	strings, exists := localizationDB[lang]
	if !exists {
		strings = localizationDB["en"]
	}

	// Get only the required keys for this component
	componentStrings := make(map[string]string)
	for _, key := range template.RequiredKeys {
		if value, ok := strings[key]; ok {
			componentStrings[key] = value
		} else {
			componentStrings[key] = fmt.Sprintf("[%s]", key)
		}
	}

	// Interpolate template
	localizedTemplate := interpolateTemplate(template.Template, componentStrings)

	// Generate component ID with full timestamp to avoid collisions
	componentID := fmt.Sprintf("%s_%s_%d", componentType, lang, time.Now().UnixMilli())

	return &LocalizedComponent{
		ComponentName: template.ComponentName,
		ComponentType: template.ComponentType,
		Language:      lang,
		Template:      localizedTemplate,
		LocalizedData: componentStrings,
		Metadata: ComponentMetadata{
			ComponentID:  componentID,
			LastUpdated:  time.Now().Format(time.RFC3339),
			RequiredKeys: template.RequiredKeys,
		},
	}, nil
}

// Health check handler
func healthCheck(c *gin.Context) {
	redisStatus := "disconnected"
	if redisClient != nil {
		timeoutCtx, cancel := context.WithTimeout(ctx, RedisTimeout)
		defer cancel()
		if err := redisClient.Ping(timeoutCtx).Err(); err == nil {
			redisStatus = "connected"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            "healthy",
		"service":           "localization-manager-backend",
		"version":           "0.1.0",
		"cache_size":        componentCache.Size(),
		"concurrency_limit": getConcurrencyLimit(),
		"redis_status":      redisStatus,
	})
}

// Get localized component handler
func getLocalizedComponentEndpoint(c *gin.Context) {
	componentType := c.Param("component_type")
	lang := c.DefaultQuery("lang", "en")

	// Validate language code
	if _, exists := localizationDB[lang]; !exists {
		availableLangs := make([]string, 0, len(localizationDB))
		for key := range localizationDB {
			availableLangs = append(availableLangs, key)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":               fmt.Sprintf("invalid language code: %s", lang),
			"available_languages": availableLangs,
		})
		c.Abort()
		return
	}

	cacheKey := fmt.Sprintf("component:%s:%s", componentType, lang)

	// Check TTL cache first
	if cached, found := componentCache.Get(cacheKey); found {
		component := cached.(*LocalizedComponent)
		// Get already updated LRU, no need to call Put again
		// Do NOT write to Redis here - Redis should only be written on cache misses/generation
		response := *component
		response.Cached = true
		c.JSON(http.StatusOK, response)
		return
	}

	// TTL cache miss, check Redis
	if redisClient != nil {
		component, err := getFromRedis(cacheKey)
		if err == nil && component != nil {
			// Found in Redis, store in TTL cache
			componentCache.Put(cacheKey, component)

			// Refresh Redis TTL asynchronously to not block response
			go func() {
				setInRedis(cacheKey, component)
			}()

			response := *component
			response.Cached = true
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Both caches missed, generate component
	component, err := getLocalizedComponent(componentType, lang)
	if err != nil {
		availableComponents := make([]string, 0, len(componentTemplates))
		for key := range componentTemplates {
			availableComponents = append(availableComponents, key)
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error":                err.Error(),
			"message":              "Component type not found",
			"available_components": availableComponents,
		})
		return
	}

	// Store in both caches
	componentCache.Put(cacheKey, component)
	if redisClient != nil {
		setInRedis(cacheKey, component)
	}

	component.Cached = false
	c.JSON(http.StatusOK, component)
}

