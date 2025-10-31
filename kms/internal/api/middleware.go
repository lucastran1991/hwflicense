package api

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// Rate limiter state
	rateLimiters = make(map[string]*rateLimiter)
	rateLimitMu  sync.RWMutex
)

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	tokens   int
	lastTime time.Time
	mu       sync.Mutex
}

const (
	// MaxRequests is the maximum number of requests per window
	MaxRequests = 100
	// WindowSize is the size of the rate limiting window
	WindowSize = 1 * time.Minute
)

// RateLimitMiddleware implements rate limiting to prevent brute force attacks
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rateLimitMu.RLock()
		limiter, exists := rateLimiters[clientIP]
		rateLimitMu.RUnlock()

		if !exists {
			rateLimitMu.Lock()
			limiter = &rateLimiter{
				tokens:   MaxRequests,
				lastTime: time.Now(),
			}
			rateLimiters[clientIP] = limiter
			rateLimitMu.Unlock()
		}

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		// Refill tokens based on time passed
		now := time.Now()
		elapsed := now.Sub(limiter.lastTime)
		if elapsed >= WindowSize {
			limiter.tokens = MaxRequests
			limiter.lastTime = now
		}

		// Check if we have tokens available
		if limiter.tokens <= 0 {
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		// Consume a token
		limiter.tokens--

		c.Next()
	}
}

// LoggingMiddleware logs requests (but never logs key material)
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Never log request body or key material - only metadata
		log.Printf("[%s] %s %s - %d - %v - %s", clientIP, method, path, status, latency, c.Errors.String())
	}
}

// RecoveryMiddleware handles panics and returns JSON error responses
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, recovered interface{}) {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		c.AbortWithStatus(500)
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Allow localhost origins (development)
		// Supports both localhost and 127.0.0.1 with any port
		allowed := false
		if origin != "" {
			// Check if origin starts with localhost or 127.0.0.1
			if strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasPrefix(origin, "https://127.0.0.1:") {
				allowed = true
			}
		}
		
		// Handle preflight requests first
		if c.Request.Method == "OPTIONS" {
			if allowed && origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			} else if origin == "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
			
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(204)
			return
		}
		
		// Set CORS headers for actual requests
		if allowed && origin != "" {
			// Set specific origin (required when credentials are enabled)
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		} else if origin == "" {
			// No origin header means same-origin request or non-browser client
			// Allow with wildcard for non-browser API clients
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		c.Next()
	}
}

