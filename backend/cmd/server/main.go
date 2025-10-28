package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"taskmaster-license/internal/api"
	"taskmaster-license/internal/config"
	"taskmaster-license/internal/database"
	"taskmaster-license/internal/middleware"
	"taskmaster-license/internal/repository"
	"taskmaster-license/internal/service"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbPath := config.AppConfig.GetDatabaseConnectionString()
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database initialized successfully at:", dbPath)

	// Generate root keys if they don't exist
	if err := ensureRootKeys(); err != nil {
		log.Printf("Warning: failed to ensure root keys: %v", err)
	}

	// Initialize repositories
	repo := repository.NewRepository(db)

	// Initialize services
	cmlService := service.NewCMLService(repo)
	siteService := service.NewSiteService(repo, cmlService)
	manifestService := service.NewManifestService(repo)

	// Initialize handlers
	authHandler := api.NewAuthHandler(config.AppConfig.JWTSecret)
	cmlHandler := api.NewCMLHandler(cmlService)
	siteHandler := api.NewSiteHandler(siteService)
	manifestHandler := api.NewManifestHandler(manifestService)
	ledgerHandler := api.NewLedgerHandler(repo)

	// Initialize license services (from merged license-server)
	licenseRepo := repository.NewLicenseRepository(db.Connection)
	licenseSiteService := service.NewLicenseSiteService(licenseRepo)
	licenseStatsService := service.NewLicenseStatsService(licenseRepo)
	licenseAlertService := service.NewLicenseAlertService(licenseRepo)
	licenseHandler := api.NewLicenseHandler(licenseSiteService, licenseStatsService, licenseAlertService, config.AppConfig.JWTSecret)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Authentication
	router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(config.AppConfig.JWTSecret))

	// CML management
	cmlGroup := protected.Group("/cml")
	{
		cmlGroup.POST("/upload", cmlHandler.UploadCML)
		cmlGroup.GET("", cmlHandler.GetCML)
		cmlGroup.POST("/refresh", cmlHandler.RefreshCML)
	}

	// Site license management
	sitesGroup := protected.Group("/sites")
	{
		sitesGroup.POST("/create", siteHandler.CreateSite)
		sitesGroup.GET("", siteHandler.ListSites)
		sitesGroup.GET("/:site_id", siteHandler.GetSite)
		sitesGroup.DELETE("/:site_id", siteHandler.DeleteSite)
		sitesGroup.POST("/:site_id/heartbeat", siteHandler.Heartbeat)
	}

	// License validation (public endpoint)
	router.POST("/api/validate", siteHandler.Validate)

	// Usage ledger
	protected.GET("/ledger", ledgerHandler.GetLedger)

	// Manifest management
	manifestsGroup := protected.Group("/manifests")
	{
		manifestsGroup.POST("/generate", manifestHandler.GenerateManifest)
		manifestsGroup.GET("", manifestHandler.ListManifests)
		manifestsGroup.GET("/:manifest_id", manifestHandler.GetManifest)
		manifestsGroup.GET("/:manifest_id/download", manifestHandler.DownloadManifest)
		manifestsGroup.POST("/send", manifestHandler.SendManifest)
	}

	// License key management (from merged license-server)
	protected.POST("/keys/create", licenseHandler.CreateSiteKey)
	protected.GET("/keys", licenseHandler.ListSiteKeys)
	protected.PUT("/keys/:id", licenseHandler.UpdateSiteKey)
	protected.POST("/keys/refresh", licenseHandler.RefreshKey)
	protected.POST("/keys/validate", licenseHandler.ValidateKey)
	protected.POST("/stats/aggregate", licenseHandler.AggregateStats)
	protected.POST("/alerts", licenseHandler.SendAlert)

	// Start server
	addr := config.AppConfig.ServerAddress()
	fmt.Printf("Starting server on %s\n", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func ensureRootKeys() error {
	keysDir := "keys"
	if err := os.MkdirAll(keysDir, 0755); err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	rootPrivatePath := filepath.Join(keysDir, "root_private.pem")
	rootPublicPath := filepath.Join(keysDir, "root_public.pem")

	// Check if keys already exist
	if _, err := os.Stat(rootPrivatePath); err == nil {
		// Keys exist, check public key
		if _, err := os.Stat(rootPublicPath); err == nil {
			return nil // Both keys exist
		}
	}

	// Keys don't exist - user should run genkeys command
	fmt.Println("Root keys not found. Please run: go run cmd/genkeys/main.go root")
	
	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
