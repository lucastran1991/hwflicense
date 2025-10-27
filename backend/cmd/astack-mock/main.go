package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CMLIssueRequest struct {
	OrgID       string   `json:"org_id" binding:"required"`
	MaxSites    int      `json:"max_sites" binding:"required"`
	Validity    string   `json:"validity" binding:"required"`
	FeaturePacks []string `json:"feature_packs"`
	KeyType     string   `json:"key_type"`
}

type ManifestReceiveRequest struct {
	OrgID    string                 `json:"org_id"`
	Period   string                 `json:"period"`
	Manifest map[string]interface{} `json:"manifest"`
}

func main() {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "mock-astack"})
	})

	// Generate and sign CML
	router.POST("/api/cml/issue", func(c *gin.Context) {
		var req CMLIssueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Create mock CML data
		cmlData := map[string]interface{}{
			"type":           "customer_master_license",
			"org_id":         req.OrgID,
			"max_sites":      req.MaxSites,
			"validity":       req.Validity,
			"feature_packs":  req.FeaturePacks,
			"key_type":       req.KeyType,
			"issued_by":      "astack_root",
			"issuer_public_key": "mock-key",
			"issued_at":      time.Now().Format(time.RFC3339),
		}

		// Generate mock signature (in production, this would be a real ECDSA signature)
		signature := fmt.Sprintf("mock_signature_%s", uuid.New().String())

		c.JSON(200, gin.H{
			"cml": map[string]interface{}{
				"cml_data": cmlData,
				"signature": signature,
			},
			"status": "issued",
		})
	})

	// Receive and validate manifest
	router.POST("/api/manifests/receive", func(c *gin.Context) {
		var req ManifestReceiveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Received manifest from org: %s, period: %s", req.OrgID, req.Period)
		log.Printf("Manifest data: %v", req.Manifest)

		// Mock validation
		c.JSON(200, gin.H{
			"status": "received",
			"validated": true,
			"message": "Manifest successfully validated and recorded",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// List received manifests (for testing)
	router.GET("/api/manifests/received", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"manifests": []map[string]interface{}{},
			"message": "Mock server - no persistent storage",
		})
	})

	fmt.Println("Starting Mock A-Stack Server on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

