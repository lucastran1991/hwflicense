package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"taskmaster-license/internal/service"
)

type ManifestHandler struct {
	manifestService *service.ManifestService
}

func NewManifestHandler(manifestService *service.ManifestService) *ManifestHandler {
	return &ManifestHandler{manifestService: manifestService}
}

type GenerateManifestRequest struct {
	Period string `json:"period" binding:"required"`
}

func (h *ManifestHandler) GenerateManifest(c *gin.Context) {
	var req GenerateManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manifest, err := h.manifestService.GenerateManifest(req.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"manifest":  string(manifest.ManifestData),
		"signature": manifest.Signature,
	})
}

func (h *ManifestHandler) ListManifests(c *gin.Context) {
	period := c.Query("period")

	manifests, err := h.manifestService.ListManifests(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manifests": manifests})
}

func (h *ManifestHandler) GetManifest(c *gin.Context) {
	id := c.Param("manifest_id")

	manifest, err := h.manifestService.GetManifest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"manifest": manifest})
}

func (h *ManifestHandler) DownloadManifest(c *gin.Context) {
	id := c.Param("manifest_id")

	manifest, err := h.manifestService.GetManifest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=manifest_%s.json", id))
	c.String(http.StatusOK, string(manifest.ManifestData))
}

type SendManifestRequest struct {
	ManifestID  string `json:"manifest_id" binding:"required"`
	AStackEndpoint string `json:"astack_endpoint" binding:"required"`
}

func (h *ManifestHandler) SendManifest(c *gin.Context) {
	var req SendManifestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get manifest
	manifest, err := h.manifestService.GetManifest(req.ManifestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Parse manifest data
	var manifestData map[string]interface{}
	if err := json.Unmarshal(manifest.ManifestData, &manifestData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid manifest data"})
		return
	}

	// Prepare payload for A-Stack
	payload := map[string]interface{}{
		"org_id":    manifest.OrgID,
		"period":    manifest.Period,
		"manifest":  manifestData,
		"signature": manifest.Signature,
	}

	// Send to A-Stack with retry logic
	err = sendToAStackWithRetry(req.AStackEndpoint, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mark as sent in database
	now := time.Now()
	if err := h.manifestService.MarkManifestSent(manifest.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark manifest as sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "sent",
		"timestamp": now.Format(time.RFC3339),
		"message":   "Manifest successfully sent to A-Stack",
	})
}

// sendToAStackWithRetry sends manifest to A-Stack with exponential backoff retry logic
func sendToAStackWithRetry(endpoint string, payload map[string]interface{}) error {
	maxRetries := 3
	initialDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			delay := initialDelay * time.Duration(1<<uint(attempt-1))
			time.Sleep(delay)
		}

		// Marshal payload
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		// Send request
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("failed to send after %d attempts: %w", maxRetries, err)
			}
			continue
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Read response body
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("A-Stack response: %s\n", string(body))
			return nil
		}

		// If not success, read error and retry
		body, _ := io.ReadAll(resp.Body)
		if attempt == maxRetries-1 {
			return fmt.Errorf("A-Stack returned status %d: %s", resp.StatusCode, string(body))
		}
	}

	return fmt.Errorf("max retries exceeded")
}
