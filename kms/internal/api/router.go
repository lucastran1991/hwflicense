package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(handler *Handler) *gin.Engine {
	// Set Gin to release mode (disable debug)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply global middleware
	// CORS must be first to handle preflight requests
	router.Use(CORSMiddleware())
	router.Use(RecoveryMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(RateLimitMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	v1 := router.Group("/keys")
	{
		v1.GET("", handler.ListKeys)                    // List all keys (must be before /:id routes)
		v1.GET("/:id/download", handler.DownloadKey)    // Download key (must be before /:id routes)
		v1.POST("", handler.RegisterKey)
		v1.POST("/validate", handler.ValidateKey)
		v1.POST("/:id/refresh", handler.RefreshKey)
		v1.DELETE("/:id", handler.RemoveKey)
	}

	// License routes
	licenses := router.Group("/licenses")
	{
		licenses.POST("/generate", handler.GenerateLicense)
		licenses.POST("/validate", handler.ValidateLicense)
	}

	return router
}

